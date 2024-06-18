package main

import (
    "fmt"
    "os"
    "path/filepath"
)

var backupDir = filepath.Join(configDir, "backup")

func backupConfigs() {
    fmt.Println("Creating backup of existing configuration files...")
    os.MkdirAll(backupDir, os.ModePerm)
    for _, file := range filesToUpdate {
        src := filepath.Join(configDir, file)
        dest := filepath.Join(backupDir, file)
        if _, err := os.Stat(src); err == nil {
            copyFile(src, dest)
        } else {
            fmt.Printf("%s does not exist in the config directory.\n", file)
        }
    }
}

func rollbackConfigs() {
    fmt.Println("Restoring configuration files from backup...")
    if _, err := os.Stat(backupDir); os.IsNotExist(err) {
        fmt.Printf("Backup directory %s does not exist or is empty. Rollback failed.\n", backupDir)
        return
    }
    for _, file := range filesToUpdate {
        src := filepath.Join(backupDir, file)
        dest := filepath.Join(configDir, file)
        if _, err := os.Stat(src); err == nil {
            copyFile(src, dest)
        } else {
            fmt.Printf("%s does not exist in the backup directory.\n", file)
        }
    }
    fmt.Println("Rollback complete! Go restart klipper from mainsail")
}

