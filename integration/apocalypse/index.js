const fs = require('fs')
const ethers = require('ethers')
const _ = require('lodash')
const { contract } = require('@chainlink/test-helpers')

const oracleJson = require('../../evm-contracts/abi/v0.6/Oracle.json')

async function deployContract({ Factory, name, signer }, ...deployArgs) {
    const contractFactory = new Factory(signer)
    const contract = await contractFactory.deploy(...deployArgs, { gasPrice: ethers.utils.parseUnits('50', 'gwei') })
    await contract.deployed()
    return contract
}

async function deployLinkTokenContract(wallet, nonce) {
    
    const linkToken = await deployContract({
        Factory: contract.LinkTokenFactory,
        name: 'LinkToken',
        signer: wallet,
        nonce: nonce,
    })
  
    return linkToken
}

async function txBombardment() {
    let tx = {
        to: '0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f',
        value: ethers.utils.parseUnits('1', 'gwei'),
        gasPrice: ethers.utils.parseUnits('100', 'gwei'),
    }
    
    //let nonce = await wallets[0].provider.getTransactionCount('0xde554b6c292f5e5794a68dc560a537dd89d3b03e', 'pending')
    //for (let i = 0; i < 333; i++) {
    //    for (let wallet of wallets) {
    //        tx.nonce = nonce
    //        nonce++
    //        wallet.sendTransaction(tx).then(tx => {
    //            console.log(tx)
    //        })
    //    }
    //}

    let providers = [
        new ethers.providers.JsonRpcProvider('http://localhost:8545'),
        new ethers.providers.JsonRpcProvider('http://localhost:18545'),
        //new ethers.providers.JsonRpcProvider('http://localhost:28545'),
    ]

    let wallets = Array(100).fill(null).map((_, i) => ethers.Wallet.createRandom().connect(providers[i % providers.length]))
    let nonces = Array(100).fill(null).map(_ => 0)

    const moneybagsPrivateKey = 'a82180a8001e2681b9feac787afaf45f1d0bb7cb61eed53f879030cca1823459'
    const moneybagsWallet = new ethers.Wallet(moneybagsPrivateKey, new ethers.providers.JsonRpcProvider('http://localhost:8545'))
    let moneybagsNonce = await wallets[0].provider.getTransactionCount('0xde554b6c292f5e5794a68dc560a537dd89d3b03e', 'pending')

    let groups = _.chunk(wallets, 10)
    for (let group of groups) {
        console.log(`Funding ${group}...`)
        let batch = group.map(async (wallet) => (await moneybagsWallet.sendTransaction({
            to: wallet.address,
            value: ethers.utils.parseUnits('1', 'ether'),
            gasPrice: ethers.utils.parseUnits('20', 'gwei'),
            nonce: moneybagsNonce++,
        })).wait())
        console.log(batch)

        await Promise.all(batch)
    }

    console.log('Initiating tx tornado...')
    for (let i = 0; i < 1000; i++) {
        let wallet = wallets[i % providers.length]
        tx.nonce = nonces[i % providers.length]
        nonces[i % providers.length]++
        //nonce++
        wallet.sendTransaction(tx).then(tx => {
            console.log(tx)
        })
    }
}

async function main() {
    const privateKey = 'a82180a8001e2681b9feac787afaf45f1d0bb7cb61eed53f879030cca1823459'
    let wallet1 = new ethers.Wallet(privateKey, new ethers.providers.JsonRpcProvider('http://localhost:8545'))
    let wallet2 = new ethers.Wallet(privateKey, new ethers.providers.JsonRpcProvider('http://localhost:18545'))
    let wallet3 = new ethers.Wallet(privateKey, new ethers.providers.JsonRpcProvider('http://localhost:28545'))
        
    //console.log('Deploying LINK token...')
    //let linkToken = await deployLinkTokenContract(wallet)
    //await linkToken.deployed()
    //console.log('LINK token:', linkToken.address)
    //
    //console.log('Deploying Oracle contract...')
    //let oracleFactory = new ethers.ContractFactory(oracleJson.compilerOutput.abi, oracleJson.compilerOutput.evm.bytecode, wallet)
    //let oracle = await oracleFactory.deploy(linkToken.address, { gasPrice: ethers.utils.parseUnits('50', 'gwei') })
    //await oracle.deployed()
    //console.log('Oracle:', oracle.address)
    
    console.log('Sending txs...')
    await txBombardment([wallet1, wallet2, wallet3])
    console.log('done')
}

main()

