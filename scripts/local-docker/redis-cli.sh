#!/bin/bash

docker compose -p identity exec -it redis redis-cli -a redis
