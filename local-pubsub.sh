#!/bin/bash

kPORT_REGEX="[0-9]:[0-9]"

while true; do
	read -p "Enter docker port mapping (\"3080:3085\" is default):" ports
	if [[ $ports =~ $kPORT_REGEX ]]; then break; fi
	echo "Port mapping please..."
done

read -p "Enter project name:" project_name
read -p "Enter topic name:" topic_name
read -p "Enter subscription name:" subscription_name

port_array=(${ports//:/ })
echo ${port_array[0]}
echo ${port_array[1]}
echo "end"
