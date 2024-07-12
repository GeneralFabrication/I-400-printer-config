#I-400 Printer Configuration
This repository contains configuration files for the I-400 3D printer running Klipper firmware. The configurations are provided as .cfg Go template files and are available for use through an executable on the printers.

#Directory Structure
/v0.2: Contains two subdirectories for different configurations:
/v0.2/with-toolboard: Configuration files for a 3D printer with a toolboard attached. Includes settings for:
Extruder
Bed heating
Kinematics
Probe
Part cooling
Print control
TMC2208 and TMC2209 stepper drivers
Temperature sensors
Bed mesh generation
Board-RPI serial communication
Gcode macros for homing behavior
/v0.2/without-toolboard: Configuration files for a 3D printer without a toolboard. Includes settings for:
Bed heater
Extruder
Part cooling fan
XY kinematics
Z-axis bed motion
Probe
General print control
Safe homing
Print cancellation
TMC2160 and TMC2209 drivers
Heaters, fans, thermistors, and endstops
#Ignored Files
The .gitignore file specifies the following files to be ignored:

crowsnest.conf
.moonraker.conf.bkp
timelapse.cfg
sonar.conf
mainsail.cfg
