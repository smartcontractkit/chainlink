#!/usr/bin/env bash

set -e

if [[ "$(basename $(pwd))" = "gethwrappers2" ]];
then
    cd ../../..
fi

SRC=../libocr/gethwrappers2/
DEST=./core/internal/gethwrappers2/

if [[ -d $SRC && -d $DEST ]]
then
    if [[ -d "$DEST/generated" ]]
    then
        rm -rf $DEST/generated.bak
        mv $DEST/generated $DEST/generated.bak
    fi
    mkdir -p $DEST/generated
    if ! cp -a $SRC $DEST/generated
    then
        mv $DEST/generated.bak $DEST/generated
        echo "Failed to copy gethwrappers2 artifacts from libocr."
        exit 1
    fi
    mkdir -p $DEST/generated/offchainaggregator
    if ! cp $DEST/generated.bak/offchainaggregator/offchainaggregator.go $DEST/generated/offchainaggregator/
    then
        mv $DEST/generated.bak $DEST/generated
        echo "Missing $DEST/generated/offchainaggregator.go"
        exit 1
    fi
    rm -rf $DEST/generated.bak
    echo "Generated gethwrappers2 artifacts copied from libocr."
else
    if [[ ! -d $SRC ]]
    then
        echo "Skipping $DEST, due to missing ${SRC}"
    else [[ ! -d $DEST ]]
        echo "Skipping $DEST, due to missing ${DEST}"
    fi
fi

