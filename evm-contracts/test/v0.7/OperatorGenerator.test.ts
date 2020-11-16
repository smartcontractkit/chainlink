import { contract, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ContractReceipt } from 'ethers/contract'
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

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operatorGenerator = await operatorGeneratorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#createOperator', () => {
    let receipt: ContractReceipt

    beforeEach(async () => {
      const tx = await operatorGenerator
        .connect(roles.oracleNode)
        .createOperator()

      receipt = await tx.wait()
    })

    it('emits an event', async () => {
      const event = receipt.events?.[0]
      assert.equal(event?.event, 'OperatorCreated')
      assert.equal(event?.args?.[1], roles.oracleNode.address)
    })

    it('sets the correct owner', async () => {
      const args = receipt.events?.[0].args

      operator = await operatorFactory
        .connect(roles.defaultAccount)
        .attach(args?.[0])
      const ownerString = await operator.owner()
      assert.equal(ownerString, roles.oracleNode.address)
    })
  })
})
