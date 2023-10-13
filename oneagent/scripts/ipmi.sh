#!/bin/bash
ipmitool sel list >/root/ipmi-new.log
cmp --silent /root/ipmi.log  /root/ipmi-new.log && echo "a1 0" ||echo "a1 1"
mv -f /root/ipmi-new.log /root/ipmi.log
