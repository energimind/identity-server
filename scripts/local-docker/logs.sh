#!/bin/bash

# shellcheck disable=SC2086
# shellcheck disable=SC2048
docker compose -p identity logs $*
