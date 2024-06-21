package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var klipperConfigFiles = []string{
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
	"local.cfg",
	"mainsail.cfg",
	"partcooling-5015fans.cfg",
	"partcooling-pump.cfg",
	"print-control.cfg",
	"printer.cfg",
	"probe-beacon-revh-SmartOrbiter-FixedMount.cfg",
	"probe-beacon-revh.cfg",
	"probe-heschen-pl-08n.cfg",
	"safe-home.cfg",
}
var backupDir = filepath.Join(configDir, "backup")

func backupConfigs() error {
	fmt.Println("Creating backup of existing configuration files...")
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating backup directory: %v\n", err)
		return err
	}

	for _, file := range klipperConfigFiles {
		src := filepath.Join(configDir, file)
		dest := filepath.Join(backupDir, fmt.Sprintf("%s_%s", file, time.Now().Format("20060102150405")))
		_, err := os.Stat(src)
		if err != nil {
			fmt.Printf("%s does not exist in the config directory.\n", file)
			return err
		}
		err = copyFile(src, dest)
		if err != nil {
			fmt.Printf("%s does not exist in the config directory.\n", file)
		}
	}
	return nil
}

func rollbackConfigs() error {
	fmt.Println("Restoring configuration files from backup...")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		fmt.Printf("Backup directory %s does not exist or is empty. Rollback failed.\n", backupDir)
		return err
	}

	// List all backup files
	files, err := os.ReadDir(backupDir)
	if err != nil {
		fmt.Printf("Error reading backup directory: %v\n", err)
		return err
	}

	// Find the latest backup for each file
	latestBackups := make(map[string]string)
	for _, file := range files {
		for _, originalFile := range klipperConfigFiles {
			if strings.HasPrefix(file.Name(), originalFile+"_") {
				if latestBackups[originalFile] == "" || file.Name() > latestBackups[originalFile] {
					latestBackups[originalFile] = file.Name()
				}
			}
		}
	}

	for originalFile, backupFile := range latestBackups {
		src := filepath.Join(backupDir, backupFile)
		dest := filepath.Join(configDir, originalFile)
		_, err := os.Stat(src)
		if err != nil {
			fmt.Printf("%s does not exist in the backup directory.\n", backupFile)
			return err
		}

		err = copyFile(src, dest)
		if err != nil {
			fmt.Printf("Error restoring %s from backup: %v\n", originalFile, err)
			return err
		}
	}

	fmt.Println("Rollback complete! Go restart klipper from mainsail")
	return nil
}
