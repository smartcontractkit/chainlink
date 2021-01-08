import { Oracle__factory } from '@chainlink/contracts/ethers/v0.4/factories/Oracle__factory'
import {
  createProvider,
  deployContract,
  DEVNET_ADDRESS,
  getArgs,
  registerPromiseHandler,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'
import { EthLog__factory } from './generated/factories/EthLog__factory'
import { RunLog__factory } from './generated/factories/RunLog__factory'

async function main() {
  registerPromiseHandler()
  const args = getArgs(['CHAINLINK_NODE_ADDRESS'])

  await deployContracts({ chainlinkNodeAddress: args.CHAINLINK_NODE_ADDRESS })
}
main()

interface Args {
  chainlinkNodeAddress: string
}
async function deployContracts({ chainlinkNodeAddress }: Args) {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  console.log('Deploying contracts from ' + DEVNET_ADDRESS)

  const linkToken = await deployLinkTokenContract()

  const oracle = await deployContract(
    { Factory: Oracle__factory, name: 'Oracle', signer },
    linkToken.address,
  )
  await oracle.setFulfillmentPermission(chainlinkNodeAddress, true)

  await deployContract({ Factory: EthLog__factory, name: 'EthLog', signer })

  await deployContract(
    { Factory: RunLog__factory, name: 'RunLog', signer },
    linkToken.address,
    oracle.address,
  )
}
