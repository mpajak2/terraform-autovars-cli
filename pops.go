package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: pops <plan|apply|output> <env>")
		os.Exit(1)
	}
	terraformAction := os.Args[1]
	envName := os.Args[2]

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	stackName := filepath.Base(cwd)
	stacksDir := filepath.Dir(cwd)
	projectRoot := filepath.Dir(stacksDir)
	envVarsDir := filepath.Join(projectRoot, "ENVVARS", envName)
	envVarsSecuredDir := filepath.Join(envVarsDir, "secured")

	fmt.Printf("Detected stack: %s\n", stackName)
	fmt.Printf("Environment: %s\n", envName)
	fmt.Printf("Project root: %s\n", projectRoot)
	fmt.Printf("ENVVARS dir: %s\n", envVarsDir)
	fmt.Printf("ENVVARS secured dir: %s\n", envVarsSecuredDir)

	varFiles, err := findVarFiles(envVarsDir, stackName)
	if err != nil {
		fmt.Printf("Error scanning directory %s: %v\n", envVarsDir, err)
		os.Exit(1)
	}

	securedFiles, err := findVarFiles(envVarsSecuredDir, stackName)
	if err != nil {
		fmt.Printf("Warning: could not scan secured directory: %v\n", err)
		securedFiles = []string{}
	}

	for _, securedFilePath := range securedFiles {
		decryptedFile, err := decryptSopsFile(securedFilePath)
		if err != nil {
			fmt.Printf("Error decrypting file %s: %v\n", securedFilePath, err)
			os.Exit(1)
		}
		varFiles = append(varFiles, decryptedFile)
	}

	terraformArgs := []string{terraformAction}
	for _, file := range varFiles {
		terraformArgs = append(terraformArgs, fmt.Sprintf("-var-file=%s", file))
	}

	fmt.Printf("Running command: terraform %s\n", strings.Join(terraformArgs, " "))
	cmd := exec.Command("terraform", terraformArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running terraform: %v\n", err)
		os.Exit(1)
	}
}

func findVarFiles(dir, stackName string) ([]string, error) {
	var matchingFiles []string

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.Contains(strings.ToLower(name), strings.ToLower(stackName)) && strings.HasSuffix(name, ".json") {
			matchingFiles = append(matchingFiles, filepath.Join(dir, name))
		}
	}
	return matchingFiles, nil
}

func decryptSopsFile(filePath string) (string, error) {
	cmd := exec.Command("sops", "-d", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file: %w", err)
	}

	tmpFile, err := ioutil.TempFile("", "decrypted-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(output); err != nil {
		return "", fmt.Errorf("failed to write decrypted data: %w", err)
	}

	return tmpFile.Name(), nil
}
