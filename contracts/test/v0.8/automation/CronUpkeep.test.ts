import moment from 'moment'
import { ethers } from 'hardhat'
import { Contract } from 'ethers'
import { assert, expect } from 'chai'
import { CronUpkeepTestHelper } from '../../../typechain/CronUpkeepTestHelper'
import { CronUpkeepDelegate } from '../../../typechain/CronUpkeepDelegate'
import { CronUpkeepFactory } from '../../../typechain/CronUpkeepFactory'
import { CronUpkeepTestHelper__factory as CronUpkeepTestHelperFactory } from '../../../typechain/factories/CronUpkeepTestHelper__factory'
import { CronInternalTestHelper } from '../../../typechain/CronInternalTestHelper'
import { CronReceiver } from '../../../typechain/CronReceiver'
import { BigNumber, BigNumberish } from '@ethersproject/bignumber'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { validCrons } from '../../test-helpers/fixtures'
import * as h from '../../test-helpers/helpers'

const { utils } = ethers
const { AddressZero } = ethers.constants

const OWNABLE_ERR = 'Only callable by owner'
const CRON_NOT_FOUND_ERR = 'CronJobIDNotFound'

let cron: CronUpkeepTestHelper
let cronFactory: CronUpkeepTestHelperFactory // the typechain factory that deploys cron upkeep contracts
let cronFactoryContract: CronUpkeepFactory // the cron factory contract
let cronDelegate: CronUpkeepDelegate
let cronTestHelper: CronInternalTestHelper
let cronReceiver1: CronReceiver
let cronReceiver2: CronReceiver

let admin: SignerWithAddress
let owner: SignerWithAddress
let stranger: SignerWithAddress

const timeStamp = 32503680000 // Jan 1, 3000 12:00AM
const basicCronString = '0 * * * *'

let handler1Sig: string
let handler2Sig: string
let revertHandlerSig: string
let basicSpec: string

async function assertJobIDsEqual(expected: number[]) {
  const ids = (await cron.getActiveCronJobIDs()).map((n) => n.toNumber())
  assert.deepEqual(ids.sort(), expected.sort())
}

function decodePayload(payload: string) {
  return utils.defaultAbiCoder.decode(
    ['uint256', 'uint256', 'address', 'bytes'],
    payload,
  ) as [BigNumber, BigNumber, string, string]
}

function encodePayload(payload: [BigNumberish, BigNumberish, string, string]) {
  return utils.defaultAbiCoder.encode(
    ['uint256', 'uint256', 'address', 'bytes'],
    payload,
  )
}

async function createBasicCron() {
  return await cron.createCronJobFromEncodedSpec(
    cronReceiver1.address,
    handler1Sig,
    basicSpec,
  )
}

