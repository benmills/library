#!/bin/bash
base_port=$1
node_count=$2

node=$(( $node_count - 1))
while [ $node -gt 0 ]; do
  curl -X PUT --data "http://localhost:$(( $node + $base_port ))" localhost:$base_port/peers/join
  node=$(( $node - 1 ))
done
