import { generated as chainlink } from 'chainlink'
import {
  createProvider,
  DEVNET_ADDRESS,
  registerPromiseHandler,
  getArgs,
  deployContract,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'
import { EthLogFactory, RunLogFactory } from './generated'

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

  const linkToken = await deployLinkTokenContract()

  const oracle = await deployContract(
    { Factory: chainlink.OracleFactory, name: 'Oracle', signer },
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
