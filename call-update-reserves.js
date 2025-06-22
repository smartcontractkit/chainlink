#!/usr/bin/env node

const { ethers } = require('ethers');
const fs = require('fs');
const path = require('path');

async function main() {
    // Parse command line arguments
    const args = process.argv.slice(2);
    
    if (args.length < 5) {
        console.error('Usage: node call-update-reserves.js <rpc-url> <private-key> <contract-address> <total-minted> <total-reserve>');
        console.error('Example: node call-update-reserves.js http://localhost:8545 0xac0974... 0x5FbDB... 1000000 2000000');
        process.exit(1);
    }
    
    const [rpcUrl, privateKey, contractAddress, totalMinted, totalReserve] = args;
    
    try {
        // Connect to the network
        console.log('Connecting to RPC:', rpcUrl);
        const provider = new ethers.JsonRpcProvider(rpcUrl);
        
        // Create wallet
        const wallet = new ethers.Wallet(privateKey, provider);
        console.log('Using wallet address:', wallet.address);
        
        // Check balance
        const balance = await provider.getBalance(wallet.address);
        console.log('Wallet balance:', ethers.formatEther(balance), 'ETH');
        
        // Load the contract ABI
        const contractPath = path.join(__dirname, 'core/services/workflows/cmd/cre/examples/v2/por/bindings/solc');
        const sources = {
            'IReserveManager.sol': {
                content: fs.readFileSync(path.join(contractPath, 'IReserveManager.sol'), 'utf8')
            },
            'ReserveManager.sol': {
                content: fs.readFileSync(path.join(contractPath, 'ReserveManager.sol'), 'utf8')
            }
        };
        
        // Compile to get ABI
        const solc = require('solc');
        const input = {
            language: 'Solidity',
            sources,
            settings: {
                outputSelection: {
                    '*': {
                        '*': ['abi']
                    }
                }
            }
        };
        
        const output = JSON.parse(solc.compile(JSON.stringify(input)));
        
        if (output.errors && output.errors.some(e => e.severity === 'error')) {
            console.error('Compilation errors:', output.errors);
            throw new Error('Failed to compile contracts for ABI');
        }
        
        const abi = output.contracts['ReserveManager.sol']['ReserveManager'].abi;
        
        // Create contract instance
        const contract = new ethers.Contract(contractAddress, abi, wallet);
        console.log('\nContract address:', contractAddress);
        
        // Read current values
        console.log('\n--- Current Contract State ---');
        const currentMinted = await contract.lastTotalMinted();
        const currentReserve = await contract.lastTotalReserve();
        console.log('Last Total Minted:', currentMinted.toString());
        console.log('Last Total Reserve:', currentReserve.toString());
        
        // Prepare the update
        const updateData = {
            totalMinted: ethers.parseUnits(totalMinted, 0), // assuming wei units
            totalReserve: ethers.parseUnits(totalReserve, 0)
        };
        
        console.log('\n--- Calling updateReserves ---');
        console.log('New Total Minted:', updateData.totalMinted.toString());
        console.log('New Total Reserve:', updateData.totalReserve.toString());
        
        // Estimate gas
        const estimatedGas = await contract.updateReserves.estimateGas(updateData);
        console.log('Estimated gas:', estimatedGas.toString());
        
        // Call updateReserves
        const tx = await contract.updateReserves(updateData);
        console.log('\nTransaction sent:', tx.hash);
        
        // Wait for confirmation
        console.log('Waiting for confirmation...');
        const receipt = await tx.wait();
        console.log('Transaction confirmed in block:', receipt.blockNumber);
        console.log('Gas used:', receipt.gasUsed.toString());
        
        // Check for events
        if (receipt.logs.length > 0) {
            console.log('\n--- Events Emitted ---');
            for (const log of receipt.logs) {
                try {
                    const parsedLog = contract.interface.parseLog(log);
                    if (parsedLog.name === 'RequestReserveUpdate') {
                        console.log('RequestReserveUpdate event:');
                        console.log('  Request ID:', parsedLog.args[0].toString());
                    }
                } catch (e) {
                    // Skip if not our contract's event
                }
            }
        }
        
        // Read updated values
        console.log('\n--- Updated Contract State ---');
        const newMinted = await contract.lastTotalMinted();
        const newReserve = await contract.lastTotalReserve();
        console.log('Last Total Minted:', newMinted.toString());
        console.log('Last Total Reserve:', newReserve.toString());
        
        console.log('\n✅ Update completed successfully!');
        
    } catch (error) {
        console.error('\n❌ Error:', error.message);
        if (error.reason) {
            console.error('Reason:', error.reason);
        }
        if (error.transaction) {
            console.error('Failed transaction:', error.transaction);
        }
        process.exit(1);
    }
}

// Check if ethers and solc are installed
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
    process.exit(1);
});