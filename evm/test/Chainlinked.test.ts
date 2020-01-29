import { helpers as h } from '@chainlink/eth-test-helpers'
import { ChainlinkedFactory } from '../src/generated/ChainlinkedFactory'

const chainlinkedFactory = new ChainlinkedFactory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    h.checkPublicABI(chainlinkedFactory, [])
  })
})
