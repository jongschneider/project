#!/bin/bash

if [ -z ${SHA+x} ]; then SHA=$(git rev-parse HEAD); fi

docker build -t jongschneider/echo1:latest-dev -f ./echoservers/echo1/Dockerfile.dev ./echoservers/echo1
docker build -t jongschneider/echo1:$SHA -f ./echoservers/echo1/Dockerfile.dev ./echoservers/echo1

docker push jongscnheider/echo1:latest-dev
docker push jongschneider/echo1:$SHA


docker build -t jongschneider/echo2:latest-dev -f ./echoservers/echo2/Dockerfile.dev ./echoservers/echo2
docker build -t jongschneider/echo2:$SHA -f ./echoservers/echo2/Dockerfile.dev ./echoservers/echo2

docker push jongscnheider/echo2:latest-dev
docker push jongschneider/echo2:$SHA
