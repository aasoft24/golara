package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const skeletonRepo = "https://github.com/aasoft24/golara.git"

var projectFolders = []string{
	"app", "bootstrap", "resources", "config", "database", "routes", "public",
	"main.go", "config.yaml", "start.sh",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goinstall <project-name>")
		return
	}

	project := os.Args[1]

	// 1️⃣ Create project folder
	if _, err := os.Stat(project); os.IsNotExist(err) {
		if err := os.Mkdir(project, 0755); err != nil {
			fmt.Printf("Failed to create project folder: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Project folder '%s' exists, continuing...\n", project)
	}

	// 2️⃣ Clone Golara repo temporarily
	tempDir := filepath.Join(project, "temp_skeleton")
	cloneCmd := exec.Command("git", "clone", skeletonRepo, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("Failed to clone skeleton: %v\n", err)
		return
	}

	// 3️⃣ Copy only required folders/files
	for _, name := range projectFolders {
		src := filepath.Join(tempDir, name)
		dest := filepath.Join(project, name)
		if _, err := os.Stat(src); err == nil {
			os.Rename(src, dest)
		}
	}

	// Remove temp skeleton
	os.RemoveAll(tempDir)

	// 4️⃣ Initialize go.mod
	fmt.Printf("Enter Go module path (e.g. github.com/username/%s): ", project)
	var modulePath string
	fmt.Scanln(&modulePath)

	modFile := filepath.Join(project, "go.mod")
	if _, err := os.Stat(modFile); os.IsNotExist(err) {
		cmd := exec.Command("go", "mod", "init", modulePath)
		cmd.Dir = project
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to init go.mod: %v\n", err)
		} else {
			fmt.Println("✅ go.mod initialized")
		}
	} else {
		fmt.Println("✅ go.mod already exists, skipping")
	}

	// 5️⃣ Add Golara dependency
	cmdGet := exec.Command("go", "get", "github.com/aasoft24/golara@latest")
	cmdGet.Dir = project
	cmdGet.Stdout = os.Stdout
	cmdGet.Stderr = os.Stderr
	_ = cmdGet.Run()

	// 6️⃣ Run go mod tidy
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = project
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	_ = cmdTidy.Run()

	fmt.Printf("🚀 Project '%s' created successfully!\n", project)
	fmt.Printf("Run: cd %s && go run main.go\n", project)
}
