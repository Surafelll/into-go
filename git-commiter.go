package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func generateRandomCode() string {
	codeSnippets := []string{
		"fmt.Println(\"Hello, world!\")",
		"fmt.Println(\"Random commit bot in action!\")",
		"fmt.Println(\"Automating Git commits!\")",
	}
	return codeSnippets[rand.Intn(len(codeSnippets))]
}

func commitAndPush() {
	cmds := [][]string{
		{"git", "add", "committer.go"},
		{"git", "commit", "-m", "Automated commit"},
	}
	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func setCommitDate(dateInput string) {
	dateInput = dateInput + " 12:00:00" // Append default time to match Git format
	os.Setenv("GIT_COMMITTER_DATE", dateInput)
	os.Setenv("GIT_AUTHOR_DATE", dateInput)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter commit date (YYYY-MM-DD): ")
	dateInput, _ := reader.ReadString('\n')
	dateInput = dateInput[:len(dateInput)-1] // Remove newline character

	for i := 1; i <= 10; i++ {
		code := generateRandomCode()
		file, _ := os.Create("committer.go") // This will overwrite the file on each iteration
		defer file.Close()
		file.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n\t" + code + "\n}\n") // Add the generated code without 'main' conflict

		// Set the commit date and make the commit
		setCommitDate(dateInput)
		commitAndPush()

		// Add different commit messages
		cmd := exec.Command("git", "commit", "--amend", "--no-edit", "-m", fmt.Sprintf("Automated commit %d", i))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		// Update the commit date for each commit
		dateInput = dateInput // Increment the date if necessary
	}

	// Finally push all the changes after the loop
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
