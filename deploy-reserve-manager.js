#!/usr/bin/env node

const { execSync, spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const { ethers } = require('ethers');
const solc = require('solc');

let gethProcess;
let shutdownInProgress = false;

async function main() {
    console.log('Starting deployment script...');
    
    const port = 8545;
    const rpcUrl = `http://127.0.0.1:${port}`;
    const dataDir = path.join(__dirname, 'geth-dev-data');
    
    try {
        // Clean up any existing data directory
        if (fs.existsSync(dataDir)) {
            execSync(`rm -rf ${dataDir}`);
        }
        
        // Start geth in dev mode
        console.log('Starting geth node in dev mode...');
        gethProcess = spawn('geth', [
            '--dev',
            '--datadir', dataDir,
            '--http',
            '--http.addr', 'localhost',
            '--http.port', port.toString(),
            '--http.api', 'eth,net,web3,personal,admin,debug,miner',
            '--http.corsdomain', '*',
            '--http.vhosts', '*',
            '--allow-insecure-unlock',
            '--nodiscover',
            '--maxpeers', '0',
            '--mine',
            '--verbosity', '3'
        ]);
        
        // Capture geth output for debugging
        let httpStarted = false;
        
        gethProcess.stdout.on('data', (data) => {
            const output = data.toString();
            console.log('Geth stdout:', output.trim());
        });
        
        gethProcess.stderr.on('data', (data) => {
            const output = data.toString();
            if (output.includes('HTTP endpoint opened') || output.includes('HTTP server started')) {
                httpStarted = true;
                console.log('✓ HTTP RPC endpoint is ready!');
            }
            // Log all output for debugging
            console.log('Geth:', output.trim());
        });
        
        gethProcess.on('error', (err) => {
            console.error('Failed to start geth:', err);
            process.exit(1);
        });
        
        // Wait for HTTP endpoint to be ready
        console.log('Waiting for HTTP endpoint to start...');
        let waited = 0;
        while (!httpStarted && waited < 30000) {
            await new Promise(resolve => setTimeout(resolve, 500));
            waited += 500;
        }
        
        if (!httpStarted) {
            throw new Error('HTTP endpoint did not start within 30 seconds');
        }
        
        // Additional wait to ensure RPC is fully ready
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Test the connection
        console.log('Testing RPC connection...');
        // Use 127.0.0.1 explicitly to avoid IPv6 issues
        const provider = new ethers.JsonRpcProvider('http://127.0.0.1:' + port);
        
        // Get network info
        const network = await provider.getNetwork();
        console.log('Connected to network:', network.chainId.toString());
        
        // Get accounts
        const accounts = await provider.send('eth_accounts', []);
        console.log('Available accounts:', accounts);
        
        if (accounts.length === 0) {
            // Create a new account if none exist
            console.log('Creating new account...');
            const newAccount = await provider.send('personal_newAccount', ['']);
            accounts.push(newAccount);
        }
        
        const devAccount = accounts[0];
        console.log('Using account:', devAccount);
        
        // Generate a new wallet to avoid nonce conflicts
        const wallet = ethers.Wallet.createRandom(provider);
        const devPrivateKey = wallet.privateKey;
        
        try {
            // Fund the new wallet from coinbase
            console.log('Funding new wallet from coinbase...');
            console.log('New wallet address:', wallet.address);
            
            await provider.send('eth_sendTransaction', [{
                from: devAccount,
                to: wallet.address,
                value: '0x' + ethers.parseEther('10').toString(16)
            }]);
            
            // Wait for the transaction
            await new Promise(resolve => setTimeout(resolve, 2000));
            const balance = await provider.getBalance(wallet.address);
            
            console.log('Wallet balance:', ethers.formatEther(balance), 'ETH');
            
            // Compile contracts
            console.log('\nCompiling contracts...');
            const contractPath = path.join(__dirname, 'core/services/workflows/cmd/cre/examples/v2/por/bindings/solc');
            const sources = {
                'IReserveManager.sol': {
                    content: fs.readFileSync(path.join(contractPath, 'IReserveManager.sol'), 'utf8')
                },
                'ReserveManager.sol': {
                    content: fs.readFileSync(path.join(contractPath, 'ReserveManager.sol'), 'utf8')
                },
                'SimpleERC20.sol': {
                    content: fs.readFileSync(path.join(contractPath, 'SimpleERC20.sol'), 'utf8')
                }
            };
            
            const input = {
                language: 'Solidity',
                sources,
                settings: {
                    optimizer: {
                        enabled: false
                    },
                    evmVersion: 'london',
                    outputSelection: {
                        '*': {
                            '*': ['*']
                        }
                    }
                }
            };
            
            const output = JSON.parse(solc.compile(JSON.stringify(input)));
            
            if (output.errors && output.errors.some(e => e.severity === 'error')) {
                console.error('Compilation errors:', output.errors);
                throw new Error('Contract compilation failed');
            }
            
            console.log('✓ Contracts compiled successfully');
            
            // Deploy ReserveManager contract
            console.log('\nDeploying ReserveManager contract...');
            const reserveContract = output.contracts['ReserveManager.sol']['ReserveManager'];
            const reserveAbi = reserveContract.abi;
            const reserveBytecode = '0x' + reserveContract.evm.bytecode.object;
            
            const reserveFactory = new ethers.ContractFactory(reserveAbi, reserveBytecode, wallet);
            const reserveManager = await reserveFactory.deploy();
            await reserveManager.waitForDeployment();
            
            const reserveManagerAddress = await reserveManager.getAddress();
            console.log('✓ ReserveManager deployed at:', reserveManagerAddress);
            
            // Deploy SimpleERC20 contract
            console.log('\nDeploying SimpleERC20 token...');
            const tokenContract = output.contracts['SimpleERC20.sol']['SimpleERC20'];
            const tokenAbi = tokenContract.abi;
            const tokenBytecode = '0x' + tokenContract.evm.bytecode.object;
            
            const tokenFactory = new ethers.ContractFactory(tokenAbi, tokenBytecode, wallet);
            const token = await tokenFactory.deploy("Test Token", "TEST", 1000000); // 1M tokens
            await token.waitForDeployment();
            
            const tokenAddress = await token.getAddress();
            console.log('✓ SimpleERC20 deployed at:', tokenAddress);
            
            console.log('\n=== Deployment Successful ===');
            console.log('ReserveManager Address:', reserveManagerAddress);
            console.log('ERC20 Token Address:', tokenAddress);
            console.log('Private Key:', devPrivateKey);
            console.log('RPC URL:', rpcUrl);
            console.log('Wallet Address:', wallet.address);
            console.log('============================\n');
            
            // Keep the process running
            console.log('Geth node is running. Press Ctrl+C to stop...');
            
            // Prevent the script from exiting
            await new Promise(() => {});
            
        } catch (walletError) {
            // If the known private key doesn't work, generate a new wallet
            console.log('Generating new wallet...');
            const newWallet = ethers.Wallet.createRandom(provider);
            
            // Fund the new wallet from coinbase
            console.log('Funding new wallet from coinbase...');
            await provider.send('eth_sendTransaction', [{
                from: devAccount,
                to: newWallet.address,
                value: '0x' + ethers.parseEther('10').toString(16)
            }]);
            
            // Wait for the transaction
            await new Promise(resolve => setTimeout(resolve, 2000));
            
            const walletBalance = await provider.getBalance(newWallet.address);
            console.log('New wallet balance:', ethers.formatEther(walletBalance), 'ETH');
            
            // Compile and deploy using coinbase
            console.log('\nCompiling contracts...');
            const contractPath = path.join(__dirname, 'core/services/workflows/cmd/cre/examples/v2/por/bindings/solc');
            const sources = {
                'IReserveManager.sol': {
                    content: fs.readFileSync(path.join(contractPath, 'IReserveManager.sol'), 'utf8')
                },
                'ReserveManager.sol': {
                    content: fs.readFileSync(path.join(contractPath, 'ReserveManager.sol'), 'utf8')
                },
                'SimpleERC20.sol': {
                    content: fs.readFileSync(path.join(contractPath, 'SimpleERC20.sol'), 'utf8')
                }
            };
            
            const input = {
                language: 'Solidity',
                sources,
                settings: {
                    optimizer: {
                        enabled: false
                    },
                    evmVersion: 'london',
                    outputSelection: {
                        '*': {
                            '*': ['*']
                        }
                    }
                }
            };
            
            const output = JSON.parse(solc.compile(JSON.stringify(input)));
            
            if (output.errors && output.errors.some(e => e.severity === 'error')) {
                console.error('Compilation errors:', output.errors);
                throw new Error('Contract compilation failed');
            }
            
            console.log('✓ Contracts compiled successfully');
            
            // Deploy ReserveManager contract
            console.log('\nDeploying ReserveManager contract...');
            const reserveContract = output.contracts['ReserveManager.sol']['ReserveManager'];
            const reserveAbi = reserveContract.abi;
            const reserveBytecode = '0x' + reserveContract.evm.bytecode.object;
            
            const reserveFactory = new ethers.ContractFactory(reserveAbi, reserveBytecode, newWallet);
            const reserveManager = await reserveFactory.deploy();
            await reserveManager.waitForDeployment();
            
            const reserveManagerAddress = await reserveManager.getAddress();
            console.log('✓ ReserveManager deployed at:', reserveManagerAddress);
            
            // Deploy SimpleERC20 contract
            console.log('\nDeploying SimpleERC20 token...');
            const tokenContract = output.contracts['SimpleERC20.sol']['SimpleERC20'];
            const tokenAbi = tokenContract.abi;
            const tokenBytecode = '0x' + tokenContract.evm.bytecode.object;
            
            const tokenFactory = new ethers.ContractFactory(tokenAbi, tokenBytecode, newWallet);
            const token = await tokenFactory.deploy("Test Token", "TEST", 1000000); // 1M tokens
            await token.waitForDeployment();
            
            const tokenAddress = await token.getAddress();
            console.log('✓ SimpleERC20 deployed at:', tokenAddress);
            
            console.log('\n=== Deployment Successful ===');
            console.log('ReserveManager Address:', reserveManagerAddress);
            console.log('ERC20 Token Address:', tokenAddress);
            console.log('Private Key:', newWallet.privateKey);
            console.log('RPC URL:', rpcUrl);
            console.log('Wallet Address:', newWallet.address);
            console.log('============================\n');
            
            // Keep the process running
            console.log('Geth node is running. Press Ctrl+C to stop...');
            
            // Prevent the script from exiting
            await new Promise(() => {});
        }
        
    } catch (error) {
        console.error('Deployment failed:', error);
        cleanup();
        process.exit(1);
    }
}

function cleanup() {
    if (shutdownInProgress) return;
    shutdownInProgress = true;
    
    console.log('\nShutting down...');
    
    if (gethProcess) {
        gethProcess.kill('SIGTERM');
        // Give it time to shut down gracefully
        setTimeout(() => {
            if (gethProcess && !gethProcess.killed) {
                gethProcess.kill('SIGKILL');
            }
        }, 5000);
    }
    
    // Clean up data directory
    const dataDir = path.join(__dirname, 'geth-dev-data');
    if (fs.existsSync(dataDir)) {
        try {
            execSync(`rm -rf ${dataDir}`);
        } catch (e) {
            console.error('Failed to clean up data directory:', e.message);
        }
    }
}

// Handle signals
process.on('SIGINT', () => {
    cleanup();
    process.exit(0);
});

process.on('SIGTERM', () => {
    cleanup();
    process.exit(0);
});

process.on('exit', () => {
    cleanup();
});

// Check if required dependencies are installed
try {
    require.resolve('ethers');
    require.resolve('solc');
} catch (e) {
    console.error('Missing dependencies. Please run:');
    console.error('npm install ethers solc');
    process.exit(1);
}

// Run the script
main().catch(error => {
    console.error('Fatal error:', error);
    cleanup();
    process.exit(1);
});