import { FluxAggregatorFactory } from '@chainlink/contracts/ethers/v0.6/FluxAggregatorFactory'
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

  const fluxAggregator = await deployContract(
    { Factory: FluxAggregatorFactory, signer, name: 'FluxAggregator' },
    linkToken.address,
    1,
    3,
    1,
    ethers.utils.formatBytes32String('ETH/USD'),
  )
  const { CHAINLINK_NODE_ADDRESS } = getArgs(['CHAINLINK_NODE_ADDRESS'])
  await fluxAggregator.addOracle(
    CHAINLINK_NODE_ADDRESS,
    CHAINLINK_NODE_ADDRESS,
    1,
    1,
    0,
  )

  return {
    linkToken,
    fluxAggregator,
  }
}
