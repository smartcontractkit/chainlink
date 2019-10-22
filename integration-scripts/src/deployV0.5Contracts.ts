import { LinkTokenV05Factory } from './generated/LinkTokenV05Factory'
import { CoordinatorFactory } from './generated/CoordinatorFactory'
import { MeanAggregatorFactory } from './generated/MeanAggregatorFactory'
import {
  registerPromiseHandler,
  DEVNET_ADDRESS,
  createProvider,
} from './common'

async function main() {
  registerPromiseHandler()

  await deployContracts()
}
main()

// export async function deployContracts(provider: ethers.providers.JsonRpcProvider, DEVNET_ADDRESS: string) {
export async function deployContracts() {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  // deploy LINK token
  const linkTokenFactory = new LinkTokenV05Factory(signer)
  const linkToken = await linkTokenFactory.deploy()
  await linkToken.deployed()
  console.log(`Deployed LinkToken at: ${linkToken.address}`)

  // deploy Coordinator
  const coordinatorFactory = new CoordinatorFactory(signer)
  const coordinator = await coordinatorFactory.deploy(linkToken.address)
  await coordinator.deployed()
  console.log(`Deployed Coordinator at: ${coordinator.address}`)

  // deploy MeanAggregator
  const meanAggregatorFactory = new MeanAggregatorFactory(signer)
  const meanAggregator = await meanAggregatorFactory.deploy()
  await meanAggregator.deployed()
  console.log(`Deployed MeanAggregator at: ${meanAggregator.address}`)

  return {
    linkToken,
    coordinator,
    meanAggregator,
  }
}
