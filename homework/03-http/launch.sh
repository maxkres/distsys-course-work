#!/usr/bin/env bash

# The following commands runs the server with the provided CLI arguments.
# $0 -- path to the script
# $@ -- substitute aguments passed to the script
 ../server/server "$@"

# For Go you will have something like this:
# go run $( dirname -- "$0"; )/server/*.go -- "$@"
