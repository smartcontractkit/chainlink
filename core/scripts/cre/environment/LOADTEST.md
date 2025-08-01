# === Variables ===
GITHUB_TOKEN ?= your_token_here
CL_IMAGE_TAG ?= chainlink-node:plugins
CAPABILITIES_REPO ?= ../capabilities
WORKFLOW_NAME ?= fetchtrueusd
WORKFLOW_ID ?= 0089c0071e8c5b535ebeab3f5102091f3a657daf2bb1778eea67f4a12b82c2cb

# === Docker Build ===
docker-image:
CL_INSTALL_PRIVATE_PLUGINS=true CL_INSTALL_TESTING_PLUGINS=true GITHUB_TOKEN=$(GITHUB_TOKEN) make docker-plugins

# === Environment Setup ===
env-start:
CTF_CONFIGS=./configs/workflow-load.toml \
go run main.go env start --topology=mock --extra-allowed-gateway-ports=16000 --with-beholder

# === Workflow Binary Build ===
build-workflow:
cd $(CAPABILITIES_REPO)/workflows/$(WORKFLOW_NAME)/cmd && \
GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0 go build -o $(WORKFLOW_NAME) && \
brotli -v $(WORKFLOW_NAME) && \
cat $(WORKFLOW_NAME).br | base64 > $(WORKFLOW_NAME).br.base64

# === Prepare and Upload Assets ===
prepare-assets:
cp $(CAPABILITIES_REPO)/workflows/$(WORKFLOW_NAME)/cmd/$(WORKFLOW_NAME).br.base64 $(WORKFLOW_NAME).br
touch empty.yaml

upload-assets:
go run main.go minio upload $(WORKFLOW_NAME).br empty.yaml

# === Workflow Registration ===
register-workflow:
go run main.go workflow register \
--binary-url="http://minio:16000/default/$(WORKFLOW_NAME).br" \
--config-url="http://minio:16000/default/empty.yaml" \
--secrets-url="http://minio:16000/default/empty.yaml" \
--id="$(WORKFLOW_ID)" \
--name="$(WORKFLOW_NAME)"

# === Capability Registration ===
register-capabilities:
go run main.go registry create --name="cron-trigger" --version="1.0.0" --type="target" --don-id=2
go run main.go registry create --name="write_ethereum-testnet-sepolia" --version="1.0.0" --type="target" --don-id=2

create-mocks:
go run main.go mock create --id="write_ethereum-testnet-sepolia@1.0.0" --description="mock target" --type="target" --addresses="192.168.48.13:7777,192.168.48.14:7777,192.168.48.15:7777"
go run main.go mock create --id="cron-trigger@1.0.0" --description="mock trigger" --type="trigger" --addresses="192.168.48.13:7777,192.168.48.14:7777,192.168.48.15:7777"

