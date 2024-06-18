package main

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "runtime"
    "strings"
)

func extractAndReplaceSerialNumbers(localConfigPath, probeConfigPath string) {
    fmt.Println("Extracting serial numbers...")

    beaconSerial := getSerial("Beacon")
    if beaconSerial == "" {
        fmt.Println("Failed to extract beacon serial number. Exiting.")
        return
    }

    mcuSerial := getSerial("Klipper")
    if mcuSerial == "" {
        fmt.Println("Failed to extract MCU serial number. Exiting.")
        return
    }

    fmt.Println("Beacon Serial:", beaconSerial)
    fmt.Println("MCU Serial:", mcuSerial)

    fmt.Println("Updating local.cfg...")
    replaceSerialInFile(localConfigPath, "serial:", mcuSerial)
    fmt.Println("Updating probe-beacon-revh-SmartOrbiter-FixedMount.cfg...")
    replaceSerialInFile(probeConfigPath, "serial:", beaconSerial)
}

func getSerial(identifier string) string {
    if runtime.GOOS == "darwin" {
        // Return fake serial numbers for macOS
        if identifier == "Beacon" {
            return "/dev/serial/by-id/usb-Beacon_Beacon_RevH_FAKE1234567890-if00"
        } else if identifier == "Klipper" {
            return "/dev/serial/by-id/usb-Klipper_stm32h723xx_FAKE1234567890-if00"
        }
    }

    cmd := exec.Command("bash", "-c", fmt.Sprintf("ls /dev/serial/by-id/* | grep %s", identifier))
    output, err := cmd.Output()
    if err != nil {
        fmt.Printf("Error extracting %s serial number: %v\n", identifier, err)
        return ""
    }
    return strings.TrimSpace(string(output))
}

func replaceSerialInFile(filePath, searchPattern, serialNumber string) {
    input, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Printf("Error reading file %s: %v\n", filePath, err)
        return
    }

    lines := strings.Split(string(input), "\n")
    for i, line := range lines {
        if strings.HasPrefix(line, searchPattern) {
            parts := strings.SplitN(line, " ", 2)
            if len(parts) == 2 {
                lines[i] = fmt.Sprintf("%s %s", searchPattern, serialNumber)
            }
        }
    }

    output := strings.Join(lines, "\n")
    err = ioutil.WriteFile(filePath, []byte(output), 0644)
    if err != nil {
        fmt.Printf("Error writing to file %s: %v\n", filePath, err)
    }
}

