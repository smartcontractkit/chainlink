import { expect } from 'chai'
import { ethers, network } from 'hardhat'
import { getUsers, Roles } from '../../test-helpers/setup'
import { AutomationForwarder } from '../../../typechain/AutomationForwarder'
import { loadFixture } from '@nomicfoundation/hardhat-network-helpers'
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
let forwarder: AutomationForwarder
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

const newRegistry = ethers.Wallet.createRandom().address

const setup = async () => {
  roles = (await getUsers()).roles
  defaultAddress = await roles.defaultAccount.getAddress()
  target = await deployMockContract(roles.defaultAccount, targetABI)
  await target.deployed()
  await target.mock.handler.returns()
  await target.mock.handlerUint.returns(100)
  await target.mock.iRevert.revertsWithReason(CUSTOM_REVERT)
  const factory = await ethers.getContractFactory('AutomationForwarder')
  forwarder = await factory
    .connect(roles.defaultAccount)
    .deploy(100, target.address)
  await forwarder.deployed()
}

describe('AutomationForwarder', () => {
  beforeEach(async () => {
    await loadFixture(setup)
  })

  describe('constructor()', () => {
    it('sets the initial values', async () => {
      expect(await forwarder.getRegistry()).to.equal(defaultAddress)
      expect(await forwarder.getTarget()).to.equal(target.address)
      expect(await forwarder.getUpkeepID()).to.equal(100)
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
