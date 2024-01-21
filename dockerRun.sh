#!/bin/bash
APP_NAME="public-re-streamer-server"
sudo docker rm -f $APP_NAME || echo ""
id=$(sudo docker run -dit \
--name $APP_NAME \
--restart="always" \
--network=6105-buttons-1 \
--mount type=bind,source="$(pwd)"/config.yaml,target=/home/morphs/config.yaml \
--mount type=bind,source="$(pwd)"/cookies.txt,target=/home/morphs/cookies.txt \
-p 5955:5955 \
$APP_NAME /home/morphs/config.yaml)
sudo docker logs -f $id