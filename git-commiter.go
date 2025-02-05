package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"io/ioutil"
)

// Struct to parse JSON request
type CommitRequest struct {
	RepoURL   string `json:"repo_url"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Date      string `json:"date,omitempty"`
	Author    string `json:"author"`
}

// Struct to hold commit messages
type CommitMessages struct {
	Messages []string `json:"messages"`
}

// Find an available port starting from a given number
func findAvailablePort(startPort int) int {
	for port := startPort; port < startPort+10; port++ { // Try 10 ports
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port
		}
	}
	return -1 // No available port found
}

// Get a random date between a range
func randomDateInRange(startDate, endDate string) string {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	randomTime := start.Add(time.Duration(rand.Int63n(int64(end.Sub(start)))))
	return randomTime.Format("2006-01-02")
}

// Set environment variables for commit date
func setCommitDate(date string, commitIndex int) {
	dateWithHour := fmt.Sprintf("%s %02d:00:00", date, commitIndex%24) // Spread commits across hours
	os.Setenv("GIT_COMMITTER_DATE", dateWithHour)
	os.Setenv("GIT_AUTHOR_DATE", dateWithHour)
}

// Clone or pull the repo
func cloneRepo(repoURL string) (string, error) {
	repoName := strings.TrimSuffix(filepath.Base(repoURL), ".git")
	repoPath := filepath.Join("repos", repoName)

	// If repo exists, pull latest changes
	if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
		fmt.Println("🔄 Repo exists. Pulling latest changes...")
		cmd := exec.Command("git", "-C", repoPath, "pull")
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to pull latest changes: %v", err)
		}
		return repoPath, nil
	}

	// Clone repo if it doesn't exist
	fmt.Println("📥 Cloning repository...")
	cmd := exec.Command("git", "clone", repoURL, repoPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone repo: %v", err)
	}

	return repoPath, nil
}

// Load commit messages from JSON file
func loadCommitMessages() ([]string, error) {
	file, err := ioutil.ReadFile("messages.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read messages.json: %v", err)
	}

	var commitMessages CommitMessages
	if err := json.Unmarshal(file, &commitMessages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commit messages: %v", err)
	}

	return commitMessages.Messages, nil
}

// Commit handler
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

	if req.RepoURL == "" || req.Author == "" {
		http.Error(w, "Missing repo URL or author", http.StatusBadRequest)
		return
	}

	// Load commit messages
	commitMessages, err := loadCommitMessages()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading commit messages: %v", err), http.StatusInternalServerError)
		return
	}

	// Determine commit dates
	var commitDates []string
	if req.Date != "" {
		commitDates = append(commitDates, req.Date) // Single date
	} else if req.StartDate != "" && req.EndDate != "" {
		start, _ := time.Parse("2006-01-02", req.StartDate)
		end, _ := time.Parse("2006-01-02", req.EndDate)
		for d := start; !d.After(end); d = d.Add(24 * time.Hour) {
			commitDates = append(commitDates, d.Format("2006-01-02"))
		}
	} else {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Clone repo
	repoPath, err := cloneRepo(req.RepoURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to clone repo: %v", err), http.StatusInternalServerError)
		return
	}

	// Create commit folder
	commitFolder := filepath.Join(repoPath, "commits", time.Now().Format("2006-01-02_15-04-05"))
	if err := os.MkdirAll(commitFolder, os.ModePerm); err != nil {
		http.Error(w, "Failed to create commit folder", http.StatusInternalServerError)
		return
	}

	// Start commit process
	totalCommits := len(commitDates) * 10 // 10 commits per day
	fmt.Println("\n🚀 **Starting Automated Commit Process** 🚀\n")

	commitIndex := 1
	for _, commitDate := range commitDates {
		for i := 0; i < 10; i++ { // Create 10 commits per day
			fileName := filepath.Join(commitFolder, fmt.Sprintf("commit_%d.go", commitIndex))
			file, _ := os.Create(fileName)
			file.WriteString(fmt.Sprintf("package main\n\n// Commit %d on %s\n", commitIndex, commitDate))
			file.Close()

			// Set commit date
			setCommitDate(commitDate, commitIndex)

			// Randomly select a commit message
			randomMessage := commitMessages[rand.Intn(len(commitMessages))]

			// Git add & commit with random message
			exec.Command("git", "-C", repoPath, "add", ".").Run()
			exec.Command("git", "-C", repoPath, "commit", "-m", fmt.Sprintf("%s - Commit %d on %s", randomMessage, commitIndex, commitDate)).Run()

			// Progress bar
			fmt.Printf("\r[██████████████████████████] %d/%d commits", commitIndex, totalCommits)
			commitIndex++
			time.Sleep(200 * time.Millisecond)
		}
	}

	// Push commits
	fmt.Println("\n✅ All commits completed! Pushing to remote...\n")
	exec.Command("git", "-C", repoPath, "push").Run()
	fmt.Println("\n🎉 **Commits Successfully Pushed!** 🎉")

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "repo": req.RepoURL})
}

// Start server
func startServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: addr, Handler: nil}

	fmt.Printf("\n🌍 Server is running on port %d... 🌍\n", port)
	fmt.Println("🔗 Send a POST request to http://localhost:" + fmt.Sprintf("%d", port) + "/commit")
	fmt.Println("💾 Example JSON Payload:")
	fmt.Println(`{"repo_url": "https://github.com/user/repo.git", "start_date": "2025-02-01", "end_date": "2025-02-10", "author": "Surafel"}`)
	fmt.Println("💾 Or for a single day:")
	fmt.Println(`{"repo_url": "https://github.com/user/repo.git", "date": "2025-02-05", "author": "Surafel"}`)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Find available port
	port := findAvailablePort(8080)
	if port == -1 {
		fmt.Println("Error: No available port found.")
		return
	}

	// Start HTTP server
	http.HandleFunc("/commit", commitHandler)
	startServer(port)
}
