const fs = require('fs')
const ethers = require('ethers')
const { contract } = require('@chainlink/test-helpers')

const oracleJson = require('../../evm-contracts/abi/v0.6/Oracle.json')

async function deployContract({ Factory, name, signer }, ...deployArgs) {
    const contractFactory = new Factory(signer)
    const contract = await contractFactory.deploy(...deployArgs, { gasPrice: ethers.utils.parseUnits('50', 'gwei') })
    await contract.deployed()
    return contract
}

async function deployLinkTokenContract(wallet) {
    
    const linkToken = await deployContract({
      Factory: contract.LinkTokenFactory,
      name: 'LinkToken',
      signer: wallet,
    })
  
    return linkToken
}

async function txBombardment(wallet) {
    let tx = {
        to: '0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f',
        value: ethers.utils.parseUnits('1', 'gwei'),
        gasPrice: ethers.utils.parseUnits('100', 'gwei'),
    }
    
    let nonce = await wallet.provider.getTransactionCount('0xde554b6c292f5e5794a68dc560a537dd89d3b03e', 'pending')
    
    for (let i = 0; i < 100; i++) {
        tx.nonce = nonce + i
        wallet.sendTransaction(tx).then(tx => {
            console.log(tx)
        })
    }
}

async function main() {
    const privateKey = 'a82180a8001e2681b9feac787afaf45f1d0bb7cb61eed53f879030cca1823459'
    let wallet = new ethers.Wallet(privateKey, new ethers.providers.JsonRpcProvider('http://localhost:8545'))

    console.log('Deploying LINK token...')
    let linkToken = await deployLinkTokenContract(wallet)
    await linkToken.deployed()
    console.log('LINK token:', linkToken.address)

    console.log('Deploying Oracle contract...')
    let oracleFactory = new ethers.ContractFactory(oracleJson.compilerOutput.abi, oracleJson.compilerOutput.evm.bytecode, wallet)
    let oracle = await oracleFactory.deploy(linkToken.address, { gasPrice: ethers.utils.parseUnits('50', 'gwei') })
    await oracle.deployed()
    console.log('Oracle:', oracle.address)
    

}

main()

