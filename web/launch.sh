#!/bin/bash

hostname=$(hostname)
command="bundle exec ruby app.rb"

if [ $hostname = "raspberrypi" ]; then
	cd /home/pi/tag-mng/web
	command="bundle exec unicorn -c config/unicorn.rb"
fi

echo $command
eval $command

if [ $? != 0 ]; then
	echo "failed to launch"
else
	echo "succeeded"
fi
