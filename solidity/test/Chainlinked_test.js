import { checkPublicABI } from './support/helpers'

contract('Chainlinked', () => {
  const sourcePath = 'Chainlinked.sol'

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [])
  })
})
