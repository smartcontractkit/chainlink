import * as h from '../src/helpersV2'
import * as ethers from 'ethers'
import ganache from 'ganache-core'
import { AbstractContract } from '../src/contract'
import { linkToken } from '../src/linkToken'
import { assert } from 'chai'
const Pointer = AbstractContract.fromArtifactName('Pointer')
const Link = AbstractContract.fromBuildArtifact(linkToken)
let roles: h.Roles
const ganacheProvider: any = ganache.provider()

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(ganacheProvider)

  roles = rolesAndPersonas.roles
})

describe('Pointer', () => {
  let contract: ethers.Contract
  let link: ethers.Contract

  beforeEach(async () => {
    link = await Link.deploy(roles.defaultAccount)
    contract = await Pointer.deploy(roles.defaultAccount, [link.address])
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
