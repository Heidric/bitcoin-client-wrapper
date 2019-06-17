#!/bin/bash

docker stop test_app
docker rm test_app
docker rmi test_app
docker rmi iron/go

rm ./app/bin/app