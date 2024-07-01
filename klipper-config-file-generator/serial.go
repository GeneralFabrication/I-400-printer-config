package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strings"
)

func extractAndReplaceSerialNumbers(configPath string) error {
    fmt.Println("Extracting serial numbers...")

    beaconSerial, err := getSerial("Beacon")
    if err != nil {
        fmt.Println("Failed to extract beacon serial number. Exiting.")
        return err
    }

    mcuSerial, err := getSerial("Klipper")
    if err != nil {
        fmt.Println("Failed to extract MCU serial number. Exiting.")
        return err
    }

    fmt.Println("Beacon Serial:", beaconSerial)
    fmt.Println("MCU Serial:", mcuSerial)

    fmt.Println("Updating mcu serial")
    err = replaceSerialInFile(configPath, "[mcu]", "serial:", mcuSerial)
    if err != nil {
        fmt.Println("Failed to update MCU serial number. Exiting.")
        return err
    }

    fmt.Println("Updating beacon serial")
    err = replaceSerialInFile(configPath, "[beacon]", "serial:", beaconSerial)
    if err != nil {
        fmt.Println("Failed to update beacon serial number. Exiting.")
        return err
    }

    return nil
}

func getSerial(identifier string) (string, error) {
    if runtime.GOOS == "darwin" {
        // Return fake serial numbers for macOS
        if identifier == "Beacon" {
            return "/dev/serial/by-id/usb-Beacon_Beacon_RevH_FAKE1234567890-if00", nil
        } else if identifier == "Klipper" {
            return "/dev/serial/by-id/usb-Klipper_stm32h723xx_FAKE1234567890-if00", nil
        }
    }

    cmd := exec.Command("bash", "-c", fmt.Sprintf("ls /dev/serial/by-id/* | grep %s", identifier))
    output, err := cmd.Output()
    if err != nil {
        fmt.Printf("Error extracting %s serial number: %v\n", identifier, err)
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}

func replaceSerialInFile(filePath, section, key, newSerial string) error {
    input, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return err
    }
    defer input.Close()

    output, err := os.Create(filePath + ".tmp")
    if err != nil {
        fmt.Println("Error creating temp file:", err)
        return err
    }
    defer output.Close()

    scanner := bufio.NewScanner(input)
    writer := bufio.NewWriter(output)

    inTargetSection := false

    for scanner.Scan() {
        line := scanner.Text()

        if strings.TrimSpace(line) == section {
            inTargetSection = true
        } else if strings.HasPrefix(strings.TrimSpace(line), "[") && strings.HasSuffix(strings.TrimSpace(line), "]") {
            inTargetSection = false
        }

        if inTargetSection && strings.HasPrefix(strings.TrimSpace(line), key) {
            fmt.Fprintln(writer, key+" "+newSerial)
        } else {
            fmt.Fprintln(writer, line)
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return err
    }

    writer.Flush()
    err = os.Rename(filePath+".tmp", filePath)
    return err
}
