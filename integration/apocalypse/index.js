const fs = require('fs')
const ethers = require('ethers')
const _ = require('lodash')
const { contract } = require('@chainlink/test-helpers')

// Contracts
const oracleJson = require('../../evm-contracts/abi/v0.6/Oracle.json')
const fluxAggregatorJson = require('../../evm-contracts/abi/v0.6/FluxAggregator.json')

// RPC
const RPC_GETH1  = 'http://localhost:8545'
const RPC_GETH2  = 'http://localhost:18545'
const RPC_PARITY = 'http://localhost:28545'
const rpcProviders = [RPC_GETH1, RPC_GETH2]

// Personas
const accounts = require('./config/accounts.json')
const carol = new ethers.Wallet(accounts.carol.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const loki = new ethers.Wallet(accounts.loki.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const lokiJr = new ethers.Wallet(accounts['loki-jr'].privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const geth1 = new ethers.Wallet(accounts.geth1.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const geth2 = new ethers.Wallet(accounts.geth2.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const parity = new ethers.Wallet(accounts.parity.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))

async function main() {
    let linkToken = await deployLINK(carol, undefined)
    console.log('LINK token:', linkToken.address)

    let oracle = await deployOracle(carol, undefined, {
        linkTokenAddress: linkToken.address,
    })
    console.log('Oracle:', oracle.address)

    let fluxAggregator = await deployFluxAggregator(carol, undefined, {
        linkTokenAddress: linkToken.address,
        paymentAmount: '100',    // LINK-sats
        timeout: 300,            // seconds
    })
    console.log('Flux Aggregator:', fluxAggregator.address)

    mimicRegularTraffic(200)
}

async function deployContract({ Factory, name, signer }, ...deployArgs) {
    const contractFactory = new Factory(signer)
    const contract = await contractFactory.deploy(...deployArgs, { gasPrice: ethers.utils.parseUnits('50', 'gwei') })
    await contract.deployed()
    return contract
}

async function deployLINK(wallet, nonce) {
    const linkToken = await deployContract({
        Factory: contract.LinkTokenFactory,
        name: 'LinkToken',
        signer: wallet,
        nonce: nonce,
    })
    return linkToken
}

async function deployOracle(wallet, nonce, { linkTokenAddress }) {
    let oracleFactory = new ethers.ContractFactory(oracleJson.compilerOutput.abi, oracleJson.compilerOutput.evm.bytecode, wallet)
    let oracle = await oracleFactory.deploy(linkTokenAddress, { gasPrice: ethers.utils.parseUnits('50', 'gwei') })
    await oracle.deployed()
    return oracle
}

async function deployFluxAggregator(wallet, nonce, { linkTokenAddress, paymentAmount, timeout }) {
    let description = ethers.utils.formatBytes32String('xyzzy')
    let fluxAggregatorFactory = new ethers.ContractFactory(fluxAggregatorJson.compilerOutput.abi, fluxAggregatorJson.compilerOutput.evm.bytecode.object, wallet)
    let fluxAggregator = await fluxAggregatorFactory.deploy(linkTokenAddress, paymentAmount, timeout, 2, description, {
        gasPrice: ethers.utils.parseUnits('50', 'gwei'),
    })
    await fluxAggregator.deployed()
    return fluxAggregator
}


async function mimicRegularTraffic(totalTxsPerSecond) {
    const numAccounts = 200

    let accounts = await makeRandomAccounts(numAccounts, rpcProviders)
    for (let account of accounts) {
        sendRandomTransactions(account, accounts)
    }
}

async function makeRandomAccounts(num, rpcProviders) {
    let accounts = Array(num).fill(null).map((_, i) => ethers.Wallet.createRandom().connect(new ethers.providers.JsonRpcProvider(rpcProviders[i % rpcProviders.length])))
    // Fund the accounts
    let nonce = await carol.provider.getTransactionCount(carol.address, 'pending')
    await Promise.all(
        accounts.map(account => carol.sendTransaction({
            to: account.address,
            value: ethers.utils.parseUnits('5', 'ether'),
            nonce: nonce++,
        }))
    )
    return accounts
}

function sendRandomTransactions(account, accounts) {
    function randomAccount() {
        let i = Math.floor(Math.random() * Math.floor(accounts.length - 1))
        return accounts[i]
    }

    async function send() {
        let msBetweenTxs = 500
        try {
            // Re-read the config each time so that we can control the congestion dynamically
            msBetweenTxs = getConfig().randomTraffic.msBetweenTxs

            let tx = await account.sendTransaction({
                to: randomAccount().address,
                value: ethers.utils.parseUnits('1', 'wei'),
                gasPrice: ethers.utils.parseUnits('20', 'gwei'),
            })
            // console.log(tx)
        } catch (err) {
            // console.log(err)
        }
        setTimeout(send, msBetweenTxs)
    }
    send()
}

async function connectPeers(rpcProvider) {
    let resp = await (await fetch('http://localhost:28545', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            method: 'parity_addReservedPeer',
            params: ['enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@172.17.0.4:30303'],
            id:1,
            jsonrpc:"2.0"
        }),
    })).text()
    console.log(resp)
    resp = await (await fetch('http://localhost:28545', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            method: 'parity_addReservedPeer',
            params: ['enode://c1cad3139b0ab583de214e3d64f7fb7793995023559f7fa1e6b01e87603145ca8e60d5d9f8e23d08df3d1c0c82294bd9515b729efec210f060b2fe3a193f9ae0@172.17.0.6:30303'],
            id:1,
            jsonrpc:"2.0"
        }),
    })).text()
    console.log(resp)
}

async function gethDAGGenerationFinished() {
    let block
    while (!block) {
        block = await geth1.provider.getBlock(2)
    }
    block = null
    while (!block) {
        block = await geth2.provider.getBlock(2)
    }
}

function getConfig() {
    return JSON.parse(fs.readFileSync('./config/config.json').toString())
}


main()

