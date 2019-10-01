import * as h from '../src/helpersV2'
import ganache from 'ganache-core'
import { assert } from 'chai'
import { PointerFactory } from 'contracts/PointerFactory'
import { LinkTokenFactory } from 'contracts/LinkTokenFactory'
import { Instance } from 'src/contract'
const pointerFactory = new PointerFactory()
const linkTokenFactory = new LinkTokenFactory()

let roles: h.Roles
const ganacheProvider: any = ganache.provider()

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(ganacheProvider)

  roles = rolesAndPersonas.roles
})

describe('Pointer', () => {
  let contract: Instance<PointerFactory>
  let link: Instance<LinkTokenFactory>

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    contract = await pointerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(contract, ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await contract.getAddress(), link.address)
    })
  })
})
