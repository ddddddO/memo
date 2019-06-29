#!/bin/bash

/home/pi/health-cheack/health-cheack/health-cheack

if [ $? != 0 ]; then
	reboot
fi

