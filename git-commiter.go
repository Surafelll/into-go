package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func generateRandomCode(index int) string {
	codeSnippets := []string{
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Hello, world!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Random commit bot in action!\") }", index),
		fmt.Sprintf("func Run_%d() { fmt.Println(\"Automating Git commits!\") }", index),
	}
	return codeSnippets[rand.Intn(len(codeSnippets))]
}

func commitAndPush() {
	cmds := [][]string{
		{"git", "add", "."}, // Add all files, not just committer.go
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
	dateInput = dateInput + " 12:00:00"
	os.Setenv("GIT_COMMITTER_DATE", dateInput)
	os.Setenv("GIT_AUTHOR_DATE", dateInput)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter commit date (YYYY-MM-DD): ")
	dateInput, _ := reader.ReadString('\n')
	dateInput = dateInput[:len(dateInput)-1] // Remove newline character

	for i := 1; i <= 10; i++ {
		code := generateRandomCode(i)
		file, _ := os.Create(fmt.Sprintf("committer_%d.go", i))
		defer file.Close()
		file.WriteString(fmt.Sprintf("package main\n\nimport \"fmt\"\n\n%s\n", code))

		setCommitDate(dateInput)
		commitAndPush()
	}

	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
