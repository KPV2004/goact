// cmd/create-goact-app/main.go
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// 💡 1. Embed the entire templates directory.
// The path is relative to this Go file.
//go:embed templates/*
var embeddedFiles embed.FS

// TemplateData holds dynamic values for template processing.
type TemplateData struct {
	ProjectName string
	ModuleName  string
}

func main() {
	// 2. Parse command-line arguments.
	if len(os.Args) < 2 {
		fmt.Println("❌ Error: Project name is required.")
		fmt.Println("Usage: create-goact-app <project-name>")
		os.Exit(1)
	}

	projectName := os.Args[1]
	// Basic module name defaults to the project name.
	moduleName := strings.ToLower(projectName) 
	
	targetDir := projectName

	// 3. Create the target directory.
	fmt.Printf("📂 Creating new Goact app in %s...\n", targetDir)
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Printf("❌ Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Data for processing template placeholders.
	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	// 4. Recursively walk through embedded files.
	// Note: go:embed prefixes paths with the embedded directory name.
	err = fs.WalkDir(embeddedFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Calculate destination path by removing the "templates" prefix.
		relPath := strings.TrimPrefix(path, "templates/")
		
		// If it's a template file, handle placeholders and rename.
		if strings.HasSuffix(relPath, ".tmpl") {
			destPath := filepath.Join(targetDir, strings.TrimSuffix(relPath, ".tmpl"))
			fmt.Printf("📄 Generating %s...\n", destPath)
			return processTemplate(path, destPath, data)
		}

		// For static files (like wasm_exec.js), copy directly.
		destPath := filepath.Join(targetDir, relPath)
		fmt.Printf("📄 Copying %s...\n", destPath)
		return copyStaticFile(path, destPath)
	})

	if err != nil {
		fmt.Printf("❌ Error scaffolding project: %v\n", err)
		os.Exit(1)
	}

	// 5. Final instructions.
	fmt.Println("\n✅ Success! Project scaffolded.")
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", targetDir)
	fmt.Println("  go mod tidy  (Note: You might need to adjust dependencies if goact is local)")
	fmt.Println("  GOOS=js GOARCH=wasm go build -o main.wasm main.go")
	fmt.Println("  python3 -m http.server 8080")
}

// processTemplate reads an embedded template, executes it with data, and writes to destination.
func processTemplate(srcPath, destPath string, data TemplateData) error {
	tmplContent, err := embeddedFiles.ReadFile(srcPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(srcPath).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// copyStaticFile reads an embedded static file and writes it directly to destination.
func copyStaticFile(srcPath, destPath string) error {
	content, err := embeddedFiles.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, content, 0644)
}
