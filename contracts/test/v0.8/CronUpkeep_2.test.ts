import moment from 'moment'
import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { CronUpkeepTestHelper } from '../../typechain/CronUpkeepTestHelper'
import { CronUpkeepDelegate } from '../../typechain/CronUpkeepDelegate'
import { CronUpkeepFactory } from '../../typechain/CronUpkeepFactory'
import { CronUpkeepTestHelper__factory as CronUpkeepTestHelperFactory } from '../../typechain/factories/CronUpkeepTestHelper__factory'
import { CronReceiver } from '../../typechain/CronReceiver'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { validCrons } from '../test-helpers/fixtures'
import * as h from '../test-helpers/helpers'

const { utils } = ethers
const { AddressZero } = ethers.constants

const OWNABLE_ERR = 'Only callable by owner'
const CRON_NOT_FOUND_ERR = 'CronJobIDNotFound'

let cron: CronUpkeepTestHelper
let cronFactory: CronUpkeepTestHelperFactory // the typechain factory that deploys cron upkeep contracts
let cronFactoryContract: CronUpkeepFactory // the cron factory contract
let cronDelegate: CronUpkeepDelegate
let cronReceiver1: CronReceiver
let cronReceiver2: CronReceiver

let admin: SignerWithAddress
let owner: SignerWithAddress
let stranger: SignerWithAddress

const basicCronString = '0 * * * *'

let handler1Sig: string
let handler2Sig: string
let basicSpec: string

async function assertJobIDsEqual(expected: number[]) {
  const ids = (await cron.getActiveCronJobIDs()).map((n) => n.toNumber())
  assert.deepEqual(ids.sort(), expected.sort())
}

async function createBasicCron() {
  return await cron.createCronJobFromEncodedSpec(
    cronReceiver1.address,
    handler1Sig,
    basicSpec,
  )
}

describe('CronUpkeep 2/2', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    admin = accounts[0]
    owner = accounts[1]
    stranger = accounts[2]
    const crFactory = await ethers.getContractFactory('CronReceiver', owner)
    cronReceiver1 = await crFactory.deploy()
    cronReceiver2 = await crFactory.deploy()
    const cronDelegateFactory = await ethers.getContractFactory(
      'CronUpkeepDelegate',
      admin,
    )
    cronDelegate = await cronDelegateFactory.deploy()
    const cronExternalFactory = await ethers.getContractFactory(
      'src/v0.8/libraries/external/Cron.sol:Cron',
      admin,
    )
    const cronExternalLib = await cronExternalFactory.deploy()
    cronFactory = await ethers.getContractFactory('CronUpkeepTestHelper', {
      signer: admin,
      libraries: { Cron: cronExternalLib.address },
    })
    cron = (
      await cronFactory.deploy(owner.address, cronDelegate.address, 5, [])
    ).connect(owner)
    const cronFactoryContractFactory = await ethers.getContractFactory(
      'CronUpkeepFactory',
      { signer: admin, libraries: { Cron: cronExternalLib.address } },
    ) // the typechain factory that creates the cron factory contract
    cronFactoryContract = await cronFactoryContractFactory.deploy()
    const fs = cronReceiver1.interface.functions
    handler1Sig = utils.id(fs['handler1()'].format('sighash')).slice(0, 10) // TODO this seems like an ethers bug
    handler2Sig = utils.id(fs['handler2()'].format('sighash')).slice(0, 10)
    basicSpec = await cronFactoryContract.encodeCronString(basicCronString)
  })

  afterEach(async () => {
    await h.reset()
  })

  describe('updateCronJob()', () => {
    const newCronString = '0 0 1 1 1'
    let newEncodedSpec: string
    beforeEach(async () => {
      await createBasicCron()
      newEncodedSpec = await cronFactoryContract.encodeCronString(newCronString)
    })

    it('updates a cron job', async () => {
      let cron1 = await cron.getCronJob(1)
      assert.equal(cron1.target, cronReceiver1.address)
      assert.equal(cron1.handler, handler1Sig)
      assert.equal(cron1.cronString, basicCronString)
      await cron.updateCronJob(
        1,
        cronReceiver2.address,
        handler2Sig,
        newEncodedSpec,
      )
      cron1 = await cron.getCronJob(1)
      assert.equal(cron1.target, cronReceiver2.address)
      assert.equal(cron1.handler, handler2Sig)
      assert.equal(cron1.cronString, newCronString)
    })

    it('emits an event', async () => {
      await expect(
        await cron.updateCronJob(
          1,
          cronReceiver2.address,
          handler2Sig,
          newEncodedSpec,
        ),
      ).to.emit(cron, 'CronJobUpdated')
    })

    it('is only callable by the owner', async () => {
      await expect(
        cron
          .connect(stranger)
          .updateCronJob(1, cronReceiver2.address, handler2Sig, newEncodedSpec),
      ).to.be.revertedWith(OWNABLE_ERR)
    })

    it('reverts if trying to update a non-existent ID', async () => {
      await expect(
        cron.updateCronJob(
          2,
          cronReceiver2.address,
          handler2Sig,
          newEncodedSpec,
        ),
      ).to.be.revertedWith(CRON_NOT_FOUND_ERR)
    })
  })

  describe('deleteCronJob()', () => {
    it("deletes a jobs by it's ID", async () => {
      await createBasicCron()
      await createBasicCron()
      await createBasicCron()
      await createBasicCron()
      await assertJobIDsEqual([1, 2, 3, 4])
      await cron.deleteCronJob(2)
      await expect(cron.getCronJob(2)).to.be.revertedWith(CRON_NOT_FOUND_ERR)
      await expect(cron.deleteCronJob(2)).to.be.revertedWith(CRON_NOT_FOUND_ERR)
      await assertJobIDsEqual([1, 3, 4])
      await cron.deleteCronJob(1)
      await assertJobIDsEqual([3, 4])
      await cron.deleteCronJob(4)
      await assertJobIDsEqual([3])
      await cron.deleteCronJob(3)
      await assertJobIDsEqual([])
    })

    it('emits an event', async () => {
      await createBasicCron()
      await expect(cron.deleteCronJob(1)).to.emit(cron, 'CronJobDeleted')
    })

    it('reverts if trying to delete a non-existent ID', async () => {
      await createBasicCron()
      await createBasicCron()
      await expect(cron.deleteCronJob(0)).to.be.revertedWith(CRON_NOT_FOUND_ERR)
      await expect(cron.deleteCronJob(3)).to.be.revertedWith(CRON_NOT_FOUND_ERR)
    })
  })

  describe('pause() / unpause()', () => {
    it('is only callable by the owner', async () => {
      await expect(cron.connect(stranger).pause()).to.be.reverted
      await expect(cron.connect(stranger).unpause()).to.be.reverted
    })

    it('pauses / unpauses the contract', async () => {
      expect(await cron.paused()).to.be.false
      await cron.pause()
      expect(await cron.paused()).to.be.true
      await cron.unpause()
      expect(await cron.paused()).to.be.false
    })
  })
})

