import { utils } from 'ethers'
import { registerPromiseHandler, createProvider, getArgs } from './common'

/**
 * This script is used to fund a chainlink node address given a developer account.
 * This is used in our integration test to fund our devnet account from a geth developer account.
 * It isnt used in our parity version of the integration test since the parity account is already funded for us.
 */
async function main() {
  registerPromiseHandler()
  const args = getArgs(['CHAINLINK_NODE_ADDRESS', 'DEVELOPER_ACCOUNT'])

  await fundAddress({
    recipient: args.CHAINLINK_NODE_ADDRESS,
    signerAddress: args.DEVELOPER_ACCOUNT,
  })
}
main()

interface Args {
  recipient: string
  signerAddress: string
}
async function fundAddress({ recipient, signerAddress }: Args) {
  const provider = createProvider()
  const signer = provider.getSigner(signerAddress)

  const tx = await signer.sendTransaction({
    to: recipient,
    value: utils.parseEther('1000'),
  })
  const receipt = await tx.wait()

  console.log(receipt)
}
