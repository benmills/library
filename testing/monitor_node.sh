#!/bin/bash

while true; do
  clear
  echo "localhost:$1"
  http localhost:$1/stats
  sleep 5
done
