#!/bin/bash

if [ $# -eq 0 ]; then
    >&2 echo "No arguments provided"
    exit 1
fi


for d in */ ; do
    git tag -f "${d}$1"
    git push origin -f "${d}$1"
    git tag -f "${d}latest"
    git push origin -f "${d}latest"
done
