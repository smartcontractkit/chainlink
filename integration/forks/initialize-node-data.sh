#!/bin/bash

# The purpose of this script is to generate the initial datadirs for ethereum
# nodes used in the forking test. It is not intended to be used in automation.
# You should find some datadirs already created in ./gethnet, corresponding node
# configs in geth-config-*.toml, and corresponding nodes in docker-compose.yaml.
# If you want to create more nodes from scratch, you will need to read this
# script and follow its instructions, to bring those four dependencies into line
# with each other.
#
# `initialize-node-data.sh <num>` expects geth-config-<num>.toml to be present
# in this directory. It should be copied from geth-config-raw.toml, and
# modified. In particular, if you are using multiple nodes, you will need to
# adjust the DataDir field. You probably also want to change the WSHost to
# whatever IP address you assign this to in docker.
#
# It attempts to delete the existing data directory, assuming that's of the form
# ./gethnet/datadir<num>. It then initializes a new chain for that directory.
#
# It creates a new account, to ensure that the node has a coinbase address to
# mine to. The account's credentials are otherwise unnecessary for this test.
#
# It prints out the enode address for the resulting geth node. Copy that into
# the StaticNodes list of the other nodes you want it to participate with. Be
# sure to edit the IP address it contains, to match the IP address you assigned
# to the node in docker-compose.yaml.
#
# If it doesn't print out something of the form
# enode://<hex>@<ip-address>:<port>, the script has failed. Append an "x" to
# "set -e", to see where the failure is happening.

set -e

# Geth arguments
datadir=./gethnet/datadir$1
config="--config /root/geth-config-$1.toml"

# Docker arguments
image=ethereum/client-go
dirmapping="-v `pwd`:/root"
transient_container="--rm"
invoke_geth="docker run -i $dirmapping $transient_container $image $config"

function initialize_new_chain() {
    $invoke_geth init /root/zero-genesis.json > /dev/null 2>&1
}

# Regexp for  enode://<hex   addr>@<ip       address>:<listening port>
enode_regexp='enode://[0-9a-fA-F]+@([[:digit:]]+\.?)+:[[:digit:]]+'

function get_enode_address()  {
    ( $invoke_geth console < /dev/null 2>&1 ) | egrep -o $enode_regexp
}

sudo rm -rf $datadir # Delete the backing image of the datadir, if needed

initialize_new_chain

enode_address=`get_enode_address`

echo "Copy the following into the the StaticNodes list of the config files for"
echo "the nodes you want this node to participate with. Be sure to change the"
echo "IP address to the assigned to the node in docker-compose.yaml."
echo $enode_address

# Copy the keystore for the devnet key into this datadir
keystore_dir=$datadir/keystore/
mkdir -p $keystore_dir
sudo cp ../../tools/clroot/keys/UTC* $keystore_dir
