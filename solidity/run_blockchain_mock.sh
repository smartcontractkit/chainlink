#!/bin/bash

this_dir=`dirname $0`
$this_dir/node_modules/.bin/ganache-cli \
    -m 'candy maple cake sugar pudding cream honey rich smooth crumble sweet treat' \
    -p 18545 "$@"
