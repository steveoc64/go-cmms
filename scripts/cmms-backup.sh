#!/bin/bash
exit

date=`date "+%Y%m%d-%H%M%S"`
file="cmms-$date"
echo $file
pg_dump -U postgres cmms | gzip > ../backup/$file.sqz
cd ../backup
echo Latest Backup File
ls -l $file.sqz
echo
echo All Backups
ls -s
echo 
echo Files older than 14 days
find . -mtime +14

