import * as h from '../src/helpersV2'
import ganache from 'ganache-core'
import { AbstractContract } from '../src/contract'
import { linkToken } from '../src/linkToken'
import { assert } from 'chai'
import { Pointer } from 'contracts/Pointer'
import { LinkTokenInterface } from 'contracts/LinkTokenInterface'
const PointerContract = AbstractContract.fromArtifactName('Pointer')
const LinkContract = AbstractContract.fromBuildArtifact(linkToken)
let roles: h.Roles
const ganacheProvider: any = ganache.provider()

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(ganacheProvider)

  roles = rolesAndPersonas.roles
})

describe('Pointer', () => {
  let contract: Pointer
  let link: LinkTokenInterface

  beforeEach(async () => {
    link = await LinkContract.deploy(roles.defaultAccount)
    contract = await PointerContract.deploy(roles.defaultAccount, [
      link.address,
    ])
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(contract, ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await contract.functions.getAddress(), link.address)
    })
  })
})
