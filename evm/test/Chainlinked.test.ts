import { matchers } from '@chainlink/eth-test-helpers/src'
import { ChainlinkedFactory } from '../src/generated/ChainlinkedFactory'

const chainlinkedFactory = new ChainlinkedFactory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    matchers.publicAbi(chainlinkedFactory, [])
  })
})
