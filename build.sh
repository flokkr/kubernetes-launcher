#!/usr/bin/env bash

IMAGE=flokkr/kubernetes-launcher
###
# Create the docker image
###
build(){
   docker build -t $IMAGE .
}

###
# Deploy the docker image to the docker hub
###
deploy(){
   docker push $IMAGE
}

$@

