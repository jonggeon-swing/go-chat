#!/bin/bash

TARGET_IP=""
TARGET_PORT=8080
MESSAGE="Lorem ipsum dolor sit amet labore dolore diam ut stet accusam lorem kasd stet erat suscipit. Sit erat sea stet ut at nonummy dolor illum kasd eros id wisi labore et tincidunt illum."
COUNT=20000
DELAY=0.001

# socat을 사용해 TCP 연결 유지
{
    for ((i=1; i<=COUNT; i++)); do
        echo "$MESSAGE $i"
        # sleep $DELAY
    done
} | socat - TCP:$TARGET_IP:$TARGET_PORT

