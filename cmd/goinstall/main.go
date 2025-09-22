package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	modulePath := "github.com/username/" + project // auto module path

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
	cloneCmd := exec.Command("git", "clone", "--depth", "1", skeletonRepo, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("Failed to clone skeleton: %v\n", err)
		return
	}

	// 3️⃣ Copy all needed files
	for _, name := range projectFolders {
		src := filepath.Join(tempDir, name)
		dest := filepath.Join(project, name)
		if _, err := os.Stat(src); err == nil {
			if err := os.Rename(src, dest); err != nil {
				fmt.Printf("Move failed for %s: %v\n", name, err)
			}
		}
	}
	os.RemoveAll(tempDir)

	// 4️⃣ Recursive replace inside project files
	filepath.WalkDir(project, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		// শুধু go, yaml, sh, txt, md ফাইলগুলোতে replace করা
		ext := filepath.Ext(path)
		if ext == ".go" || ext == ".yaml" || ext == ".sh" || ext == ".txt" || ext == ".md" {
			data, _ := os.ReadFile(path)
			content := string(data)
			content = strings.ReplaceAll(content, "your_project", project)
			content = strings.ReplaceAll(content, "your/module/path", modulePath)
			os.WriteFile(path, []byte(content), 0644)
		}
		return nil
	})

	// 5️⃣ Initialize go.mod
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

	// 6️⃣ Add Golara dependency
	cmdGet := exec.Command("go", "get", "github.com/aasoft24/golara@latest")
	cmdGet.Dir = project
	cmdGet.Stdout = os.Stdout
	cmdGet.Stderr = os.Stderr
	_ = cmdGet.Run()

	// 7️⃣ Run go mod tidy
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = project
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	_ = cmdTidy.Run()

	fmt.Printf("🚀 Project '%s' created successfully!\n", project)
	fmt.Printf("Module Path: %s\n", modulePath)
	fmt.Printf("Run: cd %s && go run main.go\n", project)
}
