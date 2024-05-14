#!/bin/bash

PROFILE=${IS_PROFILE:-default}

# shellcheck disable=SC2086
# shellcheck disable=SC2048
docker compose -p identity --profile $PROFILE up $*
