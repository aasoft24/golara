package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const skeletonRepo = "https://github.com/aasoft24/golara.git"

func main() {
	fmt.Println("Goinstall CLI works!")
	if len(os.Args) < 2 {
		fmt.Println("Usage: goinstall <project-name>")
		return
	}

	project := os.Args[1]

	// 1️⃣ Create project folder
	if err := os.Mkdir(project, 0755); err != nil {
		log.Fatalf("Failed to create project folder: %v", err)
	}

	// 2️⃣ Clone skeleton repo
	cloneCmd := exec.Command("git", "clone", skeletonRepo, filepath.Join(project, "temp_skeleton"))
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		log.Fatalf("Failed to clone skeleton: %v", err)
	}

	// 3️⃣ Move skeleton files to project root
	temp := filepath.Join(project, "temp_skeleton")
	files, _ := os.ReadDir(temp)
	for _, f := range files {
		os.Rename(filepath.Join(temp, f.Name()), filepath.Join(project, f.Name()))
	}
	os.RemoveAll(temp)

	// 4️⃣ Initialize Go module (GitHub-ready)
	fmt.Print("Enter Go module path (e.g. github.com/username/" + project + "): ")
	var modulePath string
	fmt.Scanln(&modulePath)
	if modulePath == "" {
		modulePath = project // fallback to local-only module
	}

	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = project
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to init go.mod: %v", err)
	}

	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = project
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	cmdTidy.Run()

	fmt.Printf("✅ Project %s created successfully!\n", project)
	fmt.Printf("Run: cd %s && go run main.go\n", project)
}
