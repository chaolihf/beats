#!/bin/bash
Sync=` sar -B 1 1 |awk '/pgscank/ {getline; print}'|grep -v Aver|awk '{print $9}'|grep -v ^0.00|wc -l `
echo "a1 $Sync"
