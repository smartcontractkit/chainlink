import { OracleFactory } from '@chainlink/contracts/ethers/v0.4/OracleFactory'
import {
  createProvider,
  deployContract,
  DEVNET_ADDRESS,
  getArgs,
  registerPromiseHandler,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'
import { EthLogFactory } from './generated/EthLogFactory'
import { RunLogFactory } from './generated/RunLogFactory'

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
    { Factory: OracleFactory, name: 'Oracle', signer },
    linkToken.address,
  )
  await oracle.setFulfillmentPermission(chainlinkNodeAddress, true)

  await deployContract({ Factory: EthLogFactory, name: 'EthLog', signer })

  await deployContract(
    { Factory: RunLogFactory, name: 'RunLog', signer },
    linkToken.address,
    oracle.address,
  )
}
