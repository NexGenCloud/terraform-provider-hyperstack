#!/usr/bin/env bash
set -e



echo "Hello $1"
time=$(date)
#echo "time=$time" >> $GITHUB_OUTPUT

exec "$@"
task --list
echo exec task "$@"
