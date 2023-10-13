#!/bin/bash
zombie=$(ps -ostat |grep -e '^[Zz]' |wc -l)
echo "a1 $zombie"
