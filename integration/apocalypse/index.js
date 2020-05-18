const fs = require('fs')
const ethers = require('ethers')
const execa = require('execa')
const _ = require('lodash')
const { contract } = require('@chainlink/test-helpers')

// Contracts
const oracleJson = require('../../evm-contracts/abi/v0.6/Oracle.json')
const fluxAggregatorJson = require('../../evm-contracts/abi/v0.6/FluxAggregator.json')

// RPC
const RPC_GETH1  = 'http://localhost:8545'
const RPC_GETH2  = 'http://localhost:18545'
const RPC_PARITY = 'http://localhost:28545'
const RPC_PROVIDERS = [RPC_GETH1, RPC_GETH2]

// Personas
const accounts = require('./config/accounts.json')
const carol = new ethers.Wallet(accounts.carol.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const loki = new ethers.Wallet(accounts.loki.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const lokiJr = new ethers.Wallet(accounts['loki-jr'].privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const geth1 = new ethers.Wallet(accounts.geth1.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const geth2 = new ethers.Wallet(accounts.geth2.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const parity = new ethers.Wallet(accounts.parity.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const neil = new ethers.Wallet(accounts.neil.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const nelly = new ethers.Wallet(accounts.nelly.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))

const EXTERNAL_ADAPTER_URL = 'http://external_adapter:6644'
const CHAINLINK_URL_NEIL  = 'http://localhost:6688'
const CHAINLINK_URL_NELLY = 'http://localhost:6689'

async function main() {
    console.log('Awaiting Geth DAG generation...')
    await gethDAGGenerationFinished()

    console.log('Deploying Flux Monitor contracts...')
    let { linkToken, fluxAggregator } = await deployFluxMonitorContracts(carol, [ neil.address, nelly.address ])

    console.log('Initializing price feed...')
    await setPriceFeedValue(EXTERNAL_ADAPTER_URL, 100)

    console.log('Initiating transaction tornado...')
    await mimicRegularTraffic(200, RPC_PROVIDERS)

    console.log('Waiting for things to get really bad...')
    await sleep(20000)

    console.log('Adding job specs to Chainlink nodes...')
    let envNeil = { CLIENT_NODE_URL: CHAINLINK_URL_NEIL, ROOT: '/clroot' }
    await run(`chainlink -j admin login -f /apicredentials`, envNeil)
    await run(`chainlink -j jobs create '${JSON.stringify(makeJobSpecFluxMonitor(fluxAggregator.address, EXTERNAL_ADAPTER_URL))}'`, envNeil)
    let envNelly = { CLIENT_NODE_URL: CHAINLINK_URL_NELLY, ROOT: '/clroot' }
    await run(`chainlink -j admin login -f /apicredentials`, envNelly)
    let jobSpec = JSON.parse(await run(`chainlink -j jobs create '${JSON.stringify(makeJobSpecFluxMonitor(fluxAggregator.address, EXTERNAL_ADAPTER_URL))}'`, envNelly))

    await successfulJobRuns(jobSpec.id, 1, envNeil)
    await successfulJobRuns(jobSpec.id, 1, envNelly)
    fluxAggregator.
}

async function successfulJobRuns(jobSpecID, num, env) {
    while (true) {
        let runs = JSON.parse(await run(`chainlink -j runs list`, env))
        runs = runs.filter(run => run.jobId === jobSpecID)

        for (let run of runs) {
            if (run.status === 'errored') {
                throw new Error('job run errored')
            }
        }

        await sleep(3000)
    }
}

async function setPriceFeedValue(feedURL, value) {
    const url = new URL('result', feedURL).href
    const response = await fetch(url, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ result: value }),
    })
    if (!response.ok) {
        throw new Error('failed to set external adapter price feed')
    }
}

async function deployDirectRequestContracts() {
    let linkToken = await deployLINK(carol, undefined)
    console.log('  - LINK token:', linkToken.address)

    let oracle = await deployOracle(carol, undefined, {
        linkTokenAddress: linkToken.address,
    })
    console.log('  - Oracle:', oracle.address)

    return {
        linkToken,
        oracle,
    }
}

async function deployFluxMonitorContracts(carol, oracles) {
    let linkToken = await deployLINK(carol, undefined)
    console.log('  - LINK token:', linkToken.address)

    let fluxAggregator = await deployFluxAggregator(carol, undefined, {
        linkTokenAddress: linkToken.address,
        paymentAmount: '100',    // LINK-sats
        timeout: 300,            // seconds
    })
    console.log('  - Flux Aggregator:', fluxAggregator.address)

    await (await fluxAggregator.addOracles(oracles, oracles, 1, 2, 1)).wait()

    return {
        linkToken,
        fluxAggregator,
    }
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


async function mimicRegularTraffic(totalTxsPerSecond, rpcProviders) {
    const numAccounts = 200

    let accounts = await makeRandomAccounts(numAccounts, rpcProviders)
    for (let account of accounts) {
        sendRandomTransactions(account, accounts)
    }
}

async function makeRandomAccounts(num, rpcProviders) {
    let senders = []
    for (let providerURL of rpcProviders) {
        let wallet = new ethers.Wallet(require('./config/accounts.json').carol.privkey.substring(2), new ethers.providers.JsonRpcProvider(providerURL))
        senders.push({
            providerURL: providerURL,
            nonce: await wallet.provider.getTransactionCount(wallet.address, 'pending'),
            wallet: wallet,
        })
    }
    let jobs = Array(num).fill(null).map((_, i) => {
        let sender = senders[i % senders.length]
        return {
            providerURL: sender.providerURL,
            wallet: ethers.Wallet.createRandom().connect(new ethers.providers.JsonRpcProvider(sender.providerURL)),
            sender: sender,
        }
    })
    // Fund the accounts
    await Promise.all(
        jobs.map(job => {
            let nonce = job.sender.nonce
            job.sender.nonce++
            return job.sender.wallet.sendTransaction({
                to: job.wallet.address,
                value: ethers.utils.parseUnits('5', 'ether'),
                gasPrice: ethers.utils.parseUnits('20', 'gwei'),
                nonce: nonce,
            }).catch(err => {
                console.log(err, 'nonce =', nonce, job.sender.wallet.address, job.sender.nonce, job.providerURL)
            })
        })
    )
    return jobs.map(job => job.wallet)
}

function sendRandomTransactions(fromAccount, toAccounts) {
    function randomAccount() {
        let i = Math.floor(Math.random() * Math.floor(toAccounts.length - 1))
        return toAccounts[i]
    }

    async function send() {
        let msBetweenTxs = 500
        try {
            // Re-read the config each time so that we can control the congestion dynamically
            msBetweenTxs = getConfig().randomTraffic.msBetweenTxs

            await fromAccount.sendTransaction({
                to: randomAccount().address,
                value: ethers.utils.parseUnits('1', 'wei'),
                gasPrice: ethers.utils.parseUnits('20', 'gwei'),
            })
        } catch (err) {}

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
        await sleep(5000)
    }
    block = null
    while (!block) {
        block = await geth2.provider.getBlock(2)
        await sleep(5000)
    }
}

async function sleep(ms) {
    return new Promise((resolve) => {
        setTimeout(resolve, ms)
    })
}

function getConfig() {
    return JSON.parse(fs.readFileSync('./config/config.json').toString())
}

async function run(cmd, env) {
    return (await execa.command(cmd, {all: true, env: env})).all
}


function makeJobSpecEthlog(oracleContractAddress) {
    return {
        initiators: [{ type: "ethlog", params: { address: oracleContractAddress } }],
        tasks: [
            { type: "httpGet", params: { url: "http://localhost:8000/data" } },
            { type: "jsonParse", params: { path: ["bryn", "age"] } },
            { type: "multiply", params: { times: 10000 } },
            { type: "ethtx" },
        ],
    }
}

function makeJobSpecFluxMonitor(aggregatorContractAddress, feedAddr) {
    return {
        initiators: [
            {
                type: "fluxmonitor",
                params: {
                    address: aggregatorContractAddress,
                    requestData: {
                        data: {
                            coin: 'ETH',
                            market: 'USD',
                        },
                    },
                    feeds: [ feedAddr ],
                    precision: 2,
                    threshold: 5,
                    idleTimer: {
                      disabled: true,
                    },
                    pollTimer: {
                      period: '5s',
                    },
                }
            }
        ],
        tasks: [
            {
                type: "multiply",
                confirmations: null,
                params: {
                    times: 100,
                }
            },
            {
                type: "ethuint256",
                confirmations: null,
                params: {}
            },
            {
                type: "ethtx",
                confirmations: null,
                params: {}
            },
        ],
    }
}


main()

