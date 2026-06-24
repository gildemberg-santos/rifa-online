#!/bin/sh
set -e

export PORT=8080

/server &

nginx -g "daemon off;"
