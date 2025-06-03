#!/usr/bin/env bash

#set -o pipefail

GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color
CHECK_MARK="${GREEN}✓${NC}"
CROSS_MARK="${RED}✗${NC}"
WARNING_MARK="${YELLOW}!${NC}"

JD_ECR="804282218731.dkr.ecr.us-west-2.amazonaws.com"
JD_VERSION="0.12.7"
JD_LOCAL_IMAGE="job-distributor:${JD_VERSION}"
JD_ECR_IMAGE="$JD_ECR/$JD_LOCAL_IMAGE"

AWS_PROFILE="sdlc"

echo "🔍 Checking prerequisites for CRE environment..."

# Function to check if command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to check if Docker setting is enabled
check_docker_setting() {
  local settings_file=$1
  local setting=$2
  local expected=$3
  local value=$(jq -r ".$setting" "$settings_file" 2>/dev/null)
  
  if [ "$value" = "$expected" ]; then
    echo -e "  ${CHECK_MARK} $setting is correctly set to $expected"
    return 0
  else
    echo -e "  ${CROSS_MARK} $setting is set to $value (should be $expected)"
    return 1
fi
}

missing_docker_setting() {
    local setting=$1
    echo -e "  ${CROSS_MARK} $setting is missing from Docker settings file"
    DOCKER_SETTINGS_OK=false # global variable to track Docker settings status
}

build_jd_image() {
    echo "Building Job Distributor image..."
    
    # Create a temporary directory for cloning
    temp_dir=$(mktemp -d)
    
    # Clone the repository
    git clone https://github.com/smartcontractkit/job-distributor "$temp_dir"
    cd "$temp_dir"
    git checkout v0.9.0
    
    # Build the Docker image
    docker build -t $JD_LOCAL_IMAGE -f e2e/Dockerfile.e2e .
    
    # Clean up
    cd - > /dev/null
    rm -rf "$temp_dir"
    
    echo -e "${CHECK_MARK} Job Distributor image built successfully"
}

pull_jd_image() {
    aws configure list-profiles | grep -q "$AWS_PROFILE"
    if [ $? -ne 0 ]; then
        echo -e "${CROSS_MARK} AWS profile '$AWS_PROFILE' not found"
        echo -e "Please ensure you have the correct AWS profile configured"
        exit 1
    fi
    echo "AWS SSO Login for profile $AWS_PROFILE..."
    if ! aws sso login --profile "$AWS_PROFILE"; then
        echo -e "${CROSS_MARK} AWS SSO login failed for profile $AWS_PROFILE"
        echo -e "Please ensure you have the correct AWS profile configured and try again"
        exit 1
    fi
    echo -e "${CHECK_MARK} AWS SSO login successful for profile $AWS_PROFILE"
    echo "🔍 Pulling Job Distributor image from ECR..."
    aws ecr get-login-password --region us-west-2 --profile $AWS_PROFILE | docker login --username AWS --password-stdin $JD_ECR
    if [ $? -ne 0 ]; then
        echo -e "${CROSS_MARK} Docker login to ECR failed"
        echo -e "Please ensure you have the correct AWS credentials and try again"
        exit 1
    fi
    echo -e "${CHECK_MARK} Docker login to ECR successful"

    echo "Pulling Job Distributor image from ECR..."
    if ! docker pull 804282218731.dkr.ecr.us-west-2.amazonaws.com/job-distributor:0.12.7; then
        echo -e "${CROSS_MARK} Failed to pull Job Distributor image from ECR"
        echo -e "Please ensure you have access to the ECR repository and try again"
        exit 1
    fi
    docker tag $JD_ECR_IMAGE $JD_LOCAL_IMAGE
    echo -e "${CHECK_MARK} Job Distributor image pulled successfully"
}

# Check if Docker is installed
if command_exists docker; then
  echo -e "${CHECK_MARK} Docker is installed"
