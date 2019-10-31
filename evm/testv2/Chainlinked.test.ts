import { checkPublicABI } from '../src/helpersV2'
import { ChainlinkedFactory } from '../src/generated/ChainlinkedFactory'

const chainlinkedFactory = new ChainlinkedFactory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    checkPublicABI(chainlinkedFactory, [])
  })
})
