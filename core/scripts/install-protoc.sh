#!/bin/bash
set -x


VERSION=$1

if [ "$VERSION" == "" ]; then
    echo "version required"
    exit 1
fi

os=$(uname)
arch=$(uname -m)

install_dir=$HOME/.local
$install_dir/bin/protoc --version | grep $VERSION
rc=$?
if [ $rc -eq 0 ]; then
    # we have the current VERSION
    echo "protoc up-to-date @ $VERSION"
    exit 0
fi


if [ "$os" == "Linux" ] ; then
	os="linux"
    if [ "$arch" != "x86_64" ]; then
        echo "unsupported os $os-$arch update $0"
        exit 1
    fi
elif [ "$os" == "Darwin" ] ; then
	os="osx"
    # make life simply and download the universal binary
    arch="universal_binary"
else
    echo "unsupported os $os. update $0"
    exit 1
fi

workdir=$(mktemp -d)
pushd $workdir
pb_url="https://github.com/protocolbuffers/protobuf/releases"
artifact=protoc-$VERSION-$os-$arch.zip
curl -LO $pb_url/download/v${VERSION}/$artifact
if [[ ! -d $install_dir ]]; then
    mkdir $install_dir
fi
unzip -o $artifact -d $install_dir
rm $artifact

echo "protoc $VERSION installed in $install_dir"
echo "Add $install_dir/bin to PATH"
export PATH=$install_dir/bin:$PATH
popd
