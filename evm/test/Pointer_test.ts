import * as h from '../src/helpers'
const Pointer = artifacts.require('Pointer.sol')

let roles: h.Roles

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas()

  roles = rolesAndPersonas.roles
})

contract('Pointer', () => {
  let contract: any
  let link: any

  beforeEach(async () => {
    link = await h.linkContract(roles.defaultAccount)
    contract = await Pointer.new(link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Pointer, ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await contract.getAddress(), link.address)
    })
  })
})
