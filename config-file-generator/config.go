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
    configDir    = filepath.Join(homeDir, "printer_data/config")
    subdirTool   = "with-toolboard"
    subdirNoTool = "without-toolboard"
    branch       string
)

var filesToUpdate = []string{
    "aliases-openbus-v02-notoolboard.cfg",
    "bedheat-i400.cfg",
    "extruder-dyze500.cfg",
    "extruder-rapido-orbiter2.cfg",
    "extruder-rapido2-orbiter2.cfg",
    "extruder-smart-orbiter-v3.cfg",
    "kinematics-xy-i400-SmartOrbiter-FixedMount.cfg",
    "kinematics-xy-i400.cfg",
    "kinematics-z-i400-SmartOrbiter-FixedMount.cfg",
    "kinematics-z-i400.cfg",
    "partcooling-5015fans.cfg",
    "partcooling-pump.cfg",
    "print-control.cfg",
    "printer.cfg",
    "probe-heschen-pl-08n.cfg",
    "safe-home.cfg",
    "local.cfg",
    "probe-beacon-revh-SmartOrbiter-FixedMount.cfg",
    "probe-beacon-revh.cfg",
}


func addConfig() {
    fmt.Println("Adding configuration files...")
    configure("addconfig")
}

func updateConfig() {
    fmt.Println("Updating configuration files...")
    configure("updateconfig")
}

func configure(action string) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Available branches:")
    for i, br := range availableBranches {
        fmt.Printf("%d. %s\n", i+1, br)
    }

    fmt.Print("Enter the number of the branch you want to use: ")
    inputBranch, _ := reader.ReadString('\n')
    inputBranch = strings.TrimSpace(inputBranch)
    branchChoice, err := strconv.Atoi(inputBranch)
    if err != nil || branchChoice < 1 || branchChoice > len(availableBranches) {
        fmt.Println("Invalid choice. Exiting.")
        return
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
        return
    }

    subdir := subdirNoTool
    if choice == 1 {
        subdir = subdirTool
    }

    if action == "updateconfig" {
        backupConfigs()
    }

    fmt.Println("Pulling latest changes from branch", branch)

    // Fetch latest changes
    err = runCommand("git", "fetch")
    if err != nil {
        fmt.Printf("Error fetching: %v\n", err)
        return
    }

    // Checkout the specified branch
    err = runCommand("git", "checkout", branch)
    if err != nil {
        fmt.Printf("Error checking out branch %s: %v\n", branch, err)
        return
    }

    // Pull latest changes from origin branch
    err = runCommand("git", "pull", "origin", branch)
    if err != nil {
        fmt.Printf("Error pulling latest changes from branch %s: %v\n", branch, err)
        return
    }

    fmt.Println("Generating consolidated configuration file...")
    consolidatedContent, err := generateConsolidatedConfig(subdir)
    if err != nil {
        fmt.Printf("Error generating consolidated configuration: %v\n", err)
        return
    }

    outputFile := filepath.Join(configDir, "printer.cfg")
    err = os.WriteFile(outputFile, []byte(consolidatedContent), 0644)
    if err != nil {
        fmt.Printf("Error writing to %s: %v\n", outputFile, err)
    } else {
        fmt.Printf("Consolidated configuration written to %s\n", outputFile)
    }

    fmt.Println("Extracting and replacing serial numbers...")
    extractAndReplaceSerialNumbers(filepath.Join(configDir, "local.cfg"), filepath.Join(configDir, "probe-beacon-revh-SmartOrbiter-FixedMount.cfg"))

    fmt.Println("Go restart klipper in mainsail for changes to take effect")
    fmt.Println("Configuration update complete!")
}

func generateConsolidatedConfig(subdir string) (string, error) {
    inputFile := filepath.Join("..", subdir, "printer.cfg")
    content, err := os.ReadFile(inputFile)
    if err != nil {
        return "", fmt.Errorf("error reading printer.cfg: %v", err)
    }

    consolidatedContent := string(content)
    lines := strings.Split(consolidatedContent, "\n")

    for i, line := range lines {
        if strings.HasPrefix(line, "[include ") {
            parts := strings.Split(line, " ")
            if len(parts) == 2 {
                fileName := strings.TrimSuffix(parts[1], "]")
                filePath := filepath.Join("..", subdir, fileName)
                fileContent, err := os.ReadFile(filePath)
                if err != nil {
                    return "", fmt.Errorf("error reading %s: %v", filePath, err)
                }
                lines[i] = string(fileContent)
            }
        }
    }

    return strings.Join(lines, "\n"), nil
}
