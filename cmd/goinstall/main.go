package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const skeletonRepo = "https://github.com/aasoft24/golara.git"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goinstall <project-name>")
		return
	}

	project := os.Args[1]

	// 1️⃣ Create project folder if not exists
	if _, err := os.Stat(project); os.IsNotExist(err) {
		if err := os.Mkdir(project, 0755); err != nil {
			fmt.Printf("Failed to create project folder: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Project folder '%s' already exists, continuing...\n", project)
	}

	// 2️⃣ Clone skeleton repo
	tempDir := filepath.Join(project, "temp_skeleton")
	cloneCmd := exec.Command("git", "clone", skeletonRepo, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("Failed to clone skeleton: %v\n", err)
		return
	}

	// 3️⃣ Move skeleton files to project root
	files, _ := os.ReadDir(tempDir)
	for _, f := range files {
		os.Rename(filepath.Join(tempDir, f.Name()), filepath.Join(project, f.Name()))
	}
	os.RemoveAll(tempDir)

	// 4️⃣ Initialize Go module if not exists
	modFile := filepath.Join(project, "go.mod")
	fmt.Print("Enter Go module path (e.g. github.com/username/" + project + "): ")
	var modulePath string
	fmt.Scanln(&modulePath)
	if modulePath == "" {
		modulePath = project
	}

	if _, err := os.Stat(modFile); err == nil {
		fmt.Println("✅ go.mod already exists, skipping module init")
	} else {
		cmd := exec.Command("go", "mod", "init", modulePath)
		cmd.Dir = project
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to init go.mod: %v\n", err)
			return
		}
	}

	// 5️⃣ Run go mod tidy
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = project
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	cmdTidy.Run()

	fmt.Printf("✅ Project '%s' created successfully!\n", project)
	fmt.Printf("Run: cd %s && go run main.go\n", project)
}
