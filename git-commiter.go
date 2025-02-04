package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

// Generate a random function with a unique name for each file
func generateRandomCode(index int) string {
	codeSnippets := []string{
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Hello, world!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Random commit bot in action!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Automating Git commits!\") }", index),
	}
	return codeSnippets[rand.Intn(len(codeSnippets))]
}

// Set the Git commit date
func setCommitDate(dateInput string) {
	dateInput = dateInput + " 12:00:00"
	os.Setenv("GIT_COMMITTER_DATE", dateInput)
	os.Setenv("GIT_AUTHOR_DATE", dateInput)
}

// Commit and push changes
func commitAndPush(author string, commitIndex int) {
	commitMessage := fmt.Sprintf("Automated commit %d by %s", commitIndex, author)

	cmds := [][]string{
		{"git", "add", "."}, // Add all files
		{"git", "commit", "-m", commitMessage},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get commit date from user
	fmt.Print("Enter commit date (YYYY-MM-DD): ")
	dateInput, _ := reader.ReadString('\n')
	dateInput = dateInput[:len(dateInput)-1] // Remove newline character

	// Get author name from user
	fmt.Print("Enter author name: ")
	author, _ := reader.ReadString('\n')
	author = author[:len(author)-1] // Remove newline character

	for i := 1; i <= 10; i++ {
		code := generateRandomCode(i)
		fileName := fmt.Sprintf("committer_%d.go", i)

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", fileName, err)
			continue
		}

		_, writeErr := file.WriteString(fmt.Sprintf("package main\n\nimport \"fmt\"\n\n%s\n", code))
		file.Close()
		if writeErr != nil {
			fmt.Printf("Error writing to file %s: %v\n", fileName, writeErr)
			continue
		}

		setCommitDate(dateInput)
		commitAndPush(author, i)
	}

	// Push all commits after the loop
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
