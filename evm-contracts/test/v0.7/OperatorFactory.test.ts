import { contract, setup, helpers, matchers } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { utils } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { Operator__factory } from '../../ethers/v0.7/factories/Operator__factory'
import { AuthorizedForwarder__factory } from '../../ethers/v0.7/factories/AuthorizedForwarder__factory'
import { OperatorFactory__factory } from '../../ethers/v0.7/factories/OperatorFactory__factory'

const linkTokenFactory = new contract.LinkToken__factory()
const operatorGeneratorFactory = new OperatorFactory__factory()
const operatorFactory = new Operator__factory()
const forwarderFactory = new AuthorizedForwarder__factory()

let roles: setup.Roles
const provider = setup.provider()

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('OperatorFactory', () => {
  let link: contract.Instance<contract.LinkToken__factory>
  let operatorGenerator: contract.Instance<OperatorFactory__factory>
  let operator: contract.Instance<Operator__factory>
  let forwarder: contract.Instance<AuthorizedForwarder__factory>
  let receipt: ContractReceipt
  let emittedOperator: string
  let emittedForwarder: string

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operatorGenerator = await operatorGeneratorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(operatorGenerator, [
      'created',
      'deployNewOperator',
      'deployNewOperatorAndForwarder',
      'deployNewForwarder',
      'deployNewForwarderAndTransferOwnership',
      'getChainlinkToken',
    ])
  })

  describe('#deployNewOperator', () => {
    beforeEach(async () => {
      const tx = await operatorGenerator
        .connect(roles.oracleNode)
        .deployNewOperator()

      receipt = await tx.wait()
      emittedOperator = helpers.evmWordToAddress(receipt.logs?.[0].topics?.[1])
    })

    it('emits an event', async () => {
      assert.equal(receipt?.events?.[0]?.event, 'OperatorCreated')
      assert.equal(emittedOperator, receipt.events?.[0].args?.[0])
      assert.equal(roles.oracleNode.address, receipt.events?.[0].args?.[1])
      assert.equal(roles.oracleNode.address, receipt.events?.[0].args?.[2])
    })

    it('sets the correct owner', async () => {
      operator = await operatorFactory
        .connect(roles.defaultAccount)
        .attach(emittedOperator)
      const ownerString = await operator.owner()
      assert.equal(ownerString, roles.oracleNode.address)
    })

    it('records that it deployed that address', async () => {
      assert.isTrue(await operatorGenerator.created(emittedOperator))
    })
  })

  describe('#deployNewOperatorAndForwarder', () => {
    beforeEach(async () => {
      const tx = await operatorGenerator
        .connect(roles.oracleNode)
        .deployNewOperatorAndForwarder()

      receipt = await tx.wait()
      emittedOperator = helpers.evmWordToAddress(receipt.logs?.[0].topics?.[1])
      emittedForwarder = helpers.evmWordToAddress(receipt.logs?.[1].topics?.[1])
    })

    it('emits an event recording that the operator was deployed', async () => {
      assert.equal(roles.oracleNode.address, receipt.events?.[0].args?.[1])
      assert.equal(receipt?.events?.[0]?.event, 'OperatorCreated')
      assert.equal(receipt?.events?.[0]?.args?.[0], emittedOperator)
      assert.equal(receipt?.events?.[0]?.args?.[1], roles.oracleNode.address)
      assert.equal(receipt?.events?.[0]?.args?.[2], roles.oracleNode.address)
    })

    it('emits an event recording that the forwarder was deployed', async () => {
      assert.equal(roles.oracleNode.address, receipt.events?.[0].args?.[1])
      assert.equal(receipt?.events?.[1]?.event, 'AuthorizedForwarderCreated')
      assert.equal(receipt?.events?.[1]?.args?.[0], emittedForwarder)
      assert.equal(receipt?.events?.[1]?.args?.[1], emittedOperator)
      assert.equal(receipt?.events?.[1]?.args?.[2], roles.oracleNode.address)
    })

    it('sets the correct owner on the operator', async () => {
      operator = await operatorFactory
        .connect(roles.defaultAccount)
        .attach(receipt?.events?.[0]?.args?.[0])
      assert.equal(roles.oracleNode.address, await operator.owner())
    })

    it('sets the operator as the owner of the forwarder', async () => {
      forwarder = await forwarderFactory
        .connect(roles.defaultAccount)
        .attach(receipt?.events?.[1]?.args?.[0])
      const operatorAddress = receipt?.events?.[0]?.args?.[0]
      assert.equal(operatorAddress, await forwarder.owner())
    })

    it('records that it deployed that address', async () => {
      assert.isTrue(await operatorGenerator.created(emittedOperator))
      assert.isTrue(await operatorGenerator.created(emittedForwarder))
    })
  })

  describe('#deployNewForwarder', () => {
    beforeEach(async () => {
      const tx = await operatorGenerator
        .connect(roles.oracleNode)
        .deployNewForwarder()

      receipt = await tx.wait()
      emittedForwarder = receipt.events?.[0].args?.[0]
    })

    it('emits an event', async () => {
      assert.equal(receipt?.events?.[0]?.event, 'AuthorizedForwarderCreated')
      assert.equal(roles.oracleNode.address, receipt.events?.[0].args?.[1]) // owner
      assert.equal(roles.oracleNode.address, receipt.events?.[0].args?.[2]) // sender
    })

    it('sets the caller as the owner', async () => {
      forwarder = await forwarderFactory
        .connect(roles.defaultAccount)
        .attach(emittedForwarder)
      const ownerString = await forwarder.owner()
      assert.equal(ownerString, roles.oracleNode.address)
    })

    it('records that it deployed that address', async () => {
      assert.isTrue(await operatorGenerator.created(emittedForwarder))
    })
  })

  describe('#deployNewForwarderAndTransferOwnership', () => {
    const message = '0x42'

    beforeEach(async () => {
      const tx = await operatorGenerator
        .connect(roles.oracleNode)
        .deployNewForwarderAndTransferOwnership(roles.stranger.address, message)
      receipt = await tx.wait()

      emittedForwarder = helpers.evmWordToAddress(receipt.logs?.[2].topics?.[1])
    })

    it('emits an event', async () => {
      assert.equal(receipt?.events?.[2]?.event, 'AuthorizedForwarderCreated')
      assert.equal(roles.oracleNode.address, receipt.events?.[2].args?.[1]) // owner
      assert.equal(roles.oracleNode.address, receipt.events?.[2].args?.[2]) // sender
    })

    it('sets the caller as the owner', async () => {
      forwarder = await forwarderFactory
        .connect(roles.defaultAccount)
        .attach(emittedForwarder)
      const ownerString = await forwarder.owner()
      assert.equal(ownerString, roles.oracleNode.address)
    })

    it('proposes a transfer to the recipient', async () => {
      const emittedOwner = helpers.evmWordToAddress(
        receipt.logs?.[0].topics?.[1],
      )
      assert.equal(emittedOwner, roles.oracleNode.address)
      const emittedRecipient = helpers.evmWordToAddress(
        receipt.logs?.[0].topics?.[2],
      )
      assert.equal(emittedRecipient, roles.stranger.address)
    })

    it('proposes a transfer to the recipient with the specified message', async () => {
      const emittedOwner = helpers.evmWordToAddress(
        receipt.logs?.[1].topics?.[1],
      )
      assert.equal(emittedOwner, roles.oracleNode.address)
      const emittedRecipient = helpers.evmWordToAddress(
        receipt.logs?.[1].topics?.[2],
      )
      assert.equal(emittedRecipient, roles.stranger.address)

      const encodedMessage = utils.defaultAbiCoder.encode(['bytes'], [message])
      assert.equal(receipt?.logs?.[1]?.data, encodedMessage)
    })

    it('records that it deployed that address', async () => {
      assert.isTrue(await operatorGenerator.created(emittedForwarder))
    })
  })
})
