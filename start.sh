#!/bin/sh
set -e

/server &

nginx -g "daemon off;"
