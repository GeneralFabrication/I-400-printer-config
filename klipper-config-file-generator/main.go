package main

import (
    "fmt"
    "os"
)

var (
    availableBranches = []string{"v0.1", "v0.2"}
)

func main() {
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "addconfig":
            addConfig()
        case "updateconfig":
            updateConfig()
        case "rollbackconfig":
            rollbackConfigs()
        default:
            fmt.Println("Invalid command. Use addconfig, updateconfig, or rollbackconfig.")
        }
    } else {
        fmt.Println("Please provide a command: addconfig, updateconfig, or rollbackconfig.")
    }
}

