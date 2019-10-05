import { utils } from 'ethers'
import {
  DEVNET_ADDRESS,
  registerPromiseHandler,
  createProvider,
  getArgs,
} from './common'

async function main() {
  registerPromiseHandler()
  const args = getArgs(['CHAINLINK_NODE_ADDRESS'])

  await fundAddress({ recipient: args.CHAINLINK_NODE_ADDRESS })
}
main()

interface Args {
  recipient: string
}
async function fundAddress({ recipient }: Args) {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  const tx = await signer.sendTransaction({
    to: recipient,
    value: utils.parseEther('1000'),
  })
  const receipt = await tx.wait()

  console.log(receipt)
}
