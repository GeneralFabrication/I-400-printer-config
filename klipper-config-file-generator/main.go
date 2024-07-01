package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	homeDir, _   = os.UserHomeDir()
	configPath   = filepath.Join(homeDir, "printer_data/config/printer.cfg")
	configDir    = filepath.Join(homeDir, "printer_data/config")
	subdirTool   = "with-toolboard"
	subdirNoTool = "without-toolboard"
	branch       string
)

var availableBranches = []string{"v0.1", "v0.2"}

func main() {
	fmt.Println("Do you want to add, update, or rollback your config files?")
	fmt.Println("1. Add")
	fmt.Println("2. Update")
	fmt.Println("3. Rollback")

	var configChoice int
	fmt.Print("Enter your choice (1, 2, or 3): ")
	_, err := fmt.Scanf("%d", &configChoice)
	if err != nil || (configChoice < 1 || configChoice > 3) {
		fmt.Println("Invalid choice. Exiting.")
		panic(err)
	}

	switch configChoice {
	case 1:
		err = addConfig()
	case 2:
		err = updateConfig()
	case 3:
		err = rollbackConfig()
	}

	if err != nil {
		panic(err)
	}
}

func addConfig() error {
	fmt.Println("Adding config files...")
	return configure("add")
}

func updateConfig() error {
	fmt.Println("Updating config files...")
	return configure("update")
}

func rollbackConfig() error {
	fmt.Println("Rolling back config files...")
	return rollbackConfigs()
}

func configure(action string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Available openbus versions:")
	for i, br := range availableBranches {
		fmt.Printf("%d. %s\n", i+1, br)
	}

	fmt.Print("Enter the number of the openbus version you want to use: ")
	inputBranch, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return err
	}
	inputBranch = strings.TrimSpace(inputBranch)
	branchChoice, err := strconv.Atoi(inputBranch)
	if err != nil || branchChoice < 1 || branchChoice > len(availableBranches) {
		fmt.Println("Invalid choice. Exiting.")
		return err
	}
	branch = availableBranches[branchChoice-1]

	fmt.Println("Please choose the configuration type:")
	fmt.Println("1. with-toolboard")
	fmt.Println("2. without-toolboard")

	var choice int
	fmt.Print("Enter your choice (1 or 2): ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || (choice != 1 && choice != 2) {
		fmt.Println("Invalid choice. Exiting.")
		return err
	}

	subdir := subdirNoTool
	if choice == 1 {
		subdir = subdirTool
	}

	fmt.Println("Do you want to use the smart orbiter or the dyze?")
	fmt.Println("1. Smart Orbiter")
	fmt.Println("2. Dyze")

	fmt.Print("Enter your choice (1 or 2): ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || (choice != 1 && choice != 2) {
		fmt.Println("Invalid choice. Exiting.")
		return err
	}

	selectedFile := "extruder-dyze500.cfg"
	if choice == 1 {
		selectedFile = "extruder-smart-orbiter-v3.cfg"
	}

	fmt.Println("Pulling latest changes from branch", branch)

	// Fetch latest changes
	err = runCommand("git", "fetch")
	if err != nil {
		fmt.Printf("Error fetching: %v\n", err)
		return err
	}

	// Checkout the specified branch
	err = runCommand("git", "checkout", branch)
	if err != nil {
		fmt.Printf("Error checking out branch %s: %v\n", branch, err)
		return err
	}

	// Pull latest changes from origin branch
	err = runCommand("git", "pull", "origin", branch)
	if err != nil {
		fmt.Printf("Error pulling latest changes from branch %s: %v\n", branch, err)
		return err
	}

	fmt.Println("Generating consolidated configuration file...")
	consolidatedContent, err := generateConsolidatedConfig(subdir, selectedFile)
	if err != nil {
		fmt.Printf("Error generating consolidated configuration: %v\n", err)
		return err
	}

	err = runCommand("mkdir", "-p", configDir)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", configDir, err)
		return err
	}

	outputFile := filepath.Join(configPath)
	err = os.WriteFile(outputFile, []byte(consolidatedContent), 0644)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", outputFile, err)
		return err
	} else {
		fmt.Printf("Consolidated configuration written to %s\n", outputFile)

	}

	err = copyMoonrakerConfig(subdir)
	if err != nil {
		fmt.Printf("Error copying moonraker.conf: %v\n", err)
		return err
	}

	fmt.Println("Extracting and replacing serial numbers...")
	err = extractAndReplaceSerialNumbers(configPath)
	if err != nil {
		fmt.Printf("Error extracting and replacing serial numbers: %v\n", err)
		return err
	}

	fmt.Println("Go restart klipper in mainsail for changes to take effect")
	fmt.Println("Configuration update complete!")

	return nil
}

func copyMoonrakerConfig(subdir string) error {
	srcFile := filepath.Join(subdir, "moonraker.conf")
	dstFile := filepath.Join(configDir, "moonraker.conf")

	err := copyFile(srcFile, dstFile)
	if err != nil {
		return fmt.Errorf("error copying moonraker.conf: %v", err)
	}

	fmt.Printf("Successfully copied moonraker.conf to %s\n", dstFile)
	return nil
}