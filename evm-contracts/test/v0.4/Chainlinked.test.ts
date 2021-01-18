import { matchers } from '@chainlink/test-helpers'
import { Chainlinked__factory } from '../../ethers/v0.4/factories/Chainlinked__factory'

const chainlinkedFactory = new Chainlinked__factory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    matchers.publicAbi(chainlinkedFactory, [])
  })
})
