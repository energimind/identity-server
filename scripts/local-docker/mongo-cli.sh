#!/bin/bash

user=mongo
pass="-p mongo"

# shellcheck disable=SC2086
docker compose -p identity exec -it mongo mongosh -u $user $pass
