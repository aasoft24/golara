package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const skeletonRepo = "https://github.com/aasoft24/golara.git"

// Folders and files to copy
var projectFolders = []string{
	"app",
	"bootstrap",
	"resources",
	"config",
	"database",
	"routes",
	"public",
	"main.go",
	"config.yaml",
	"start.sh",
}

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

	// 2️⃣ Clone skeleton repo temporarily
	tempDir := filepath.Join(project, "temp_skeleton")
	cloneCmd := exec.Command("git", "clone", skeletonRepo, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("Failed to clone skeleton: %v\n", err)
		return
	}

	// 3️⃣ Copy only defined folders/files
	for _, name := range projectFolders {
		src := filepath.Join(tempDir, name)
		dest := filepath.Join(project, name)
		if _, err := os.Stat(src); err == nil {
			if err := copyRecursive(src, dest); err != nil {
				fmt.Printf("Failed to copy %s: %v\n", name, err)
			}
		}
	}

	// Remove temp skeleton
	os.RemoveAll(tempDir)

	// 4️⃣ Module path prompt
	fmt.Print("Enter Go module path (e.g. github.com/username/" + project + "): ")
	var modulePath string
	fmt.Scanln(&modulePath)
	if modulePath == "" {
		modulePath = project
	}
	fmt.Printf("✅ Using module path: %s\n", modulePath)

	// 5️⃣ Initialize Go module if not exists
	modFile := filepath.Join(project, "go.mod")
	if _, err := os.Stat(modFile); err == nil {
		fmt.Println("✅ go.mod already exists, skipping module init")
	} else {
		cmd := exec.Command("go", "mod", "init", modulePath)
		cmd.Dir = project
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️ Failed to init go.mod: %v\n", err)
		} else {
			fmt.Println("✅ go.mod initialized successfully")
		}
	}

	// 6️⃣ Add golara dependency
	addDep := exec.Command("go", "get", "github.com/aasoft24/golara@latest")
	addDep.Dir = project
	addDep.Stdout = os.Stdout
	addDep.Stderr = os.Stderr
	_ = addDep.Run()

	// 7️⃣ Run go mod tidy
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = project
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	_ = cmdTidy.Run()

	fmt.Printf("🚀 Project '%s' created successfully!\n", project)
	fmt.Printf("Run: cd %s && go run main.go\n", project)
}

// copyRecursive copies folders and files recursively
func copyRecursive(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if err := os.MkdirAll(dest, info.Mode()); err != nil {
			return err
		}
		entries, _ := os.ReadDir(src)
		for _, e := range entries {
			if err := copyRecursive(filepath.Join(src, e.Name()), filepath.Join(dest, e.Name())); err != nil {
				return err
			}
		}
	} else {
		from, _ := os.Open(src)
		defer from.Close()
		to, _ := os.Create(dest)
		defer to.Close()
		if _, err := io.Copy(to, from); err != nil {
			return err
		}
	}
	return nil
}
