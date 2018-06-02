'use strict'

require('./support/helpers.js')

contract('Chainlinked', () => {
  const sourcePath = 'Chainlinked.sol'

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [])
  })
})
