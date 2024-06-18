package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generateConsolidatedConfig(subdir string, selectedFile string) (string, error) {
	inputFile := filepath.Join(subdir, "printer.cfg")
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
				filePath := filepath.Join(subdir, fileName)

				// Check if the included file is one of the extruder files
				if fileName == "extruder-dyze500.cfg" || fileName == "extruder-smart-orbiter-v3.cfg" {
					filePath = filepath.Join(subdir, selectedFile)
				}

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
