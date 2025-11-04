#!/bin/bash
set -euox pipefail

export CL_DIR=$(dirname $0)
export SUI_DIR=$(realpath $CL_DIR/../chainlink-sui)

export TEST_LOG_LEVEL=debug
export LOG_LEVEL=debug
export SETH_LOG_LEVEL=debug
export TEST_SETH_LOG_LEVEL=debug
export CL_DATABASE_URL="postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable"
export CL_SUI_CMD="$SUI_DIR/chainlink-sui"

rm -vf $HOME/ram/sui.txt $HOME/ram/sui-log.txt $HOME/ram/loop_*

pushd $CL_DIR/integration-tests/smoke/ccip/logs/
rm -vf *.log
popd

pushd $SUI_DIR
rm -vf chainlink-sui
go build -o chainlink-sui ./relayer/cmd/chainlink-sui/main.go
popd

cd integration-tests/smoke/ccip
exec go test -v -tags=integration -count=1 -run Test_CCIP_Upgrade_EVM2Sui ./... -timeout=20m
# if [ "${1:-}" = "dest" ]; then
#   exec go test -v -tags=integration -count=1 -run Test_CCIPMessaging_EVM2Sui ./...
# else
#   exec go test -v -tags=integration -count=1 -run Test_CCIPMessaging_Sui2EVM ./...
# fi
