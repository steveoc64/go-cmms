#!/bin/bash
mpstat -P ALL
iostat
cs=`pgrep -d',' cmms-server`
ps=`pgrep -d',' postmaster`
top -bHcs -n 1 -p $cs,$ps
