#!/bin/bash

/home/pi/health-cheack/health-cheack/health-cheack

if [ $? != 0 ]; then
	sudo reboot
fi

