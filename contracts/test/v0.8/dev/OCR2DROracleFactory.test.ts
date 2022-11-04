import { ethers } from 'hardhat'
import { evmWordToAddress, publicAbi } from '../../test-helpers/helpers'
import { assert } from 'chai'
import { Contract, ContractFactory, ContractReceipt } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let oracleGeneratorFactory: ContractFactory
let oracleFactory: ContractFactory
let roles: Roles

before(async () => {
  const users = await getUsers()

  roles = users.roles
  oracleGeneratorFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ocr2dr/OCR2DROracleFactory.sol:OCR2DROracleFactory',
    roles.defaultAccount,
  )
  oracleFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ocr2dr/OCR2DROracle.sol:OCR2DROracle',
    roles.defaultAccount,
  )
})

describe('OCR2DROracleFactory', () => {
  let factory: Contract
  let oracle: Contract
  let receipt: ContractReceipt
  let emittedOracle: string

  beforeEach(async () => {
    factory = await oracleGeneratorFactory
      .connect(roles.defaultAccount)
      .deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(factory, ['created', 'deployNewOracle', 'typeAndVersion'])
  })

  describe('#typeAndVersion', () => {
    it('describes the authorized forwarder', async () => {
      assert.equal(await factory.typeAndVersion(), 'OCR2DROracleFactory 0.0.0')
    })
  })

  describe('#deployNewOracle', () => {
    beforeEach(async () => {
      const tx = await factory.connect(roles.oracleNode).deployNewOracle()

      receipt = await tx.wait()
      emittedOracle = evmWordToAddress(receipt.logs?.[0].topics?.[1])
    })

    it('emits an event', async () => {
      assert.equal(receipt?.events?.[0]?.event, 'OracleCreated')
      assert.equal(emittedOracle, receipt.events?.[0].args?.[0])
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
      oracle = await oracleFactory
        .connect(roles.defaultAccount)
        .attach(emittedOracle)
      const ownerString = await oracle.owner()
      assert.equal(ownerString, await factory.address)
    })

    it('records that it deployed that address', async () => {
      assert.isTrue(await factory.created(emittedOracle))
      assert.isFalse(
        await factory.created(await roles.oracleNode1.getAddress()),
      )
    })
  })
})
