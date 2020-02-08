import { CoordinatorFactory } from '@chainlink/contracts/ethers/v0.5/CoordinatorFactory'
import { MeanAggregatorFactory } from '@chainlink/contracts/ethers/v0.5/MeanAggregatorFactory'
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
main()

export async function deployContracts() {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  const linkToken = await deployLinkTokenContract()

  const coordinator = await deployContract(
    { Factory: CoordinatorFactory, signer, name: 'Coordinator' },
    linkToken.address,
  )

  const meanAggregator = await deployContract({
    Factory: MeanAggregatorFactory,
    signer,
    name: 'MeanAggregator',
  })

  const prepaidAggregator = await deployContract(
    { Factory: PrepaidAggregatorFactory, signer, name: 'PrepaidAggregator' },
    linkToken.address,
    1,
    3,
    1,
    ethers.utils.formatBytes32String('ETH/USD'),
  )
  const { CHAINLINK_NODE_ADDRESS } = getArgs(['CHAINLINK_NODE_ADDRESS'])
  await prepaidAggregator.addOracle(CHAINLINK_NODE_ADDRESS, 1, 1, 0)

  return {
    linkToken,
    coordinator,
    meanAggregator,
    prepaidAggregator,
  }
}
