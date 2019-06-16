#!/bin/bash

# deploy path
# /home/lbfdeatq/tag-mng
  
# generate db bk file
ssh lbfdeatq@ddddddo.work "cd /home/pi/tag-mng/db/bk; pg_dump tag-mng > `date +%Y%m%d_%H`-dump"

# fetch generated db bk file

# rm generated db bk file

exit 0

