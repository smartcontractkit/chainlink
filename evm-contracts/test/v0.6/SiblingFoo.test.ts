import { SiblingFooFactory } from '../../ethers/v0.6/SiblingFooFactory'
import { contract, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'

let personas: setup.Personas
const provider = setup.provider()
const siblingFooFactory = new SiblingFooFactory()
let siblingFoo: contract.Instance<SiblingFooFactory>

beforeAll(async () => {
  personas = await setup.users(provider).then(x => x.personas)
})

describe('Foo', () => {
  const deployment = setup.snapshot(provider, async () => {
    siblingFoo = await siblingFooFactory.connect(personas.Carol).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  it('sets the inital value', async () => {
    assert((await siblingFoo.bar()).eq(utils.bigNumberify(5)))
  })
})