// only run during yarn test:gas
describe.skip('Cron Gas Usage', () => {
  before(async () => {
    const accounts = await ethers.getSigners()
    admin = accounts[0]
    owner = accounts[1]
    const crFactory = await ethers.getContractFactory('CronReceiver', owner)
    cronReceiver1 = await crFactory.deploy()
    const cronDelegateFactory = await ethers.getContractFactory(
      'CronUpkeepDelegate',
      owner,
    )
    const cronDelegate = await cronDelegateFactory.deploy()
    const cronExternalFactory = await ethers.getContractFactory(
      'src/v0.8/libraries/external/Cron.sol:Cron',
      admin,
    )
    const cronExternalLib = await cronExternalFactory.deploy()
    const cronFactory = await ethers.getContractFactory(
      'CronUpkeepTestHelper',
      {
        signer: owner,
        libraries: { Cron: cronExternalLib.address },
      },
    )
    cron = await cronFactory.deploy(owner.address, cronDelegate.address, 5, [])
    const fs = cronReceiver1.interface.functions
    handler1Sig = utils
      .id(fs['handler1()'].format('sighash')) // TODO this seems like an ethers bug
      .slice(0, 10)
  })

  describe('checkUpkeep() / performUpkeep()', () => {
    it('uses gas', async () => {
      for (let idx = 0; idx < validCrons.length; idx++) {
        const cronString = validCrons[idx]
        const cronID = idx + 1
        await cron.createCronJobFromString(
          cronReceiver1.address,
          handler1Sig,
          cronString,
        )
        await h.fastForward(moment.duration(100, 'years').asSeconds()) // long enough that at least 1 tick occurs
        const [needsUpkeep, data] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isTrue(needsUpkeep, `failed for cron string ${cronString}`)
        await cron.txCheckUpkeep('0x')
        await cron.performUpkeep(data)
        await cron.deleteCronJob(cronID)
      }
    })
  })
})
