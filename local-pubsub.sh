#!/bin/bash

kPORT_REGEX="[0-9]:[0-9]"
kGCP_REGISTRY_URL="https://console.cloud.google.com/gcr/images/google.com:cloudsdktool/GLOBAL/cloud-sdk?gcrImageListsize=30"

while true; do
	read -p "Enter docker port mapping (\"3080:3085\" is default):" ports
	if [[ $ports =~ $kPORT_REGEX ]]; then break; fi
	echo "Port mapping please..."
done

while true; do
	read -p "Enter docker image version, or type browse to see a list of available versions:" version
	if [[ "$version" != "browse" ]]; then break; fi
	if [[ "$version" == "browse" ]]; then
		if which xdg-open >/dev/null; then
			xdg-open "$kGCP_REGISTRY_URL"
		elif which open >/dev/null; then
			open "$kGCP_REGISTRY_URL"
		else
			echo "nothing to open with"
		fi
	fi
done

#read -p "Enter project name:" project_name
#read -p "Enter topic name:" topic_name
#read -p "Enter subscription name:" subscription_name

port_array=(${ports//:/ })
host_port="${port_array[1]}"

docker run -d --rm -ti -p $ports \
    gcr.io/google.com/cloudsdktool/cloud-sdk:$version-emulators \
    gcloud beta emulators pubsub start \
      --project=abc \
      --host-port=0.0.0.0:$host_port

echo "end"
