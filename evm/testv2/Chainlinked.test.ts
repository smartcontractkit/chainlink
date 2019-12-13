import { checkPublicABI } from '../src/helpers'
import { ChainlinkedFactory } from '../src/generated/ChainlinkedFactory'

const chainlinkedFactory = new ChainlinkedFactory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    checkPublicABI(chainlinkedFactory, [])
  })
})
