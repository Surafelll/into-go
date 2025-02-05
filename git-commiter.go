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
		{"git", "push"},
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

	code := generateRandomCode()
	file, _ := os.Create("committer.go")
	defer file.Close()
	file.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n\t" + code + "\n}\n")

	setCommitDate(dateInput)
	commitAndPush()
}