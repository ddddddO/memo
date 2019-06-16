#!/bin/bash

# deploy path
# /home/lbfdeatq/tag-mng

# scheduled
# 13 10,16,22,4 * * * 

function backup () {
        echo "generate db bk file"
        ssh lbfdeatq@ddddddo.work "cd /home/pi/tag-mng/db/bk; pg_dump tag-mng > `date +%Y%m%d_%H`-dump"


        echo "fetch generated db bk file"
        scp lbfdeatq@ddddddo.work:/home/pi/tag-mng/db/bk/*-dump /home/lbfdeatq/tag-mng/db-bk

        return
}

echo "start db bk"

backup

if [ $? != 0 ]; then
        echo "restart db bk after 18min"
        sleep 1080
        backup
fi

if [ $? != 0 ]; then
        echo "failed db bk(generate bk file OR fetch db bk file)"
        exit 1
fi

# rm generated db bk file
ssh lbfdeatq@ddddddo.work "cd /home/pi/tag-mng/db/bk; rm /home/pi/tag-mng/db/bk/*-dump"

if [ $? != 0 ]; then
        echo "failed to rm generated db bk file"
        exit 1
fi

echo "succeeded db bk!"
exit 0
