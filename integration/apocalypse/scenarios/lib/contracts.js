const assert = require('assert')
const ethers = require('ethers')
const { contract } = require('@chainlink/test-helpers')
const oracleJson = require('../../../../evm-contracts/abi/v0.6/Oracle.json')
const fluxAggregatorJson = require('../../../../evm-contracts/abi/v0.6/FluxAggregator.json')

module.exports = {
  deployDirectRequestContracts,
  deployFluxMonitorContracts,
}

async function deployDirectRequestContracts(carol) {
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
    paymentAmount: '100', // LINK-sats
    timeout: 300, // seconds
  })
  console.log('  - Flux Aggregator:', fluxAggregator.address)

  console.log('Funding FluxAggregator...')
  let resp = await (
    await linkToken.transfer(fluxAggregator.address, '100000000000000000000')
  ).wait()
  assert(resp.status === 1)

  console.log('Calling .updateAvailableFunds...')
  resp = await (await fluxAggregator.updateAvailableFunds()).wait()
  assert(resp.status === 1)

  console.log('Adding oracles:', oracles)
  let resp3 = await (
    await fluxAggregator.addOracles(oracles, oracles, 1, 2, 1)
  ).wait()
  assert(resp3.status === 1)

  console.log('Calling .updateFutureRounds...')
  let resp2 = await (
    await fluxAggregator.updateFutureRounds(
      '1000000000000000000',
      1,
      2,
      1,
      3000,
    )
  ).wait()
  assert(resp2.status === 1)

  return {
    linkToken,
    fluxAggregator,
  }
}

async function deployContract({ Factory, name, signer }, ...deployArgs) {
  const contractFactory = new Factory(signer)
  const contract = await contractFactory.deploy(...deployArgs, {
    gasPrice: ethers.utils.parseUnits('50', 'gwei'),
    gasLimit: 10 * 1000 * 1000, // 10 million
  })
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
  let oracleFactory = new ethers.ContractFactory(
    oracleJson.compilerOutput.abi,
    oracleJson.compilerOutput.evm.bytecode,
    wallet,
  )
  let oracle = await oracleFactory.deploy(linkTokenAddress, {
    gasPrice: ethers.utils.parseUnits('50', 'gwei'),
  })
  await oracle.deployed()
  return oracle
}

async function deployFluxAggregator(
  wallet,
  nonce,
  { linkTokenAddress, paymentAmount, timeout },
) {
  let description = ethers.utils.formatBytes32String('xyzzy')
  let fluxAggregatorFactory = new ethers.ContractFactory(
    fluxAggregatorJson.compilerOutput.abi,
    fluxAggregatorJson.compilerOutput.evm.bytecode.object,
    wallet,
  )
  let fluxAggregator = await fluxAggregatorFactory.deploy(
    linkTokenAddress,
    paymentAmount,
    timeout,
    2,
    description,
    {
      gasPrice: ethers.utils.parseUnits('50', 'gwei'),
    },
  )
  await fluxAggregator.deployed()
  return fluxAggregator
}
