import { matchers } from '@chainlink/test-helpers'
import { ChainlinkedFactory } from '../../ethers/v0.4/ChainlinkedFactory'

const chainlinkedFactory = new ChainlinkedFactory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    matchers.publicAbi(chainlinkedFactory, [])
  })
})
