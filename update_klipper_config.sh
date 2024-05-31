#!/bin/bash

CONFIG_DIR=~/printer_data/config
BACKUP_DIR=$CONFIG_DIR/backup
BRANCH="v0.2"
SUBDIR="without-toolboard"

echo "Please choose the configuration type:"
echo "1. with-toolboard"
echo "2. without-toolboard"
read -p "Enter your choice (1 or 2): " choice

if [ "$choice" -eq 1 ]; then
    SUBDIR="with-toolboard"
elif [ "$choice" -eq 2 ]; then
    SUBDIR="without-toolboard"
else
    echo "Invalid choice. Exiting."
    exit 1
fi

FILES_TO_UPDATE=(
    "aliases-openbus-v02-notoolboard.cfg"
    "bedheat-i400.cfg"
    "extruder-dyze500.cfg"
    "extruder-rapido-orbiter2.cfg"
    "extruder-rapido2-orbiter2.cfg"
    "extruder-smart-orbiter-v3.cfg"
    "kinematics-xy-i400-SmartOrbiter-FixedMount.cfg"
    "kinematics-xy-i400.cfg"
    "kinematics-z-i400-SmartOrbiter-FixedMount.cfg"
    "kinematics-z-i400.cfg"
    "partcooling-5015fans.cfg"
    "partcooling-pump.cfg"
    "print-control.cfg"
    "printer.cfg"
    "probe-heschen-pl-08n.cfg"
    "safe-home.cfg"
)

function print_msg {
    echo "---------------------------------------------------"
    echo $1
    echo "---------------------------------------------------"
}

function backup_configs {
    print_msg "Creating backup of existing configuration files..."
    mkdir -p $BACKUP_DIR
    for file in "${FILES_TO_UPDATE[@]}"; do
        if [ -f "$CONFIG_DIR/$file" ]; then
            cp "$CONFIG_DIR/$file" "$BACKUP_DIR/"
        else
            print_msg "$file does not exist in the config directory."
        fi
    done
}

function rollback_configs {
    print_msg "Restoring configuration files from backup..."
    if [ -d "$BACKUP_DIR" ] && [ "$(ls -A $BACKUP_DIR)" ]; then
        for file in "${FILES_TO_UPDATE[@]}"; do
            if [ -f "$BACKUP_DIR/$file" ]; then
                cp "$BACKUP_DIR/$file" "$CONFIG_DIR/"
            else
                print_msg "$file does not exist in the backup directory."
            fi
        done
        print_msg "Rollback complete! Go restart klipper from mainsail"
    else
        print_msg "Backup directory $BACKUP_DIR does not exist or is empty. Rollback failed."
    fi
}

function update_configs {
    print_msg "Pulling latest changes from branch $BRANCH..."
    git fetch
    git checkout $BRANCH
    git pull origin $BRANCH

    print_msg "Copying new configuration files..."
    for file in "${FILES_TO_UPDATE[@]}"; do
        if [ -f "$SUBDIR/$file" ]; then
            cp "$SUBDIR/$file" "$CONFIG_DIR/"
        else
            print_msg "$file does not exist in the repository subdirectory."
        fi
    done

    print_msg "Go restart klipper in mainsail for changes to take effect"

    print_msg "Configuration update complete!"
}

if [ "$1" == "rollback" ]; then
    rollback_configs
else
    backup_configs
    update_configs
fi
