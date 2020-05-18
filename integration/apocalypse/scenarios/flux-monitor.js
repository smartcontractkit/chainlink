const fs = require('fs')
const ethers = require('ethers')
const execa = require('execa')
const _ = require('lodash')
const assert = require('assert')

const { gethDAGGenerationFinished } = require('./lib/ethnodes')
const { deployFluxMonitorContracts } = require('./lib/contracts')
const { chainlinkLogin, addJobSpec, setPriceFeedValue, successfulJobRuns, makeJobSpecFluxMonitor } = require('./lib/chainlink')
const { mimicRegularTraffic } = require('./lib/apocalypse')
const { run, sleep, getConfig } = require('./lib/utils')

// RPC
const RPC_GETH1  = 'http://localhost:8545'
const RPC_GETH2  = 'http://localhost:18545'
const RPC_PARITY = 'http://localhost:28545'
const RPC_ETHEREUM_PROVIDERS = [RPC_GETH1, RPC_GETH2]
const EXTERNAL_ADAPTER_INTERNAL_URL = 'http://external_adapter:6644'
const EXTERNAL_ADAPTER_EXTERNAL_URL = 'http://localhost:6644'
const CHAINLINK_URL_NEIL  = 'http://localhost:6688'
const CHAINLINK_URL_NELLY = 'http://localhost:6689'

// Personas
const accounts = require('../config/accounts.json')
const carol = new ethers.Wallet(accounts.carol.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const loki = new ethers.Wallet(accounts.loki.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const lokiJr = new ethers.Wallet(accounts['loki-jr'].privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const geth1 = new ethers.Wallet(accounts.geth1.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH1))
const geth2 = new ethers.Wallet(accounts.geth2.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const parity = new ethers.Wallet(accounts.parity.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const neil = new ethers.Wallet(accounts.neil.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))
const nelly = new ethers.Wallet(accounts.nelly.privkey.substring(2), new ethers.providers.JsonRpcProvider(RPC_GETH2))


async function main() {
    console.log('Awaiting Geth DAG generation...')
    await gethDAGGenerationFinished([ geth1.provider, geth2.provider ])

    console.log('Deploying Flux Monitor contracts...')
    let { linkToken, fluxAggregator } = await deployFluxMonitorContracts(carol, [ neil.address, nelly.address ])

    console.log('Initializing price feed...')
    await setPriceFeedValue(EXTERNAL_ADAPTER_EXTERNAL_URL, 100)

    // console.log('Initiating transaction tornado...')
    // await mimicRegularTraffic({
    //     funderPrivkey: accounts.carol.privkey.substring(2),
    //     numAccounts: 200,
    //     ethereumRPCProviders: RPC_ETHEREUM_PROVIDERS,
    // })

    // console.log('Waiting for things to get really bad (20s)...')
    // await sleep(20000)

    console.log('Logging in to Chainlink nodes...')
    await chainlinkLogin(CHAINLINK_URL_NEIL,  '/tmp/neil')
    await chainlinkLogin(CHAINLINK_URL_NELLY, '/tmp/nelly')

    await sleep(5000)

    console.log('Adding job specs to Chainlink nodes...')
    let jobSpecNeil  = await addJobSpec(CHAINLINK_URL_NEIL,  makeJobSpecFluxMonitor(fluxAggregator.address, EXTERNAL_ADAPTER_INTERNAL_URL), '/tmp/neil')
    let jobSpecNelly = await addJobSpec(CHAINLINK_URL_NELLY, makeJobSpecFluxMonitor(fluxAggregator.address, EXTERNAL_ADAPTER_INTERNAL_URL), '/tmp/nelly')

    await successfulJobRuns(CHAINLINK_URL_NEIL,  jobSpecNeil.id,  1, '/tmp/neil')
    await successfulJobRuns(CHAINLINK_URL_NELLY, jobSpecNelly.id, 1, '/tmp/nelly')
    let answer = await fluxAggregator.latestAnswer()
    console.log('answer ~>', answer)
    console.log('type answer ~>', typeof answer)
    // assert(answer)
}

main()

