import { contract, matchers, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { Pointer__factory } from '../../ethers/v0.4/factories/Pointer__factory'

const pointerFactory = new Pointer__factory()
const linkTokenFactory = new contract.LinkToken__factory()
const provider = setup.provider()

let roles: setup.Roles

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('Pointer', () => {
  let pointer: contract.Instance<Pointer__factory>
  let link: contract.Instance<contract.LinkToken__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    pointer = await pointerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(pointer, ['getAddress'])
  })

  describe('#getAddress', () => {
    it('returns the LINK token address', async () => {
      assert.equal(await pointer.getAddress(), link.address)
    })
  })
})
