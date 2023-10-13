#!/bin/bash
Load=` uptime | tr -s ' ' | cut -d ' ' -f 11 | cut -d ',' -f 1 `
echo "a1 $Load"
