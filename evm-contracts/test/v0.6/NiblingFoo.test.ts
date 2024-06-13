import { NiblingFooFactory } from '../../ethers/v0.6/NiblingFooFactory'
import { contract, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'

let personas: setup.Personas
const provider = setup.provider()
const niblingFooFactory = new NiblingFooFactory()
let niblingFoo: contract.Instance<NiblingFooFactory>

beforeAll(async () => {
  personas = await setup.users(provider).then(x => x.personas)
})

describe('Foo', () => {
  const deployment = setup.snapshot(provider, async () => {
    niblingFoo = await niblingFooFactory.connect(personas.Carol).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  it('sets the inital value', async () => {
    assert((await niblingFoo.bar()).eq(utils.bigNumberify(4)))
  })
})
