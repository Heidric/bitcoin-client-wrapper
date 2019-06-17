#!/bin/bash

if [ -z "$MAIN_PORT" ]; then
    export MAIN_PORT="7777"
fi

if [ -z "$ENV" ]; then
    export ENV="test"
fi

docker-compose build && docker-compose up -d