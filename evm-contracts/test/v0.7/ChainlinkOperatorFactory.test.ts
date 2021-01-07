import { contract, setup, helpers } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ContractReceipt } from 'ethers/contract'
import { Operator__factory } from '../../ethers/v0.7/factories/Operator__factory'
import { ChainlinkOperatorFactory__factory } from '../../ethers/v0.7/factories/ChainlinkOperatorFactory__factory'

const linkTokenFactory = new contract.LinkToken__factory()
const operatorGeneratorFactory = new ChainlinkOperatorFactory__factory()
const operatorFactory = new Operator__factory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('ChainlinkOperatorFactory', () => {
  let link: contract.Instance<contract.LinkToken__factory>
  let operatorGenerator: contract.Instance<ChainlinkOperatorFactory__factory>
  let operator: contract.Instance<Operator__factory>

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
      const tx = await operatorGenerator.connect(roles.oracleNode).fallback()

      receipt = await tx.wait()
    })

    it('emits an event', async () => {
      const emittedOwner = helpers.evmWordToAddress(
        receipt.logs?.[0].topics?.[2],
      )
      assert.equal(emittedOwner, roles.oracleNode.address)
    })

    it('sets the correct owner', async () => {
      const emittedAddress = helpers.evmWordToAddress(
        receipt.logs?.[0].topics?.[1],
      )

      operator = await operatorFactory
        .connect(roles.defaultAccount)
        .attach(emittedAddress)
      const ownerString = await operator.owner()
      assert.equal(ownerString, roles.oracleNode.address)
    })
  })
})
