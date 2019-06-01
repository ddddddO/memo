#!/bin/bash

hostname=$(hostname)
command="bundle exec ruby app.rb"

if [ $hostname = "raspberrypi" ]; then
	sleep 1080
	command="bundle exec unicorn -c config/unicorn.rb"
fi

echo $command
eval $command

if [ $? != 0 ]; then
	echo "failed to launch"
else
	echo "succeeded"
fi
