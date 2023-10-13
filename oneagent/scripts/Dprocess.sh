#!/bin/bash

# D process monitor
Dprocess=` ps -eL -o lwp,pid,ppid,state,comm | grep -E " D " |wc -l `
echo "a1 $Dprocess"
