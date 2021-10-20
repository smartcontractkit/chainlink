import { contract } from '@chainlink/test-helpers'
import { createProvider, deployContract, DEVNET_ADDRESS } from './common'

export async function deployLinkTokenContract(): Promise<
  contract.Instance<contract.LinkToken__factory>
> {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)
  if (process.env.LINK_TOKEN_ADDRESS) {
    console.log(
      `LinkToken already deployed at: ${process.env.LINK_TOKEN_ADDRESS}, fetching contract...`,
    )
    const factory = new contract.LinkToken__factory(signer)
    const linkToken = factory.attach(process.env.LINK_TOKEN_ADDRESS)
    console.log(`Deployed LinkToken at: ${linkToken.address}`)

    return linkToken
  }

  const linkToken = await deployContract({
    Factory: contract.LinkToken__factory,
    name: 'LinkToken',
    signer,
  })

  return linkToken
}
