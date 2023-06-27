#!/usr/bin/env bash

apk update
apk add curl jq

celestia bridge init --node.store /bridge
export CELESTIA_NODE_AUTH_TOKEN=$(celestia bridge auth admin --node.store /bridge)
echo "WARNING: Keep this auth token secret **DO NOT** log this auth token outside of development. CELESTIA_NODE_AUTH_TOKEN=$CELESTIA_NODE_AUTH_TOKEN"
/wait-for-it.sh 192.167.10.10:26657 -t 90 -- \
  curl -s http://192.167.10.10:26657/block?height=1 | jq '.result.block_id.hash' | tr -d '"' > genesis.hash

curl -s http://192.167.10.10:26657/block_by_hash?hash=0x`cat genesis.hash`
echo  # newline

export CELESTIA_CUSTOM=test:`cat genesis.hash`
echo $CELESTIA_CUSTOM
celestia bridge start \
  --node.store /bridge --gateway --gateway.deprecated-endpoints \
  --core.ip 192.167.10.10 \
  --keyring.accname validator
