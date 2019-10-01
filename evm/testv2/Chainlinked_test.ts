import { checkPublicABI } from '../src/helpersV2'
import { AbstractContract } from '../src/contract'

const Chainlinked = AbstractContract.fromArtifactName('Chainlinked')

describe('Chainlinked', () => {
  it('has a limited public interface', async () => {
    const contractFactory = Chainlinked.getContractFactory()
    checkPublicABI(contractFactory, [])
  })
})
