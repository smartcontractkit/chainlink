import { ethers } from 'hardhat'
import { Contract } from 'ethers'
import { assert, expect } from 'chai'
import { CronUpkeepFactory } from '../../typechain/CronUpkeepFactory'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { reset } from '../test-helpers/helpers'

let cronExternalLib: Contract
let factory: CronUpkeepFactory

let admin: SignerWithAddress
let owner: SignerWithAddress

describe('CronUpkeepFactory', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    admin = accounts[0]
    owner = accounts[1]
    const cronExternalFactory = await ethers.getContractFactory(
      'src/v0.8/libraries/external/Cron.sol:Cron',
      admin,
    )
    cronExternalLib = await cronExternalFactory.deploy()
    const cronUpkeepFactoryFactory = await ethers.getContractFactory(
      'CronUpkeepFactory',
      {
        signer: admin,
        libraries: {
          Cron: cronExternalLib.address,
        },
      },
    )
    factory = await cronUpkeepFactoryFactory.deploy()
  })

  afterEach(async () => {
    await reset()
  })

  describe('constructor()', () => {
    it('deploys a delegate contract', async () => {
      assert.notEqual(
        await factory.cronDelegateAddress(),
        ethers.constants.AddressZero,
      )
    })
  })

  describe('newCronUpkeep()', () => {
    it('emits an event', async () => {
      await expect(factory.connect(owner).newCronUpkeep()).to.emit(
        factory,
        'NewCronUpkeepCreated',
      )
    })
    it('sets the deployer as the owner', async () => {
      const response = await factory.connect(owner).newCronUpkeep()
      const { events } = await response.wait()
      if (!events) {
        assert.fail('no events emitted')
      }
      const upkeepAddress = events[0].args?.upkeep
      const cronUpkeepFactory = await ethers.getContractFactory('CronUpkeep', {
        libraries: { Cron: cronExternalLib.address },
      })
      assert(
        await cronUpkeepFactory.attach(upkeepAddress).owner(),
        owner.address,
      )
    })
  })
})