else
  echo -e "${CROSS_MARK} Docker is not installed"
  echo -e "Please install Docker and run this script again"
  exit 1
fi

if command_exists aws; then
    echo -e "${CHECK_MARK} AWS CLI is installed"
else
    echo -e "${CROSS_MARK} AWS CLI is not installed"
    echo -e "Please install AWS CLI and run this script again"
    exit 1
fi

# Check if Docker is running
if docker info > /dev/null 2>&1; then
echo -e "${CHECK_MARK} Docker is running"    
else
echo -e "${CROSS_MARK} Docker is not running"
echo -e "Please start Docker and run this script again"
exit 1
fi
echo "home $HOME"
DOCKER_SETTINGS_OK=true
echo "🔍 Checking Docker settings..."
p="$HOME/Library/Group Containers/group.com.docker/settings-store.json"
cat "$p"
x=$(realpath "$p")
echo $x
cat "$x"
if  [[ ! -f "$p" ]]; then
    p="~/Library/Group Containers/group.com.docker/settings.json"
fi
echo "  Checking for Docker settings file at $p"
if [ -f "$p" ]; then
    echo "  Found Docker settings file at $p"
    # Check if the settings are correct
    check_docker_setting "$p" "UseVirtualizationFramework" "true" || missing_docker_setting "Virtualization Framework"
    check_docker_setting "$p" "UseVirtualizationFrameworkVirtioFS" "true" || missing_docker_setting VirtioFS
    check_docker_setting "$p" "EnableDefaultDockerSocket" "true" || missing_docker_setting "enable default Docker socket"
else
    echo -e "${CROSS_MARK} Docker settings file not found at expected locations"
    echo -e "Please ensure Docker is properly configured and run this script again"
    exit 1
fi

# Check if Job Distributor image is available
echo ""


echo "🔍 Checking for Job Distributor image..."
if docker image inspect $JD_ECR_IMAGE > /dev/null 2>&1; then
  echo -e "${CHECK_MARK} Job Distributor image ($JD_IMAGE) is available"
elif docker image inspect $JD_LOCAL_IMAGE > /dev/null 2>&1; then
  echo -e "${CHECK_MARK} Job Distributor image ($JD_LOCAL_IMAGE) is available from local build"
else
  echo -e "${CROSS_MARK} Job Distributor image ($JD_IMAGE) is not available"
  echo -e "Would you like to build the Job Distributor image now? (y/n)"
  read -r build_jd
  if [[ "$build_jd" =~ ^[Yy]$ ]]; then
    echo "Building Job Distributor image..."
    
    # Check if AWS CLI is configured
    if command_exists aws; then
      pull_jd_image
    else
      build_jd_image
    fi
  else
    echo -e "${WARNING_MARK} You will need to build or pull the Job Distributor image manually before starting the CRE environment"
  fi
  if [[ "$build_jd" =~ ^[AB]$ ]]; then
    echo "Building Job Distributor image..."
    
    # Create a temporary directory for cloning
    temp_dir=$(mktemp -d)
    
    # Clone the repository
    git clone https://github.com/smartcontractkit/job-distributor "$temp_dir"
    cd "$temp_dir"
    git checkout v0.9.0
    
    # Build the Docker image
    docker build -t job-distributor:0.12.7 -f e2e/Dockerfile.e2e .
    
    # Clean up
    cd - > /dev/null
    rm -rf "$temp_dir"
    
    echo -e "${CHECK_MARK} Job Distributor image built successfully"
  else
    echo -e "${WARNING_MARK} You will need to build or pull the Job Distributor image manually before starting the CRE environment"
  fi
fi

# Check if CRE CLI is installed
echo ""
echo "🔍 Checking for CRE CLI..."

if command_exists cre_v0.2.1_darwin_arm64 || command_exists cre; then
  echo -e "${CHECK_MARK} CRE CLI is already installed"
