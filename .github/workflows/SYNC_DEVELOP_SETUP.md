# Sync Develop Workflow Setup

This document provides instructions for setting up the "Sync develop from smartcontractkit/chainlink" workflow.

## Overview

The `sync-develop-from-smartcontractkit-chainlink.yml` workflow automatically syncs the `develop` branch from the upstream repository (`smartcontractkit/chainlink`) to your fork every 30 minutes.

## Required Setup

### Prerequisites: Create the develop branch

Before the sync workflow can run, the `develop` branch must exist in your fork. Follow these steps to create it:

1. Go to the Actions tab in your forked repository
2. Select "Create develop branch (one-time setup)" workflow from the left sidebar
3. Click "Run workflow" → "Run workflow" button
4. Wait for the workflow to complete (it should take less than a minute)
5. Verify the develop branch was created by checking the branches in your repository

Once the develop branch exists, the sync workflow will be able to run automatically.

### Setup Automatic Syncing

To enable automatic syncing, you need to create a Personal Access Token (PAT) and add it as a repository secret.

### Step 1: Create a Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Direct link: https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Configure the token:
   - **Note**: `Chainlink Fork Sync Token`
   - **Expiration**: Choose an appropriate expiration (recommended: 90 days or No expiration for continuous sync)
   - **Scopes**: Select the following permissions:
     - ✅ `repo` (Full control of private repositories) - Required to push to your repository
4. Click "Generate token"
5. **Important**: Copy the token immediately - you won't be able to see it again!

### Step 2: Add the Token as a Repository Secret

1. Go to your forked repository on GitHub
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Configure the secret:
   - **Name**: `PAT_TOKEN` (must be exactly this name)
   - **Value**: Paste the Personal Access Token you created in Step 1
5. Click "Add secret"

### Step 3: Verify the Setup

After adding the secret, the workflow will automatically use it on the next scheduled run (every 30 minutes).

To manually trigger a test:
1. Go to Actions tab in your repository
2. Select "Sync develop from smartcontractkit/chainlink" workflow
3. If the workflow file includes a `workflow_dispatch` trigger, you can click "Run workflow" to trigger it manually. (By default, this workflow only runs on a schedule.)

Alternatively, wait for the next scheduled run and check the workflow logs to ensure it completes successfully.

## Troubleshooting

### Authentication Failed Error
- Verify the `PAT_TOKEN` secret exists and is spelled correctly
- Ensure the token has the `repo` scope enabled
- Check if the token has expired and create a new one if needed

### Push Permission Denied
- The PAT must have write access to your fork
- Verify you're using a token associated with an account that has push permissions to the repository

### Workflow Not Running
- This workflow only runs on forks (not on `smartcontractkit/chainlink`)
- Check the Actions tab to see if the workflow is enabled
- Verify the workflow file is present in the `.github/workflows` directory

## Security Note

Never commit your Personal Access Token directly in code or configuration files. Always use GitHub Secrets to store sensitive credentials.
