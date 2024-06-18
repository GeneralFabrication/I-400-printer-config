package main

import (
    "fmt"
    "runtime"
    "strings"
    "os"
    "bufio"
    "os/exec"
)

func extractAndReplaceSerialNumbers(configPath string) {
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

    fmt.Println("Updating mcu serial")
    replaceSerialInFile(configPath, "[mcu]", "serial:", mcuSerial)
    fmt.Println("Updating beacon serial")
    replaceSerialInFile(configPath, "[beacon]", "serial:", beaconSerial)
}

func getSerial(identifier string) string {
    if runtime.GOOS == "darwin" {
        // Return fake serial numbers for macOS
        if identifier == "Beacon" {
            return "/dev/serial/by-id/usb-Beacon_Beacon_RevH_FAKE1234567890-if00"
        } else if identifier == "MCU" {
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

func replaceSerialInFile(filePath, section, key, newSerial string) {
	input, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer input.Close()

	output, err := os.Create(filePath + ".tmp")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
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
		return
	}

	writer.Flush()
	os.Rename(filePath+".tmp", filePath)
}

