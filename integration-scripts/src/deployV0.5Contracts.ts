import { LinkTokenV05Factory } from './generated/LinkTokenV05Factory'
import { CoordinatorFactory } from './generated/CoordinatorFactory'
import {
  createProvider,
  DEVNET_ADDRESS,
  registerPromiseHandler,
} from './common'

async function main() {
  registerPromiseHandler()

  await deployContracts()
}
main()

async function deployContracts() {
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
  console.log(`Deployed coordinator at: ${coordinator.address}`)
}
