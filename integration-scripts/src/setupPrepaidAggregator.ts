import { PrepaidAggregatorFactory } from '@chainlink/contracts/ethers/v0.5/PrepaidAggregatorFactory'
import { ethers } from 'ethers'
import {
  createProvider,
  deployContract,
  DEVNET_ADDRESS,
  getArgs,
  registerPromiseHandler,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'

async function main() {
  registerPromiseHandler()
  await deployContracts()
}

// TESTING ONLY
// const privKey =
//   'e08e067c08a8b145704cad7cd2129a8fbb753cd513b53cca7ff89dca724ef0a9'
const oracle2 = '0x6c1676F8492119C16eA62E07351E0094e93c6667'

async function deployContracts() {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  const linkToken = await deployLinkTokenContract()

  const prepaidAggregator = await deployContract(
    { Factory: PrepaidAggregatorFactory, signer, name: 'PrepaidAggregator' },
    linkToken.address,
    1,
    3600,
    1,
    ethers.utils.formatBytes32String('ETH/USD'),
  )
  const { CHAINLINK_NODE_ADDRESS } = getArgs(['CHAINLINK_NODE_ADDRESS'])
  const tx1 = await prepaidAggregator.addOracle(CHAINLINK_NODE_ADDRESS, 1, 1, 0)
  await tx1.wait()
  const tx2 = await prepaidAggregator.addOracle(oracle2, 2, 2, 0)
  await tx2.wait()
}

main()
