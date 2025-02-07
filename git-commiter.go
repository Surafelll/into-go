package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CommitRequest struct {
	Date   string `json:"date"`
	Author string `json:"author"`
}

// Generate random Go code for commits
func generateRandomCode(index int) string {
	timestamp := fmt.Sprintf("// Timestamp: %d\n", time.Now().UnixNano()) // Unique identifier
	codeSnippets := []string{
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Hello, world!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Automating Git commits!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Random commit bot in action!\") }", index),
	}
	randomSnippet := codeSnippets[rand.Intn(len(codeSnippets))]

	return timestamp + randomSnippet
}

// Set commit date uniquely per commit
func setCommitDate(dateInput string, commitIndex int) {
	dateInput = fmt.Sprintf("%s %02d:00:00", dateInput, commitIndex)
	os.Setenv("GIT_COMMITTER_DATE", dateInput)
	os.Setenv("GIT_AUTHOR_DATE", dateInput)
}

// Cool animated progress bar
func showProgressBar(current, total int) {
	width := 30 // Progress bar width
	progress := int(float64(current) / float64(total) * float64(width))
	bar := "[" + strings.Repeat("█", progress) + strings.Repeat("-", width-progress) + "]"
	fmt.Printf("\r%s %d/%d commits", bar, current, total)
}

// Save commit history to a structured folder
func saveCommitHistory(author string, commitMessage string) {
	historyPath := filepath.Join("history", time.Now().Format("2006-01-02")) // history/YYYY-MM-DD
	if err := os.MkdirAll(historyPath, os.ModePerm); err != nil {
		fmt.Println("⚠️ Error creating history folder:", err)
		return
	}

	filename := filepath.Join(historyPath, fmt.Sprintf("%s.txt", time.Now().Format("15-04-05")))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("⚠️ Error saving commit history:", err)
		return
	}
	defer file.Close()

	// Write commit details
	file.WriteString(fmt.Sprintf("Author: %s\nTime: %s\nMessage: %s\n\n", author, time.Now().Format(time.RFC1123), commitMessage))
}

// Commit and push changes with clean console output
func commitAndPush(author string, commitIndex int, totalCommits int) {
	commitMessage := fmt.Sprintf("🚀 Automated commit %d by %s", commitIndex, author)

	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", commitMessage},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Run() // Run command silently
	}

	showProgressBar(commitIndex, totalCommits) // Update progress bar
	saveCommitHistory(author, commitMessage)  // Save commit message
	time.Sleep(300 * time.Millisecond)        // Smooth animation effect
}

// API endpoint to handle commit requests
func commitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Date == "" || req.Author == "" {
		http.Error(w, "Missing date or author", http.StatusBadRequest)
		return
	}

	totalCommits := 10
	fmt.Println("\n🚀 **Starting Automated Commit Process** 🚀\n")

	for i := 1; i <= totalCommits; i++ {
		code := generateRandomCode(i)
		fileName := fmt.Sprintf("committer_%d.go", i)

		file, _ := os.Create(fileName)
		file.WriteString(fmt.Sprintf("package main\n\nimport \"fmt\"\n\n%s\n", code))
		file.Close()

		setCommitDate(req.Date, i)
		commitAndPush(req.Author, i, totalCommits)
	}

	fmt.Println("\n\n✅ All commits completed! Pushing to remote...\n")

	// Push all commits silently (no extra console noise)
	cmd := exec.Command("git", "push")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()

	fmt.Println("\n🎉 **Commits Successfully Pushed!** 🎉")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Commits pushed successfully"})
}

func main() {
	http.HandleFunc("/commit", commitHandler)

	fmt.Println("\n🌍 Server is running on port 8080... 🌍")
	fmt.Println("🔗 Send a POST request to http://localhost:8080/commit")
	fmt.Println("💾 Example JSON Payload: {\"date\": \"2025-02-05\", \"author\": \"Surafel\"}\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("❌ Error starting server:", err)
	}
}
