import { ethers } from 'hardhat'
import { evmWordToAddress, publicAbi } from '../test-helpers/helpers'
import { assert } from 'chai'
import { Contract, ContractFactory, ContractReceipt } from 'ethers'
import { getUsers, Roles } from '../test-helpers/setup'

let linkTokenFactory: ContractFactory
let operatorGeneratorFactory: ContractFactory
let operatorFactory: ContractFactory
let forwarderFactory: ContractFactory

let roles: Roles

before(async () => {
  const users = await getUsers()

  roles = users.roles
  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    roles.defaultAccount,
  )
  operatorGeneratorFactory = await ethers.getContractFactory(
    'src/v0.7/OperatorFactory.sol:OperatorFactory',
    roles.defaultAccount,
  )
  operatorFactory = await ethers.getContractFactory(
    'src/v0.7/Operator.sol:Operator',
    roles.defaultAccount,
  )
  forwarderFactory = await ethers.getContractFactory(
    'src/v0.7/AuthorizedForwarder.sol:AuthorizedForwarder',
    roles.defaultAccount,
  )
})

describe('OperatorFactory', () => {
  let link: Contract
  let operatorGenerator: Contract
  let operator: Contract
  let forwarder: Contract
  let receipt: ContractReceipt
  let emittedOperator: string
  let emittedForwarder: string

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    operatorGenerator = await operatorGeneratorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address)
  })

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(operatorGenerator, [
      'created',
      'deployNewOperator',
      'deployNewOperatorAndForwarder',
      'deployNewForwarder',
      'deployNewForwarderAndTransferOwnership',
      'getChainlinkToken',
      'typeAndVersion',
    ])
  })

  describe('#typeAndVersion', () => {
    it('describes the authorized forwarder', async () => {
      assert.equal(
        await operatorGenerator.typeAndVersion(),
        'OperatorFactory 1.0.0',
      )
    })
  })

  describe('#deployNewOperator', () => {
    beforeEach(async () => {
      const tx = await operatorGenerator
        .connect(roles.oracleNode)
        .deployNewOperator()

      receipt = await tx.wait()
      emittedOperator = evmWordToAddress(receipt.logs?.[0].topics?.[1])
    })

    it('emits an event', async () => {
      assert.equal(receipt?.events?.[0]?.event, 'OperatorCreated')
      assert.equal(emittedOperator, receipt.events?.[0].args?.[0])
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[0].args?.[1],
      )
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[0].args?.[2],
      )
    })

    it('sets the correct owner', async () => {
      operator = await operatorFactory
        .connect(roles.defaultAccount)
        .attach(emittedOperator)
      const ownerString = await operator.owner()
      assert.equal(ownerString, await roles.oracleNode.getAddress())
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
      emittedOperator = evmWordToAddress(receipt.logs?.[0].topics?.[1])
      emittedForwarder = evmWordToAddress(receipt.logs?.[3].topics?.[1])
    })

    it('emits an event recording that the operator was deployed', async () => {
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[0].args?.[1],
      )
      assert.equal(receipt?.events?.[0]?.event, 'OperatorCreated')
      assert.equal(receipt?.events?.[0]?.args?.[0], emittedOperator)
      assert.equal(
        receipt?.events?.[0]?.args?.[1],
        await roles.oracleNode.getAddress(),
      )
      assert.equal(
        receipt?.events?.[0]?.args?.[2],
        await roles.oracleNode.getAddress(),
      )
    })

    it('proposes the transfer of the forwarder to the operator', async () => {
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[0].args?.[1],
      )
      assert.equal(
        receipt?.events?.[1]?.topics?.[0],
        '0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278', //OwnershipTransferRequested(address,address)
      )
      assert.equal(
        evmWordToAddress(receipt?.events?.[1]?.topics?.[1]),
        operatorGenerator.address,
      )
      assert.equal(
        evmWordToAddress(receipt?.events?.[1]?.topics?.[2]),
        emittedOperator,
      )

      assert.equal(
        receipt?.events?.[2]?.topics?.[0],
        '0x4e1e878dc28d5f040db5969163ff1acd75c44c3f655da2dde9c70bbd8e56dc7e', //OwnershipTransferRequestedWithMessage(address,address,bytes)
      )
      assert.equal(
        evmWordToAddress(receipt?.events?.[2]?.topics?.[1]),
        operatorGenerator.address,
      )
      assert.equal(
        evmWordToAddress(receipt?.events?.[2]?.topics?.[2]),
        emittedOperator,
      )
    })

    it('emits an event recording that the forwarder was deployed', async () => {
      assert.equal(receipt?.events?.[3]?.event, 'AuthorizedForwarderCreated')
      assert.equal(receipt?.events?.[3]?.args?.[0], emittedForwarder)
      assert.equal(receipt?.events?.[3]?.args?.[1], operatorGenerator.address)
      assert.equal(
        receipt?.events?.[3]?.args?.[2],
        await roles.oracleNode.getAddress(),
      )
    })

    it('sets the correct owner on the operator', async () => {
      operator = await operatorFactory
        .connect(roles.defaultAccount)
        .attach(receipt?.events?.[0]?.args?.[0])
      assert.equal(await roles.oracleNode.getAddress(), await operator.owner())
    })

    it('sets the operator as the owner of the forwarder', async () => {
      forwarder = await forwarderFactory
        .connect(roles.defaultAccount)
        .attach(emittedForwarder)
      assert.equal(operatorGenerator.address, await forwarder.owner())
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
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[0].args?.[1],
      ) // owner
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[0].args?.[2],
      ) // sender
    })

    it('sets the caller as the owner', async () => {
      forwarder = await forwarderFactory
        .connect(roles.defaultAccount)
        .attach(emittedForwarder)
      const ownerString = await forwarder.owner()
      assert.equal(ownerString, await roles.oracleNode.getAddress())
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
        .deployNewForwarderAndTransferOwnership(
          await roles.stranger.getAddress(),
          message,
        )
      receipt = await tx.wait()

      emittedForwarder = evmWordToAddress(receipt.logs?.[2].topics?.[1])
    })

    it('emits an event', async () => {
      assert.equal(receipt?.events?.[2]?.event, 'AuthorizedForwarderCreated')
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[2].args?.[1],
      ) // owner
      assert.equal(
        await roles.oracleNode.getAddress(),
        receipt.events?.[2].args?.[2],
      ) // sender
    })

    it('sets the caller as the owner', async () => {
      forwarder = await forwarderFactory
        .connect(roles.defaultAccount)
        .attach(emittedForwarder)
      const ownerString = await forwarder.owner()
      assert.equal(ownerString, await roles.oracleNode.getAddress())
    })

    it('proposes a transfer to the recipient', async () => {
      const emittedOwner = evmWordToAddress(receipt.logs?.[0].topics?.[1])
      assert.equal(emittedOwner, await roles.oracleNode.getAddress())
      const emittedRecipient = evmWordToAddress(receipt.logs?.[0].topics?.[2])
      assert.equal(emittedRecipient, await roles.stranger.getAddress())
    })

    it('proposes a transfer to the recipient with the specified message', async () => {
      const emittedOwner = evmWordToAddress(receipt.logs?.[1].topics?.[1])
      assert.equal(emittedOwner, await roles.oracleNode.getAddress())
      const emittedRecipient = evmWordToAddress(receipt.logs?.[1].topics?.[2])
      assert.equal(emittedRecipient, await roles.stranger.getAddress())

      const encodedMessage = ethers.utils.defaultAbiCoder.encode(
        ['bytes'],
        [message],
      )
      assert.equal(receipt?.logs?.[1]?.data, encodedMessage)
    })

    it('records that it deployed that address', async () => {
      assert.isTrue(await operatorGenerator.created(emittedForwarder))
    })
  })
})
