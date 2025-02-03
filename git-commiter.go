package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type CommitRequest struct {
	Date   string `json:"date"`
	Author string `json:"author"`
}

// Generate unique file content to ensure Git recognizes changes
func generateRandomCode(index int) string {
	timestamp := fmt.Sprintf("// Timestamp: %d\n", time.Now().UnixNano()) // Unique identifier
	codeSnippets := []string{
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Hello, world!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Random commit bot in action!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Automating Git commits!\") }", index),
	}
	randomSnippet := codeSnippets[rand.Intn(len(codeSnippets))]

	return timestamp + randomSnippet
}

// Set Git commit date uniquely per commit
func setCommitDate(dateInput string, commitIndex int) {
	dateInput = fmt.Sprintf("%s %02d:00:00", dateInput, commitIndex) // Different hour per commit
	os.Setenv("GIT_COMMITTER_DATE", dateInput)
	os.Setenv("GIT_AUTHOR_DATE", dateInput)
}

// Commit and push changes
func commitAndPush(author string, commitIndex int) {
	commitMessage := fmt.Sprintf("Automated commit %d by %s", commitIndex, author)

	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", commitMessage},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
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

	for i := 1; i <= 10; i++ {
		code := generateRandomCode(i)
		fileName := fmt.Sprintf("committer_%d.go", i)

		file, err := os.Create(fileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating file: %v", err), http.StatusInternalServerError)
			return
		}

		_, writeErr := file.WriteString(fmt.Sprintf("package main\n\nimport \"fmt\"\n\n%s\n", code))
		file.Close()
		if writeErr != nil {
			http.Error(w, fmt.Sprintf("Error writing to file: %v", writeErr), http.StatusInternalServerError)
			return
		}

		setCommitDate(req.Date, i)
		commitAndPush(req.Author, i)
	}

	// Push all commits
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Commits pushed successfully"})
}

func main() {
	http.HandleFunc("/commit", commitHandler)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
