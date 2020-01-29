import { contract, helpers as h, providers } from '@chainlink/eth-test-helpers'
import { assert } from 'chai'
import { LinkTokenFactory } from '../src/generated/LinkTokenFactory'
import { PointerFactory } from '../src/generated/PointerFactory'

const pointerFactory = new PointerFactory()
const linkTokenFactory = new LinkTokenFactory()
const provider = providers.makeTestProvider()

let roles: h.Roles

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)

  roles = rolesAndPersonas.roles
})

describe('Pointer', () => {
  let pointer: contract.Instance<PointerFactory>
  let link: contract.Instance<LinkTokenFactory>
  const deployment = providers.useSnapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    pointer = await pointerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(pointer, ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await pointer.getAddress(), link.address)
    })
  })
})
