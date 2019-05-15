#!/bin/bash

# piユーザーで実行しないこと(実行不可)

hostname=$(hostname)
command="bundle exec ruby app.rb"

if [ $hostname = "raspberrypi" ]; then
	command="bundle exec unicorn app.rb -c config/unicorn.rb"
fi

eval $command

if [ $? != 0 ]; then
	echo "failed to launch"
else
	echo "succeeded"
fi
