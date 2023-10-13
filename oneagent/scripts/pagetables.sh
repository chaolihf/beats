#!/bin/bash
pagetables=$(awk '/PageTables/{print$2}' /proc/meminfo)
echo "a1 $pagetables"
