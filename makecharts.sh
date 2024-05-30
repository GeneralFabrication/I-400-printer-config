#!/bin/bash
echo calibrating for X axis:
~/klipper/scripts/calibrate_shaper.py /tmp/calibration_data_x_*.csv -o ~/printer_data/config/charts/shaper_calibrate_x_$(date "+%F_%H:%M").png
echo
echo calibrating for y axis:
~/klipper/scripts/calibrate_shaper.py /tmp/calibration_data_y_*.csv -o ~/printer_data/config/charts/shaper_calibrate_y_$(date "+%F_%H:%M").png
rm /tmp/calibration_data_x_*.csv /tmp/calibration_data_y_*.csv
