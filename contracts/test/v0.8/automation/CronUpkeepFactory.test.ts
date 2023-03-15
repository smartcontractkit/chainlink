import { ethers } from 'hardhat'
import { Contract } from 'ethers'
import { assert, expect } from 'chai'
import { CronUpkeepFactory } from '../../../typechain/CronUpkeepFactory'
import type { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import { reset } from '../../test-helpers/helpers'
import * as h from '../../test-helpers/helpers'

const OWNABLE_ERR = 'Only callable by owner'

let cronExternalLib: Contract
let factory: CronUpkeepFactory

let admin: SignerWithAddress
let owner: SignerWithAddress
let stranger: SignerWithAddress

describe('CronUpkeepFactory', () => {
  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    admin = accounts[0]
    owner = accounts[1]
    stranger = accounts[2]
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

  it('has a limited public ABI [ @skip-coverage ]', () => {
    h.publicAbi(factory as unknown as Contract, [
      's_maxJobs',
      'newCronUpkeep',
      'newCronUpkeepWithJob',
      'setMaxJobs',
      'cronDelegateAddress',
      'encodeCronString',
      'encodeCronJob',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
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

  describe('setMaxJobs()', () => {
    it('sets the max jobs value', async () => {
      expect(await factory.s_maxJobs()).to.equal(5)
      await factory.setMaxJobs(6)
      expect(await factory.s_maxJobs()).to.equal(6)
    })

    it('is only callable by the owner', async () => {
      await expect(factory.connect(stranger).setMaxJobs(6)).to.be.revertedWith(
        OWNABLE_ERR,
      )
    })
  })
})
