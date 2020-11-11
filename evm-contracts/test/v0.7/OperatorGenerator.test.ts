import { contract, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { OperatorFactory } from '../../ethers/v0.7/OperatorFactory'
import { OperatorGeneratorFactory } from '../../ethers/v0.7/OperatorGeneratorFactory'

const linkTokenFactory = new contract.LinkTokenFactory()
const operatorGeneratorFactory = new OperatorGeneratorFactory()
const operatorFactory = new OperatorFactory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('OperatorGenerator', () => {
  let link: contract.Instance<contract.LinkTokenFactory>
  let operatorGenerator: contract.Instance<OperatorGeneratorFactory>
  let operator: contract.Instance<OperatorFactory>

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operatorGenerator = await operatorGeneratorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  it('creates a new Operator', async () => {
    const tx = await operatorGenerator
      .connect(roles.oracleNode)
      .createOperator()

    const receipt = await tx.wait()
    const args = receipt.events?.[0].args

    operator = await operatorFactory
      .connect(roles.defaultAccount)
      .attach(args?.[0])

    const ownerString = await operator.owner()

    assert.equal(ownerString, roles.oracleNode.address)
  })
})
