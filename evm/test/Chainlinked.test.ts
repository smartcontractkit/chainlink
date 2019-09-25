import { checkPublicABI } from '../src/helpers'
const Chainlinked = artifacts.require('Chainlinked.sol')

contract('Chainlinked', () => {
  it('has a limited public interface', () => {
    checkPublicABI(Chainlinked, [])
  })
})
