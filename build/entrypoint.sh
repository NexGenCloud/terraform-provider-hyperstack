#!/usr/bin/env bash
set -e



echo "Hello $1"
time=$(date)
echo "time=$time" >> $GITHUB_OUTPUT

task --list
echo exec task "$@"
