#!/bin/bash
APP_NAME="public-re-streamer-server"
id=$(sudo docker restart $APP_NAME)
sudo docker logs -f $id