describe('CronUpkeep', () => {
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
      'src/v0.8/automation/libraries/external/Cron.sol:Cron',
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
    handler1Sig = utils.id(fs['handler1()'].format('sighash')).slice(0, 10)
    handler2Sig = utils.id(fs['handler2()'].format('sighash')).slice(0, 10)
    revertHandlerSig = utils
      .id(fs['revertHandler()'].format('sighash'))
      .slice(0, 10)
    const cronTHFactory = await ethers.getContractFactory(
      'CronInternalTestHelper',
    )
    cronTestHelper = await cronTHFactory.deploy()
    basicSpec = await cronFactoryContract.encodeCronString(basicCronString)
  })

  afterEach(async () => {
    await h.reset()
  })

  it('has a limited public ABI [ @skip-coverage ]', () => {
    // Casting cron is necessary due to a tricky versioning mismatch issue, likely between ethers
    // and typechain. Remove once the version issue is resolved.
    // https://smartcontract-it.atlassian.net/browse/ARCHIVE-22094
    h.publicAbi(cron as unknown as Contract, [
      's_maxJobs',
      'performUpkeep',
      'createCronJobFromEncodedSpec',
      'updateCronJob',
      'deleteCronJob',
      'checkUpkeep',
      'getActiveCronJobIDs',
      'getCronJob',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // Pausable methods
      'paused',
      'pause',
      'unpause',
      // Cron helper methods
      'createCronJobFromString',
      'txCheckUpkeep',
    ])
  })

  describe('constructor()', () => {
    it('sets the initial values', async () => {
      expect(await cron.owner()).to.equal(owner.address)
      expect(await cron.s_maxJobs()).to.equal(5)
    })

    it('optionally creates a first job', async () => {
      const payload = await cronFactoryContract.encodeCronJob(
        cronReceiver1.address,
        handler1Sig,
        basicCronString,
      )
      cron = (
        await cronFactory.deploy(
          owner.address,
          cronDelegate.address,
          5,
          payload,
        )
      ).connect(owner)
      const job = await cron.getCronJob(1)
      assert.equal(job.target, cronReceiver1.address)
      assert.equal(job.handler, handler1Sig)
      assert.equal(job.cronString, basicCronString)
    })
  })

  describe('checkUpkeep() / performUpkeep()', () => {
    beforeEach(async () => {
      await h.setTimestamp(timeStamp)
      // id 1
      await cron.createCronJobFromString(
        cronReceiver1.address,
        handler1Sig,
        '0 0 31 * *', // 31st day of every month
      )
      // id 2
      await cron.createCronJobFromString(
        cronReceiver1.address,
        handler2Sig,
        '10 * * * *', // on the 10 min mark
      )
      // id 3
      await cron.createCronJobFromString(
        cronReceiver2.address,
        handler1Sig,
        '0 0 * 7 *', // every day in July
      )
      // id 4
      await cron.createCronJobFromString(
        cronReceiver2.address,
        revertHandlerSig,
        '20 * * * *', // on the 20 min mark
      )
    })

    describe('checkUpkeep()', () => {
      it('returns false if no one is elligible', async () => {
        const [needsUpkeep] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isFalse(needsUpkeep)
      })

      it('returns the id of eligible cron jobs', async () => {
        await h.fastForward(moment.duration(11, 'minutes').asSeconds())
        const [needsUpkeep, payload] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isTrue(needsUpkeep)
        const [id, ..._] = decodePayload(payload)
        assert.equal(id.toNumber(), 2)
      })

      describe('when mutiple crons are elligible', () => {
        it('cycles through the cron IDs based on block number', async () => {
          await h.fastForward(moment.duration(1, 'year').asSeconds())
          let [_, payload] = await cron
            .connect(AddressZero)
            .callStatic.checkUpkeep('0x')
          const [id1] = decodePayload(payload)
          await h.mineBlock(ethers.provider)
          ;[_, payload] = await cron
            .connect(AddressZero)
            .callStatic.checkUpkeep('0x')
          const [id2] = decodePayload(payload)
          await h.mineBlock(ethers.provider)
          ;[_, payload] = await cron
            .connect(AddressZero)
            .callStatic.checkUpkeep('0x')
          const [id3] = decodePayload(payload)
          await h.mineBlock(ethers.provider)
          ;[_, payload] = await cron
            .connect(AddressZero)
            .callStatic.checkUpkeep('0x')
          const [id4] = decodePayload(payload)
          assert.deepEqual(
            [id1, id2, id3, id4].map((n) => n.toNumber()).sort(),
            [1, 2, 3, 4],
          )
        })
      })
    })

    describe('performUpkeep()', () => {
      it('forwards the call to the appropriate target/handler', async () => {
        await h.fastForward(moment.duration(11, 'minutes').asSeconds())
        const [needsUpkeep, payload] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isTrue(needsUpkeep)
        await expect(cron.performUpkeep(payload)).to.emit(
          cronReceiver1,
          'Received2',
        )
      })

      it('emits an event', async () => {
        await h.fastForward(moment.duration(11, 'minutes').asSeconds())
        const [needsUpkeep, payload] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isTrue(needsUpkeep)
        await expect(cron.performUpkeep(payload))
          .to.emit(cron, 'CronJobExecuted')
          .withArgs(2, true)
      })

      it('succeeds even if the call to the target fails', async () => {
        await cron.deleteCronJob(2)
        await h.fastForward(moment.duration(21, 'minutes').asSeconds())
        const payload = encodePayload([
          4,
          moment.unix(timeStamp).add(20, 'minutes').unix(),
          cronReceiver2.address,
          revertHandlerSig,
        ])
        await expect(cron.performUpkeep(payload))
          .to.emit(cron, 'CronJobExecuted')
          .withArgs(4, false)
      })

      it('is only callable by anyone', async () => {
        await h.fastForward(moment.duration(11, 'minutes').asSeconds())
        const [needsUpkeep, payload] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isTrue(needsUpkeep)
        await cron.connect(stranger).performUpkeep(payload)
      })

      it('is only callable once for a given tick', async () => {
        await h.fastForward(moment.duration(10, 'minutes').asSeconds())
        const [needsUpkeep, payload] = await cron
          .connect(AddressZero)
          .callStatic.checkUpkeep('0x')
        assert.isTrue(needsUpkeep)
        const maliciousPayload = encodePayload([
          2,
          moment.unix(timeStamp).add(10, 'minutes').add(59, 'seconds').unix(),
          cronReceiver1.address,
          handler2Sig,
        ])
        await cron.performUpkeep(payload)
        await expect(cron.performUpkeep(payload)).to.be.reverted
        await expect(cron.performUpkeep(maliciousPayload)).to.be.reverted
        await h.fastForward(moment.duration(1, 'minute').asSeconds())
        await expect(cron.performUpkeep(payload)).to.be.reverted
        await expect(cron.performUpkeep(maliciousPayload)).to.be.reverted
        await h.fastForward(moment.duration(10, 'minute').asSeconds())
        await expect(cron.performUpkeep(payload)).to.be.reverted
        await expect(cron.performUpkeep(maliciousPayload)).to.be.reverted
      })
    })
  })

  describe('createCronJobFromEncodedSpec()', () => {
    it('creates jobs with sequential IDs', async () => {
      const cronString1 = '0 * * * *'
      const cronString2 = '0 1,2,3 */4 5-6 1-2'
      const encodedSpec1 =
        await cronFactoryContract.encodeCronString(cronString1)
      const encodedSpec2 =
        await cronFactoryContract.encodeCronString(cronString2)
      const nextTick1 = (
        await cronTestHelper.calculateNextTick(cronString1)
      ).toNumber()
      const nextTick2 = (
        await cronTestHelper.calculateNextTick(cronString2)
      ).toNumber()
      await cron.createCronJobFromEncodedSpec(
        cronReceiver1.address,
        handler1Sig,
        encodedSpec1,
      )
      await assertJobIDsEqual([1])
      await cron.createCronJobFromEncodedSpec(
        cronReceiver1.address,
        handler2Sig,
        encodedSpec1,
      )
      await assertJobIDsEqual([1, 2])
      await cron.createCronJobFromEncodedSpec(
        cronReceiver2.address,
        handler1Sig,
        encodedSpec2,
      )
      await assertJobIDsEqual([1, 2, 3])
      await cron.createCronJobFromEncodedSpec(
        cronReceiver2.address,
        handler2Sig,
        encodedSpec2,
      )
      await assertJobIDsEqual([1, 2, 3, 4])
      const cron1 = await cron.getCronJob(1)
      const cron2 = await cron.getCronJob(2)
      const cron3 = await cron.getCronJob(3)
      const cron4 = await cron.getCronJob(4)
      assert.equal(cron1.target, cronReceiver1.address)
      assert.equal(cron1.handler, handler1Sig)
      assert.equal(cron1.cronString, cronString1)
      assert.equal(cron1.nextTick.toNumber(), nextTick1)
      assert.equal(cron2.target, cronReceiver1.address)
      assert.equal(cron2.handler, handler2Sig)
      assert.equal(cron2.cronString, cronString1)
      assert.equal(cron2.nextTick.toNumber(), nextTick1)
      assert.equal(cron3.target, cronReceiver2.address)
      assert.equal(cron3.handler, handler1Sig)
      assert.equal(cron3.cronString, cronString2)
      assert.equal(cron3.nextTick.toNumber(), nextTick2)
      assert.equal(cron4.target, cronReceiver2.address)
      assert.equal(cron4.handler, handler2Sig)
      assert.equal(cron4.cronString, cronString2)
      assert.equal(cron4.nextTick.toNumber(), nextTick2)
    })

    it('emits an event', async () => {
      await expect(createBasicCron()).to.emit(cron, 'CronJobCreated')
    })

    it('is only callable by the owner', async () => {
      await expect(
        cron
          .connect(stranger)
          .createCronJobFromEncodedSpec(
            cronReceiver1.address,
            handler1Sig,
            basicSpec,
          ),
      ).to.be.revertedWith(OWNABLE_ERR)
    })

    it('errors if trying to create more jobs than allowed', async () => {
      for (let idx = 0; idx < 5; idx++) {
        await createBasicCron()
      }
      await expect(createBasicCron()).to.be.revertedWithCustomError(
        cron,
        'ExceedsMaxJobs',
      )
    })
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
      ).to.be.revertedWithCustomError(cron, CRON_NOT_FOUND_ERR)
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
      await expect(cron.getCronJob(2)).to.be.revertedWithCustomError(
        cron,
        CRON_NOT_FOUND_ERR,
      )
      await expect(cron.deleteCronJob(2)).to.be.revertedWithCustomError(
        cron,
        CRON_NOT_FOUND_ERR,
      )
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
      await expect(cron.deleteCronJob(0)).to.be.revertedWithCustomError(
        cron,
        CRON_NOT_FOUND_ERR,
      )
      await expect(cron.deleteCronJob(3)).to.be.revertedWithCustomError(
        cron,
        CRON_NOT_FOUND_ERR,
      )
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

// only run during pnpm test:gas
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
      'src/v0.8/automation/libraries/external/Cron.sol:Cron',
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