else
  echo -e "${CROSS_MARK} CRE CLI is not installed"
  echo -e "Would you like to download and install the CRE CLI now? (y/n)"
  read -r install_cre
  
  if [[ "$install_cre" =~ ^[Yy]$ ]]; then
    # Detect architecture
    arch=$(uname -m)
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    
    case "$arch" in
      arm64)
        arch_str="arm64"
        ;;
      x86_64)
        arch_str="amd64"
        ;;
      *)
        echo -e "${CROSS_MARK} Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    
    # Download CRE CLI
    echo "Downloading CRE CLI v0.2.1 for ${os}_${arch_str}..."
    
    if command_exists gh; then
      gh release download v0.2.1 --repo smartcontractkit/dev-platform --pattern "*${os}_${arch_str}.tar.gz"
    else
      echo -e "${CROSS_MARK} GitHub CLI is not installed"
      echo -e "Please install GitHub CLI or download CRE CLI manually from https://github.com/smartcontractkit/dev-platform/releases/tag/v0.2.1"
      exit 1
    fi
    
    # Extract the archive
    echo "Extracting CRE CLI..."
    tar -xf cre_v0.2.1_${os}_${arch_str}.tar.gz
    rm cre_v0.2.1_${os}_${arch_str}.tar.gz
    
    # Remove quarantine attribute on macOS
    if [[ "$os" == "darwin" ]]; then
      xattr -d com.apple.quarantine cre_v0.2.1_${os}_${arch_str} 2>/dev/null || true
    fi
    
    # Make executable
    chmod +x cre_v0.2.1_${os}_${arch_str}
    
    # Create a symlink in current directory
    ln -sf cre_v0.2.1_${os}_${arch_str} cre
    
    echo -e "${CHECK_MARK} CRE CLI installed to $(pwd)/cre"
    echo -e "${WARNING_MARK} Add this directory to your PATH or move the CRE binary to a directory in your PATH"
    echo -e "   You can run: export PATH=\"$(pwd):\$PATH\""
  else
    echo -e "${WARNING_MARK} You will need to install CRE CLI manually"
  fi
fi

# Check for capability binaries
echo ""
echo "🔍 Checking for capability binaries..."
echo -e "${WARNING_MARK} You may need capability binaries depending on your use case"
echo -e "   Some capabilities like cron, log-event-trigger, or read-contract might not be embedded in your Chainlink image"
echo -e "   If needed, download from https://github.com/smartcontractkit/capabilities/releases/tag/v1.0.2-alpha"
echo -e "   Or use: gh release download v1.0.2-alpha --repo smartcontractkit/capabilities --pattern 'amd64_cron'"

# Print summary
echo ""
echo "✅ Setup Summary:"
if [[ "$DOCKER_SETTINGS_OK" == "true" ]]; then
  echo -e "   ${CHECK_MARK} Docker settings are correctly configured"
else
  echo -e "   ${CROSS_MARK} Some Docker settings need adjustment (see above)"
fi

if docker image inspect job-distributor:0.9.0 > /dev/null 2>&1; then
  echo -e "   ${CHECK_MARK} Job Distributor image is available"
else
  echo -e "   ${CROSS_MARK} Job Distributor image is not available"
fi

if command_exists cre_v0.2.1_darwin_arm64 || command_exists cre; then
  echo -e "   ${CHECK_MARK} CRE CLI is installed"
else
  echo -e "   ${CROSS_MARK} CRE CLI is not installed"
fi

echo ""
echo "🚀 Next Steps:"
echo "1. Navigate to the CRE environment directory: cd core/scripts/cre/environment"
echo "2. Start the environment: go run main.go env start"
echo "   Optional: Add --with-example to start with an example workflow"
echo "   Optional: Add --with-plugins-docker-image to use a pre-built image with capabilities"
echo ""
echo "For more information, see the documentation in core/scripts/cre/environment/docs.md"