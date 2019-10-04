import { checkPublicABI } from '../src/helpersV2'
import { ChainlinkedFactory } from 'contracts/ChainlinkedFactory'

const chainlinkedFactory = new ChainlinkedFactory()

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    checkPublicABI(chainlinkedFactory, [])
  })
})
