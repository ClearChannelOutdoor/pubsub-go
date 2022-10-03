#!/bin/bash

# Color Constants
NC="\033[0m"
BLACK="\033[0;30m"
RED="\033[0;31m"
GREEN="\033[0;32m"
BROWN="\033[0;33m"
BLUE="\033[0;34m"
PURPLE="\033[0;35m"
CYAN="\033[0;36m"
LIGHT_GREY="\033[0;37m"
DARK_GREY="\033[1;30m"
LIGHT_RED="\033[1;31m"
LIGHT_GREEN="\033[1;32m"
YELLOW="\033[1;33m"
LIGHT_BLUE="\033[1;34m"
LIGHT_PURPLE="\033[1;35m"
LIGHT_CYAN="\033[1;36m"
WHITE="\033[1;37m"

kPORT_REGEX="[0-9]:[0-9]"
kGCP_REGISTRY_URL="https://console.cloud.google.com/gcr/images/google.com:cloudsdktool/GLOBAL/cloud-sdk?gcrImageListsize=30"
kWAIT_TIME=20


while true; do
	read -p "Enter docker port mapping: " ports
	if [[ $ports =~ $kPORT_REGEX ]]; then break; fi
	echo "Port mapping please...\n"
done

while true; do
	read -p "$(echo -e "\nEnter docker image version, or ${RED}type browse${NC} to see a list of available versions:")" version
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

read -p "$(echo -e "\nEnter ${CYAN}Project Name${NC}: ")" project_name
read -p "$(echo -e "\nEnter ${CYAN}Topic Name${NC}: ")" topic_name
read -p "$(echo -e "\nEnter ${CYAN}Subscription Name${NC}: ")" subscription_name

port_array=(${ports//:/ })
localhost_port=${port_array[0]}
host_port="${port_array[1]}"

docker run -d --rm -ti -p $ports \
    gcr.io/google.com/cloudsdktool/cloud-sdk:$version-emulators \
    gcloud beta emulators pubsub start \
      --project=abc \
      --host-port=0.0.0.0:$host_port

echo -e "\nWaiting ${RED}$kWAIT_TIME${NC} seconds for container to start"
sleep $kWAIT_TIME

echo "\nCreating Topic"
curl -X PUT -v "http://localhost:$localhost_port/v1/projects/$project_name/topics/$topic_name"

echo "\nCreating Subscription"
curl -X PUT -H "Content-Type:application/json" -v \
	--data "{\"topic\":\"projects/$project_name/topics/$topic_name\"}"\
	"http://localhost:$localhost_port/v1/projects/$project_name/subscriptions/$subscription_name"

echo -e "${GREEN}Complete${NC}"

echo -e "\nPlease run ${GREEN}export PUBSUB_EMULATOR_HOST=\"localhost:$localhost_port\"${NC} in any open terminal windows"

exit
