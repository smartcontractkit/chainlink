import * as h from '../src/helpersV2'
import { assert } from 'chai'
import { PointerFactory } from '../src/generated/PointerFactory'
import { LinkTokenFactory } from '../src/generated/LinkTokenFactory'
import { Instance } from '../src/contract'
import env from '@nomiclabs/buidler'
import { EthersProviderWrapper } from '../src/provider'

const pointerFactory = new PointerFactory()
const linkTokenFactory = new LinkTokenFactory()
const provider = new EthersProviderWrapper(env.ethereum)

let roles: h.Roles

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)

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
