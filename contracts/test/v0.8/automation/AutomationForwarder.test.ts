import { expect } from 'chai'
import { ethers, network } from 'hardhat'
import { ContractReceipt } from 'ethers'
import { getUsers, Roles } from '../../test-helpers/setup'
import { AutomationForwarder } from '../../../typechain/AutomationForwarder'
import { AutomationForwarderFactory } from '../../../typechain/AutomationForwarderFactory'
import { AutomationForwarder__factory as AutomationForwarder_Factory } from '../../../typechain/factories/AutomationForwarder__factory'
import { AutomationForwarderFactory__factory as AutomationForwarderFactory_Factory } from '../../../typechain/factories/AutomationForwarderFactory__factory'
import {
  deployMockContract,
  MockContract,
} from '@ethereum-waffle/mock-contract'

/**
 * @dev note there are two types of factories in this test: the AutomationForwarderFactory contract, which
 * deploys new instances of the AutomationForwarder, and the ethers-js javascript factory, which deploys new
 * contracts of any type from JS. Therefore, the forwarderFactoryFactory (below) is a javascript object that deploys
 * the contract factory, which deploys forwarder instances :)
 */

const NOT_AUTHORIZED_ERR = 'NotAuthorized()'
const CUSTOM_REVERT = 'this is a custom revert message'

let roles: Roles
let defaultAddress: string
let strangerAddress: string
let forwarder_factory: AutomationForwarder_Factory
let forwarder: AutomationForwarder
let forwarderFactory_factory: AutomationForwarderFactory_Factory
let forwarderFactory: AutomationForwarderFactory
let target: MockContract

const targetABI = [
  'function handler()', // 0xc80916d4
  'function handlerUint(uint256) returns(uint256)', // 0x28da6d29
  'function iRevert()', // 0xc07d0f94
]

const HANDLER = '0xc80916d4'
const HANDLER_UINT =
  '0x28da6d290000000000000000000000000000000000000000000000000000000000000001'
const HANDLER_BYTES = '0x100e0465'
const HANDLER_REVERT = '0xc07d0f94'

const randBytes = ethers.utils.randomBytes(500)
const newRegistry = ethers.Wallet.createRandom().address

function getForwarderFromDeploy(receipt: ContractReceipt) {
  const sig = 'NewForwarderDeployed(address)'
  for (const log of receipt.logs) {
    const result = forwarderFactory.interface.parseLog(log)
    if (result && result.signature === sig) {
      return forwarder_factory.attach(result.args[0])
    }
  }
  throw new Error(`couldn't find ${sig} error in tx receipt`)
}

before(async () => {
  roles = (await getUsers()).roles
  defaultAddress = await roles.defaultAccount.getAddress()
  strangerAddress = await roles.stranger.getAddress()
  forwarder_factory = await ethers.getContractFactory(
    'AutomationForwarder',
    roles.defaultAccount,
  )
  forwarderFactory_factory = await ethers.getContractFactory(
    'AutomationForwarderFactory',
    roles.defaultAccount,
  )
})

describe('AutomationForwarder', () => {
  beforeEach(async () => {
    await network.provider.send('hardhat_reset')
    target = await deployMockContract(roles.defaultAccount, targetABI)
    await target.mock.handler.returns()
    await target.mock.handlerUint.returns(100)
    await target.mock.iRevert.revertsWithReason(CUSTOM_REVERT)
    forwarder = await forwarder_factory.deploy(defaultAddress, target.address)
    forwarderFactory = await forwarderFactory_factory.deploy()
  })

  describe('constructor()', () => {
    it('sets the initial values', async () => {
      expect(await forwarder.getRegistry()).to.equal(defaultAddress)
      expect(await forwarder.getTarget()).to.equal(target.address)
    })
  })

  describe('typeAndVersion()', () => {
    it('has the correct type and version', async () => {
      expect(await forwarder.typeAndVersion()).to.equal(
        'AutomationForwarder 1.0.0',
      )
    })
  })

  describe('forward()', () => {
    const gas = 100_000
    it('is only callable by the registry', async () => {
      await expect(
        forwarder.connect(roles.stranger).forward(gas, HANDLER),
      ).to.be.revertedWith(NOT_AUTHORIZED_ERR)
      await forwarder.connect(roles.defaultAccount).forward(gas, HANDLER)
    })

    it('forwards the call to the target', async () => {
      await forwarder.connect(roles.defaultAccount).forward(gas, HANDLER)
      await forwarder.connect(roles.defaultAccount).forward(gas, HANDLER_UINT)
      await forwarder.connect(roles.defaultAccount).forward(gas, HANDLER_BYTES)
    })

    it('returns the success value of the target call', async () => {
      const result = await forwarder
        .connect(roles.defaultAccount)
        .callStatic.forward(gas, HANDLER)
      expect(result).to.be.true

      const result2 = await forwarder
        .connect(roles.defaultAccount)
        .callStatic.forward(gas, HANDLER_UINT)
      expect(result2).to.be.true

      const result3 = await forwarder
        .connect(roles.defaultAccount)
        .callStatic.forward(gas, HANDLER_REVERT)
      expect(result3).to.be.false
    })

    it('reverts if too little gas is supplied', async () => {
      await expect(
        forwarder
          .connect(roles.defaultAccount)
          .forward(100_000, HANDLER, { gasLimit: 99_999 }),
      ).to.be.reverted
    })
  })

  describe('updateRegistry()', () => {
    it('is only callable by the existing registry', async () => {
      await expect(
        forwarder.connect(roles.stranger).updateRegistry(newRegistry),
      ).to.be.revertedWith(NOT_AUTHORIZED_ERR)
      await forwarder.connect(roles.defaultAccount).updateRegistry(newRegistry)
    })

    it('is updates the registry', async () => {
      expect(await forwarder.getRegistry()).to.equal(defaultAddress)
      await forwarder.connect(roles.defaultAccount).updateRegistry(newRegistry)
      expect(await forwarder.getRegistry()).to.equal(newRegistry)
    })
  })
})

describe('AutomationForwarderFactory', () => {
  describe('typeAndVersion()', () => {
    it('has the correct type and version', async () => {
      expect(await forwarderFactory.typeAndVersion()).to.equal(
        'AutomationForwarderFactory 1.0.0',
      )
    })
  })

  describe('deploy', () => {
    it('is callable by anyone', async () => {
      await forwarderFactory
        .connect(roles.defaultAccount)
        .deploy(target.address)
      await forwarderFactory.connect(roles.stranger).deploy(target.address)
    })

    it('sets the caller as the registry', async () => {
      const tx = await forwarderFactory
        .connect(roles.stranger)
        .deploy(target.address)
      const forwarder = getForwarderFromDeploy(await tx.wait())
      expect(await forwarder.getRegistry()).to.equal(strangerAddress)
    })
  })
})
