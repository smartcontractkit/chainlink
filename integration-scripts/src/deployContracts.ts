import { LinkTokenFactory } from 'chainlink/dist/src/generated/LinkTokenFactory'
import { OracleFactory } from './generated/OracleFactory'
import {
  createProvider,
  DEVNET_ADDRESS,
  registerPromiseHandler,
  getArgs,
  deployContract,
} from './common'
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

  const linkToken = await deployContract({
    Factory: LinkTokenFactory,
    name: 'LinkToken',
    signer,
  })

  const oracle = await deployContract(
    { Factory: OracleFactory, name: 'Oracle', signer },
    linkToken.address,
  )
  await oracle.setFulfillmentPermission(chainlinkNodeAddress, true)
  console.log(`Deployed Oracle at: ${oracle.address}`)

  await deployContract({ Factory: EthLogFactory, name: 'EthLog', signer })

  await deployContract(
    { Factory: RunLogFactory, name: 'RunLog', signer },
    linkToken.address,
    oracle.address,
  )
}
