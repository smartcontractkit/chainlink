import { ethers } from 'hardhat'
import { Signer, Contract, BigNumber } from 'ethers'
import { assert, expect } from 'chai'
import { publicAbi } from '../../test-helpers/helpers'
import { randomAddressString } from 'hardhat/internal/hardhat-network/provider/fork/random'

describe('VRFCoordinatorV2', () => {
  let vrfCoordinatorV2: Contract
  let vrfCoordinatorV2TestHelper: Contract
  let linkToken: Contract
  let blockHashStore: Contract
  let mockLinkEth: Contract
  let owner: Signer
  let subOwner: Signer
  let subOwnerAddress: string
  let consumer: Signer
  let random: Signer
  let randomAddress: string
  let oracle: Signer
  const linkEth = BigNumber.from(300000000)
  type config = {
    minimumRequestBlockConfirmations: number
    maxGasLimit: number
    stalenessSeconds: number
    gasAfterPaymentCalculation: number
    weiPerUnitLink: BigNumber
  }
  let c: config

  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    owner = accounts[0]
    subOwner = accounts[1]
    subOwnerAddress = await subOwner.getAddress()
    consumer = accounts[2]
    random = accounts[3]
    randomAddress = await random.getAddress()
    oracle = accounts[4]
    const ltFactory = await ethers.getContractFactory(
      'src/v0.4/LinkToken.sol:LinkToken',
      accounts[0],
    )
    linkToken = await ltFactory.deploy()
    const bhFactory = await ethers.getContractFactory(
      'src/v0.6/dev/BlockhashStore.sol:BlockhashStore',
      accounts[0],
    )
    blockHashStore = await bhFactory.deploy()
    const mockAggregatorV3Factory = await ethers.getContractFactory(
      'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
      accounts[0],
    )
    mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth)
    const vrfCoordinatorV2Factory = await ethers.getContractFactory(
      'src/v0.8/VRFCoordinatorV2.sol:VRFCoordinatorV2',
      accounts[0],
    )
    vrfCoordinatorV2 = await vrfCoordinatorV2Factory.deploy(
      linkToken.address,
      blockHashStore.address,
      mockLinkEth.address,
    )
    const vrfCoordinatorV2TestHelperFactory = await ethers.getContractFactory(
      'src/v0.8/tests/VRFCoordinatorV2TestHelper.sol:VRFCoordinatorV2TestHelper',
      accounts[0],
    )
    vrfCoordinatorV2TestHelper = await vrfCoordinatorV2TestHelperFactory.deploy(
      linkToken.address,
      blockHashStore.address,
      mockLinkEth.address,
    )
    await linkToken.transfer(
      subOwnerAddress,
      BigNumber.from('1000000000000000000'),
    ) // 1 link
    await linkToken.transfer(
      randomAddress,
      BigNumber.from('1000000000000000000'),
    ) // 1 link
    c = {
      minimumRequestBlockConfirmations: 1,
      maxGasLimit: 1000000,
      stalenessSeconds: 86400,
      gasAfterPaymentCalculation:
        21000 + 5000 + 2100 + 20000 + 2 * 2100 - 15000 + 7315,
      weiPerUnitLink: BigNumber.from('10000000000000000'),
    }
    // Note if you try and use an object, ethers
    // confuses that with an override object and will error.
    // It appears that only arrays work for struct args.
    const fc = [0, 0, 0, 0, 0, 0, 0, 0, 0]
    await vrfCoordinatorV2
      .connect(owner)
      .setConfig(
        c.minimumRequestBlockConfirmations,
        c.maxGasLimit,
        c.stalenessSeconds,
        c.gasAfterPaymentCalculation,
        c.weiPerUnitLink,
        fc,
      )
  })

  it('has a limited public interface [ @skip-coverage ]', async () => {
    publicAbi(vrfCoordinatorV2, [
      // Public constants
      'MAX_CONSUMERS',
      'MAX_NUM_WORDS',
      'MAX_REQUEST_CONFIRMATIONS',
      // Owner
      'acceptOwnership',
      'transferOwnership',
      'owner',
      'getConfig',
      'getFeeConfig',
      'getFallbackWeiPerUnitLink',
      'getCurrentSubId',
      'setConfig',
      'getRequestConfig',
      'recoverFunds',
      'ownerCancelSubscription',
      'getFeeTier',
      'pendingRequestExists',
      'getTotalBalance',
      // Oracle
      'requestRandomWords',
      'getCommitment', // Note we use this to check if a request is already fulfilled.
      'hashOfKey',
      'fulfillRandomWords',
      'registerProvingKey',
      'deregisterProvingKey',
      'oracleWithdraw',
      // Subscription management
      'createSubscription',
      'addConsumer',
      'removeConsumer',
      'getSubscription',
      'onTokenTransfer', // Effectively the fundSubscription.
      'cancelSubscription',
      'requestSubscriptionOwnerTransfer',
      'acceptSubscriptionOwnerTransfer',
      // Misc
      'typeAndVersion',
      'BLOCKHASH_STORE',
      'LINK',
      'LINK_ETH_FEED',
    ])
  })

  describe('#setConfig', async function () {
    it('only owner can set', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .setConfig(
            c.minimumRequestBlockConfirmations,
            c.maxGasLimit,
            c.stalenessSeconds,
            c.gasAfterPaymentCalculation,
            c.weiPerUnitLink,
            [0, 0, 0, 0, 0, 0, 0, 0, 0],
          ),
      ).to.be.revertedWith('Only callable by owner')
      // Anyone can read the config.
      const resp = await vrfCoordinatorV2.connect(random).getConfig()
      assert(resp[0] == c.minimumRequestBlockConfirmations)
      assert(resp[1] == c.maxGasLimit)
      assert(resp[2] == c.stalenessSeconds)
      assert(resp[3].toString() == c.gasAfterPaymentCalculation.toString())
    })

    it('max req confs', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(owner)
          .setConfig(
            201,
            c.maxGasLimit,
            c.stalenessSeconds,
            c.gasAfterPaymentCalculation,
            c.weiPerUnitLink,
            [0, 0, 0, 0, 0, 0, 0, 0, 0],
          ),
      ).to.be.revertedWith('InvalidRequestConfirmations(201, 201, 200)')
    })

    it('positive fallback price', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(owner)
          .setConfig(
            c.minimumRequestBlockConfirmations,
            c.maxGasLimit,
            c.stalenessSeconds,
            c.gasAfterPaymentCalculation,
            0,
            [0, 0, 0, 0, 0, 0, 0, 0, 0],
          ),
      ).to.be.revertedWith('InvalidLinkWeiPrice(0)')
      await expect(
        vrfCoordinatorV2
          .connect(owner)
          .setConfig(
            c.minimumRequestBlockConfirmations,
            c.maxGasLimit,
            c.stalenessSeconds,
            c.gasAfterPaymentCalculation,
            -1,
            [0, 0, 0, 0, 0, 0, 0, 0, 0],
          ),
      ).to.be.revertedWith('InvalidLinkWeiPrice(-1)')
    })
  })

  async function createSubscription(): Promise<number> {
    // let consumers: string[] = [await consumer.getAddress()];
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription()
    const receipt = await tx.wait()
    const subId = receipt.events[0].args['subId']
    await vrfCoordinatorV2
      .connect(subOwner)
      .addConsumer(subId, await consumer.getAddress())
    return subId
  }

  async function createSubscriptionWithConsumers(
    consumers: string[],
  ): Promise<number> {
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription()
    const receipt = await tx.wait()
    const subId = receipt.events[0].args['subId']
    for (let i = 0; i < consumers.length; i++) {
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, consumers[i])
    }
    return subId
  }

  describe('#createSubscription', async function () {
    it('can create a subscription', async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).createSubscription())
        .to.emit(vrfCoordinatorV2, 'SubscriptionCreated')
        .withArgs(1, subOwnerAddress)
      const s = await vrfCoordinatorV2.getSubscription(1)
      assert(s.balance.toString() == '0', 'invalid balance')
      assert(s.owner == subOwnerAddress, 'invalid address')
    })
    it('subscription id increments', async function () {
      await expect(vrfCoordinatorV2.connect(subOwner).createSubscription())
        .to.emit(vrfCoordinatorV2, 'SubscriptionCreated')
        .withArgs(1, subOwnerAddress)
      await expect(vrfCoordinatorV2.connect(subOwner).createSubscription())
        .to.emit(vrfCoordinatorV2, 'SubscriptionCreated')
        .withArgs(2, subOwnerAddress)
    })
    it('cannot create more than the max', async function () {
      const subId = createSubscriptionWithConsumers([])
      for (let i = 0; i < 100; i++) {
        await vrfCoordinatorV2
          .connect(subOwner)
          .addConsumer(subId, randomAddressString())
      }
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .addConsumer(subId, randomAddressString()),
      ).to.be.revertedWith(`TooManyConsumers()`)
    })
  })

  describe('#requestSubscriptionOwnerTransfer', async function () {
    let subId: number
    beforeEach(async () => {
      subId = await createSubscription()
    })
    it('rejects non-owner', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(random)
          .requestSubscriptionOwnerTransfer(subId, randomAddress),
      ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
    })
    it('owner can request transfer', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .requestSubscriptionOwnerTransfer(subId, randomAddress),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionOwnerTransferRequested')
        .withArgs(subId, subOwnerAddress, randomAddress)
      // Same request is a noop
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .requestSubscriptionOwnerTransfer(subId, randomAddress),
      ).to.not.emit(vrfCoordinatorV2, 'SubscriptionOwnerTransferRequested')
    })
  })

  describe('#acceptSubscriptionOwnerTransfer', async function () {
    let subId: number
    beforeEach(async () => {
      subId = await createSubscription()
    })
    it('subscription must exist', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .acceptSubscriptionOwnerTransfer(1203123123),
      ).to.be.revertedWith(`InvalidSubscription`)
    })
    it('must be requested owner to accept', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .requestSubscriptionOwnerTransfer(subId, randomAddress),
      )
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .acceptSubscriptionOwnerTransfer(subId),
      ).to.be.revertedWith(`MustBeRequestedOwner("${randomAddress}")`)
    })
    it('requested owner can accept', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .requestSubscriptionOwnerTransfer(subId, randomAddress),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionOwnerTransferRequested')
        .withArgs(subId, subOwnerAddress, randomAddress)
      await expect(
        vrfCoordinatorV2.connect(random).acceptSubscriptionOwnerTransfer(subId),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionOwnerTransferred')
        .withArgs(subId, subOwnerAddress, randomAddress)
    })
  })

  describe('#addConsumer', async function () {
    let subId: number
    beforeEach(async () => {
      subId = await createSubscription()
    })
    it('subscription must exist', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .addConsumer(1203123123, randomAddress),
      ).to.be.revertedWith(`InvalidSubscription`)
    })
    it('must be owner', async function () {
      await expect(
        vrfCoordinatorV2.connect(random).addConsumer(subId, randomAddress),
      ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
    })
    it('add is idempotent', async function () {
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
    })
    it('cannot add more than maximum', async function () {
      // There is one consumer, add another 99 to hit the max
      for (let i = 0; i < 99; i++) {
        await vrfCoordinatorV2
          .connect(subOwner)
          .addConsumer(subId, randomAddressString())
      }
      // Adding one more should fail
      // await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress);
      await expect(
        vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress),
      ).to.be.revertedWith(`TooManyConsumers()`)
      // Same is true if we first create with the maximum
      const consumers: string[] = []
      for (let i = 0; i < 100; i++) {
        consumers.push(randomAddressString())
      }
      subId = await createSubscriptionWithConsumers(consumers)
      await expect(
        vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress),
      ).to.be.revertedWith(`TooManyConsumers()`)
    })
    it('owner can update', async function () {
      await expect(
        vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionConsumerAdded')
        .withArgs(subId, randomAddress)
    })
  })

  describe('#removeConsumer', async function () {
    let subId: number
    beforeEach(async () => {
      subId = await createSubscription()
    })
    it('subscription must exist', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .removeConsumer(1203123123, randomAddress),
      ).to.be.revertedWith(`InvalidSubscription`)
    })
    it('must be owner', async function () {
      await expect(
        vrfCoordinatorV2.connect(random).removeConsumer(subId, randomAddress),
      ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
    })
    it('owner can update', async function () {
      const subBefore = await vrfCoordinatorV2.getSubscription(subId)
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
      await expect(
        vrfCoordinatorV2.connect(subOwner).removeConsumer(subId, randomAddress),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionConsumerRemoved')
        .withArgs(subId, randomAddress)
      const subAfter = await vrfCoordinatorV2.getSubscription(subId)
      // Subscription should NOT contain the removed consumer
      assert.deepEqual(subBefore.consumers, subAfter.consumers)
    })
    it('can remove all consumers', async function () {
      // Testing the handling of zero.
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
      await vrfCoordinatorV2
        .connect(subOwner)
        .removeConsumer(subId, randomAddress)
      await vrfCoordinatorV2
        .connect(subOwner)
        .removeConsumer(subId, await consumer.getAddress())
      // Should be empty
      const subAfter = await vrfCoordinatorV2.getSubscription(subId)
      assert.deepEqual(subAfter.consumers, [])
    })
  })

  describe('#cancelSubscription', async function () {
    let subId: number
    beforeEach(async () => {
      subId = await createSubscription()
    })
    it('subscription must exist', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .cancelSubscription(1203123123, subOwnerAddress),
      ).to.be.revertedWith(`InvalidSubscription`)
    })
    it('must be owner', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(random)
          .cancelSubscription(subId, subOwnerAddress),
      ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
    })
    it('can cancel', async function () {
      await linkToken
        .connect(subOwner)
        .transferAndCall(
          vrfCoordinatorV2.address,
          BigNumber.from('1000'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        )
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .cancelSubscription(subId, randomAddress),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionCanceled')
        .withArgs(subId, randomAddress, BigNumber.from('1000'))
      const randomBalance = await linkToken.balanceOf(randomAddress)
      assert.equal(randomBalance.toString(), '1000000000000001000')
      await expect(
        vrfCoordinatorV2.connect(subOwner).getSubscription(subId),
      ).to.be.revertedWith('InvalidSubscription')
    })
    it('can add same consumer after canceling', async function () {
      await linkToken
        .connect(subOwner)
        .transferAndCall(
          vrfCoordinatorV2.address,
          BigNumber.from('1000'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        )
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
      await vrfCoordinatorV2
        .connect(subOwner)
        .cancelSubscription(subId, randomAddress)
      subId = await createSubscription()
      // The cancel should have removed this consumer, so we can add it again.
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
    })
    it('cannot cancel with pending req', async function () {
      await linkToken
        .connect(subOwner)
        .transferAndCall(
          vrfCoordinatorV2.address,
          BigNumber.from('1000'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        )
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      await vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey)
      await vrfCoordinatorV2.connect(owner).reg
      const kh = await vrfCoordinatorV2.hashOfKey(testKey)
      await vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        1, // minReqConf
        1000000, // callbackGasLimit
        1, // numWords
      )
      // Should revert with outstanding requests
      await expect(
        vrfCoordinatorV2
          .connect(subOwner)
          .cancelSubscription(subId, randomAddress),
      ).to.be.revertedWith('PendingRequestExists()')
      // However the owner is able to cancel
      // funds go to the sub owner.
      await expect(
        vrfCoordinatorV2.connect(owner).ownerCancelSubscription(subId),
      )
        .to.emit(vrfCoordinatorV2, 'SubscriptionCanceled')
        .withArgs(subId, subOwnerAddress, BigNumber.from('1000'))
    })
  })

  describe('#recoverFunds', async function () {
    let subId: number
    beforeEach(async () => {
      subId = await createSubscription()
    })

    // Note we can't test the oracleWithdraw without fulfilling a request, so leave
    // that coverage to the go tests.
    it('function that should change internal balance do', async function () {
      type bf = [() => Promise<any>, BigNumber]
      const balanceChangingFns: Array<bf> = [
        [
          async function () {
            const s = ethers.utils.defaultAbiCoder.encode(['uint64'], [subId])
            await linkToken
              .connect(subOwner)
              .transferAndCall(
                vrfCoordinatorV2.address,
                BigNumber.from('1000'),
                s,
              )
          },
          BigNumber.from('1000'),
        ],
        [
          async function () {
            await vrfCoordinatorV2
              .connect(subOwner)
              .cancelSubscription(subId, randomAddress)
          },
          BigNumber.from('-1000'),
        ],
      ]
      for (const [fn, expectedBalanceChange] of balanceChangingFns) {
        const startingBalance = await vrfCoordinatorV2.getTotalBalance()
        await fn()
        const endingBalance = await vrfCoordinatorV2.getTotalBalance()
        assert(
          endingBalance.sub(startingBalance).toString() ==
            expectedBalanceChange.toString(),
        )
      }
    })
    it('only owner can recover', async function () {
      await expect(
        vrfCoordinatorV2.connect(subOwner).recoverFunds(randomAddress),
      ).to.be.revertedWith(`Only callable by owner`)
    })

    it('owner can recover link transferred', async function () {
      // Set the internal balance
      assert(BigNumber.from('0'), linkToken.balanceOf(randomAddress))
      const s = ethers.utils.defaultAbiCoder.encode(['uint64'], [subId])
      await linkToken
        .connect(subOwner)
        .transferAndCall(vrfCoordinatorV2.address, BigNumber.from('1000'), s)
      // Circumvent internal balance
      await linkToken
        .connect(subOwner)
        .transfer(vrfCoordinatorV2.address, BigNumber.from('1000'))
      // Should recover this 1000
      await expect(vrfCoordinatorV2.connect(owner).recoverFunds(randomAddress))
        .to.emit(vrfCoordinatorV2, 'FundsRecovered')
        .withArgs(randomAddress, BigNumber.from('1000'))
      assert(BigNumber.from('1000'), linkToken.balanceOf(randomAddress))
    })
  })

  it('subscription lifecycle', async function () {
    // Create subscription.
    const tx = await vrfCoordinatorV2.connect(subOwner).createSubscription()
    const receipt = await tx.wait()
    assert(receipt.events[0].event == 'SubscriptionCreated')
    assert(receipt.events[0].args['owner'] == subOwnerAddress, 'sub owner')
    const subId = receipt.events[0].args['subId']
    await vrfCoordinatorV2
      .connect(subOwner)
      .addConsumer(subId, await consumer.getAddress())

    // Fund the subscription
    await expect(
      linkToken
        .connect(subOwner)
        .transferAndCall(
          vrfCoordinatorV2.address,
          BigNumber.from('1000000000000000000'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        ),
    )
      .to.emit(vrfCoordinatorV2, 'SubscriptionFunded')
      .withArgs(subId, BigNumber.from(0), BigNumber.from('1000000000000000000'))

    // Non-owners cannot change the consumers
    await expect(
      vrfCoordinatorV2.connect(random).addConsumer(subId, randomAddress),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
    await expect(
      vrfCoordinatorV2.connect(random).removeConsumer(subId, randomAddress),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)

    // Non-owners cannot ask to transfer ownership
    await expect(
      vrfCoordinatorV2
        .connect(random)
        .requestSubscriptionOwnerTransfer(subId, randomAddress),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)

    // Owners can request ownership transfership
    await expect(
      vrfCoordinatorV2
        .connect(subOwner)
        .requestSubscriptionOwnerTransfer(subId, randomAddress),
    )
      .to.emit(vrfCoordinatorV2, 'SubscriptionOwnerTransferRequested')
      .withArgs(subId, subOwnerAddress, randomAddress)

    // Non-requested owners cannot accept
    await expect(
      vrfCoordinatorV2.connect(subOwner).acceptSubscriptionOwnerTransfer(subId),
    ).to.be.revertedWith(`MustBeRequestedOwner("${randomAddress}")`)

    // Requested owners can accept
    await expect(
      vrfCoordinatorV2.connect(random).acceptSubscriptionOwnerTransfer(subId),
    )
      .to.emit(vrfCoordinatorV2, 'SubscriptionOwnerTransferred')
      .withArgs(subId, subOwnerAddress, randomAddress)

    // Transfer it back to subOwner
    vrfCoordinatorV2
      .connect(random)
      .requestSubscriptionOwnerTransfer(subId, subOwnerAddress)
    vrfCoordinatorV2.connect(subOwner).acceptSubscriptionOwnerTransfer(subId)

    // Non-owners cannot cancel
    await expect(
      vrfCoordinatorV2.connect(random).cancelSubscription(subId, randomAddress),
    ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)

    await expect(
      vrfCoordinatorV2
        .connect(subOwner)
        .cancelSubscription(subId, randomAddress),
    )
      .to.emit(vrfCoordinatorV2, 'SubscriptionCanceled')
      .withArgs(subId, randomAddress, BigNumber.from('1000000000000000000'))
    const random2Balance = await linkToken.balanceOf(randomAddress)
    assert.equal(random2Balance.toString(), '2000000000000000000')
  })

  describe('#requestRandomWords', async function () {
    let subId: number
    let kh: string
    beforeEach(async () => {
      subId = await createSubscription()
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      kh = await vrfCoordinatorV2.hashOfKey(testKey)
    })
    it('invalid subId', async function () {
      await expect(
        vrfCoordinatorV2.connect(random).requestRandomWords(
          kh, // keyhash
          12301928312, // subId
          1, // minReqConf
          1000, // callbackGasLimit
          1, // numWords
        ),
      ).to.be.revertedWith(`InvalidSubscription()`)
    })
    it('invalid consumer', async function () {
      await expect(
        vrfCoordinatorV2.connect(random).requestRandomWords(
          kh, // keyhash
          subId, // subId
          1, // minReqConf
          1000, // callbackGasLimit
          1, // numWords
        ),
      ).to.be.revertedWith(
        `InvalidConsumer(${subId}, "${randomAddress.toString()}")`,
      )
    })
    it('invalid req confs', async function () {
      await expect(
        vrfCoordinatorV2.connect(consumer).requestRandomWords(
          kh, // keyhash
          subId, // subId
          0, // minReqConf
          1000, // callbackGasLimit
          1, // numWords
        ),
      ).to.be.revertedWith(`InvalidRequestConfirmations(0, 1, 200)`)
    })
    it('gas limit too high', async function () {
      await linkToken.connect(subOwner).transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from('1000000000000000000'), // 1 link > 0.1 min.
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
      await expect(
        vrfCoordinatorV2.connect(consumer).requestRandomWords(
          kh, // keyhash
          subId, // subId
          1, // minReqConf
          1000001, // callbackGasLimit
          1, // numWords
        ),
      ).to.be.revertedWith(`GasLimitTooBig(1000001, 1000000)`)
    })

    it('nonce increments', async function () {
      await linkToken.connect(subOwner).transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from('1000000000000000000'), // 1 link > 0.1 min.
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
      const r1 = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        1, // minReqConf
        1000000, // callbackGasLimit
        1, // numWords
      )
      const r1Receipt = await r1.wait()
      const seed1 = r1Receipt.events[0].args['requestId']
      const r2 = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        1, // minReqConf
        1000000, // callbackGasLimit
        1, // numWords
      )
      const r2Receipt = await r2.wait()
      const seed2 = r2Receipt.events[0].args['requestId']
      assert(seed2 != seed1)
    })

    it('emits correct log', async function () {
      await linkToken.connect(subOwner).transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from('1000000000000000000'), // 1 link > 0.1 min.
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
      const reqTx = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        1, // minReqConf
        1000, // callbackGasLimit
        1, // numWords
      )
      const reqReceipt = await reqTx.wait()
      assert(reqReceipt.events.length == 1)
      const reqEvent = reqReceipt.events[0]
      assert(reqEvent.event == 'RandomWordsRequested', 'wrong event name')
      assert(
        reqEvent.args['keyHash'] == kh,
        `wrong kh ${reqEvent.args['keyHash']} ${kh}`,
      )
      assert(
        reqEvent.args['subId'].toString() == subId.toString(),
        'wrong subId',
      )
      assert(
        reqEvent.args['minimumRequestConfirmations'].toString() ==
          BigNumber.from(1).toString(),
        'wrong minRequestConf',
      )
      assert(
        reqEvent.args['callbackGasLimit'] == 1000,
        'wrong callbackGasLimit',
      )
      assert(reqEvent.args['numWords'] == 1, 'wrong numWords')
      assert(
        reqEvent.args['sender'] == (await consumer.getAddress()),
        'wrong sender address',
      )
    })
    it('add/remove consumer invariant', async function () {
      await linkToken.connect(subOwner).transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from('1000000000000000000'), // 1 link > 0.1 min.
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
      await vrfCoordinatorV2.connect(subOwner).addConsumer(subId, randomAddress)
      await vrfCoordinatorV2
        .connect(subOwner)
        .removeConsumer(subId, randomAddress)
      await expect(
        vrfCoordinatorV2.connect(random).requestRandomWords(
          kh, // keyhash
          subId, // subId
          1, // minReqConf
          1000, // callbackGasLimit
          1, // numWords
        ),
      ).to.be.revertedWith(
        `InvalidConsumer(${subId}, "${randomAddress.toString()}")`,
      )
    })
    it('cancel/add subscription invariant', async function () {
      await linkToken.connect(subOwner).transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from('1000000000000000000'), // 1 link > 0.1 min.
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
      await vrfCoordinatorV2
        .connect(subOwner)
        .cancelSubscription(subId, randomAddress)
      subId = await createSubscriptionWithConsumers([])
      // Should not succeed because consumer was previously registered
      // i.e. cancel should be cleaning up correctly.
      await expect(
        vrfCoordinatorV2.connect(random).requestRandomWords(
          kh, // keyhash
          subId, // subId
          1, // minReqConf
          1000, // callbackGasLimit
          1, // numWords
        ),
      ).to.be.revertedWith(
        `InvalidConsumer(${subId}, "${randomAddress.toString()}")`,
      )
    })
  })

  describe('#oracleWithdraw', async function () {
    it('cannot withdraw with no balance', async function () {
      await expect(
        vrfCoordinatorV2
          .connect(oracle)
          .oracleWithdraw(randomAddressString(), BigNumber.from('100')),
      ).to.be.revertedWith(`InsufficientBalance`)
    })
  })

  describe('#calculatePaymentAmount [ @skip-coverage ]', async function () {
    it('output within sensible range', async function () {
      // By default, hardhat sends txes with the block limit as their gas limit.
      await vrfCoordinatorV2TestHelper
        .connect(oracle)
        .calculatePaymentAmountTest(
          BigNumber.from('0'), // Gas after payment
          0, // Fee PPM
          BigNumber.from('1000000000'), // Wei per unit gas (gas price)
        )
      const paymentAmount = await vrfCoordinatorV2TestHelper.getPaymentAmount()
      // The gas price is 1gwei and the eth/link price is set to 300000000 wei per unit link.
      // paymentAmount = 1e18*weiPerUnitGas*(gasAfterPaymentCalculation + startGas - gasleft()) / uint256(weiPerUnitLink);
      // So we expect x to be in the range (few thousand gas for the call)
      // 1e18*1e9*(1000 gas)/30000000 < x < 1e18*1e9*(5000 gas)/30000000
      // 3.333333333E22 < x < 1.666666667E23
      //const gss = await vrfCoordinatorV2TestHelper.getGasStart();
      assert(
        paymentAmount.gt(BigNumber.from('33333333330000000000000')),
        'payment too small',
      )
      assert(
        paymentAmount.lt(BigNumber.from('166666666600000000000000')),
        'payment too large',
      )
    })
    it('payment too large', async function () {
      // Set this gas price to be astronomical 1ETH/gas
      // That means the payment will be (even for 1gas)
      // 1e18*1e18/30000000
      // 3.333333333E28 > 1e27 (all link in existence)
      await expect(
        vrfCoordinatorV2TestHelper.connect(oracle).calculatePaymentAmountTest(
          BigNumber.from('0'), // Gas after payment
          0, // Fee PPM
          BigNumber.from('1000000000000000000'),
        ),
      ).to.be.revertedWith(`PaymentTooLarge()`)
    })

    it('non-positive link wei price should revert', async function () {
      const mockAggregatorV3Factory = await ethers.getContractFactory(
        'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
        owner,
      )
      const vrfCoordinatorV2TestHelperFactory = await ethers.getContractFactory(
        'VRFCoordinatorV2TestHelper',
        owner,
      )
      const mockLinkEthZero = await mockAggregatorV3Factory.deploy(0, 0)
      const vrfCoordinatorV2TestHelperZero =
        await vrfCoordinatorV2TestHelperFactory.deploy(
          linkToken.address,
          blockHashStore.address,
          mockLinkEthZero.address,
        )
      await expect(
        vrfCoordinatorV2TestHelperZero
          .connect(oracle)
          .calculatePaymentAmountTest(
            BigNumber.from('0'), // Gas after payment
            0, // Fee PPM
            BigNumber.from('1000000000000000000'),
          ),
      ).to.be.revertedWith(`InvalidLinkWeiPrice(0)`)
      const mockLinkEthNegative = await mockAggregatorV3Factory.deploy(0, -1)
      const vrfCoordinatorV2TestHelperNegative =
        await vrfCoordinatorV2TestHelperFactory.deploy(
          linkToken.address,
          blockHashStore.address,
          mockLinkEthNegative.address,
        )
      await expect(
        vrfCoordinatorV2TestHelperNegative
          .connect(owner)
          .calculatePaymentAmountTest(
            BigNumber.from('0'), // Gas after payment
            0, // Fee PPM
            BigNumber.from('1000000000000000000'),
          ),
      ).to.be.revertedWith(`InvalidLinkWeiPrice(-1)`)
    })
  })

  describe('#keyRegistration', async function () {
    it('register key emits log', async function () {
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      const kh = await vrfCoordinatorV2.hashOfKey(testKey)
      await expect(
        vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey),
      )
        .to.emit(vrfCoordinatorV2, 'ProvingKeyRegistered')
        .withArgs(kh, subOwnerAddress)
      const reqConfig = await vrfCoordinatorV2.getRequestConfig()
      assert(reqConfig[2].length == 1) // 1 keyhash registered
    })
    it('cannot re-register key', async function () {
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      const kh = await vrfCoordinatorV2.hashOfKey(testKey)
      await vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey)
      await expect(
        vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey),
      ).to.be.revertedWith(`ProvingKeyAlreadyRegistered("${kh}")`)
    })
    it('deregister key emits log', async function () {
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      const kh = await vrfCoordinatorV2.hashOfKey(testKey)
      await vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey)
      await expect(vrfCoordinatorV2.deregisterProvingKey(testKey))
        .to.emit(vrfCoordinatorV2, 'ProvingKeyDeregistered')
        .withArgs(kh, subOwnerAddress)
      const reqConfig = await vrfCoordinatorV2.getRequestConfig()
      assert(reqConfig[2].length == 0) // 0 keyhash registered
    })
    it('cannot deregister unregistered key', async function () {
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      const kh = await vrfCoordinatorV2.hashOfKey(testKey)
      await expect(
        vrfCoordinatorV2.deregisterProvingKey(testKey),
      ).to.be.revertedWith(`NoSuchProvingKey("${kh}")`)
    })
    it('can register after deregister', async function () {
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      await vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey)
      await vrfCoordinatorV2.deregisterProvingKey(testKey)
      await vrfCoordinatorV2.registerProvingKey(randomAddress, testKey)
    })
  })

  describe('#fulfillRandomWords', async function () {
    beforeEach(async () => {
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      await vrfCoordinatorV2.registerProvingKey(subOwnerAddress, testKey)
    })
    it('unregistered key should fail', async function () {
      const proof = [
        [BigNumber.from('1'), BigNumber.from('3')], // pk NOT registered
        [BigNumber.from('1'), BigNumber.from('2')], // gamma
        BigNumber.from('1'), // c
        BigNumber.from('1'), // s
        BigNumber.from('1'), // seed
        randomAddress, // uWitness
        [BigNumber.from('1'), BigNumber.from('2')], // cGammaWitness
        [BigNumber.from('1'), BigNumber.from('2')], // sHashWitness
        BigNumber.from('1'),
      ] // 13 words in proof
      const rc = [
        1, // blockNum
        2, // subId
        3, // callbackGasLimit
        4, // numWords
        randomAddress, // sender
      ]
      await expect(
        vrfCoordinatorV2.connect(oracle).fulfillRandomWords(proof, rc),
      ).to.be.revertedWith(
        `NoSuchProvingKey("0xa15bc60c955c405d20d9149c709e2460f1c2d9a497496a7f46004d1772c3054c")`,
      )
    })
    it('no corresponding request', async function () {
      const proof = [
        [BigNumber.from('1'), BigNumber.from('2')], // pk
        [BigNumber.from('1'), BigNumber.from('2')], // gamma
        BigNumber.from('1'), // c
        BigNumber.from('1'), // s
        BigNumber.from('1'), // seed
        randomAddress, // uWitness
        [BigNumber.from('1'), BigNumber.from('2')], // cGammaWitness
        [BigNumber.from('1'), BigNumber.from('2')], // sHashWitness
        BigNumber.from('1'),
      ] // 13 words in proof
      const rc = [
        1, // blockNum
        2, // subId
        3, // callbackGasLimit
        4, // numWords
        randomAddress, // sender
      ]
      await expect(
        vrfCoordinatorV2.connect(oracle).fulfillRandomWords(proof, rc),
      ).to.be.revertedWith(`NoCorrespondingRequest()`)
    })
    it('incorrect commitment wrong blocknum', async function () {
      const subId = await createSubscription()
      await linkToken.connect(subOwner).transferAndCall(
        vrfCoordinatorV2.address,
        BigNumber.from('1000000000000000000'), // 1 link > 0.1 min.
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
      const testKey = [BigNumber.from('1'), BigNumber.from('2')]
      const kh = await vrfCoordinatorV2.hashOfKey(testKey)
      const tx = await vrfCoordinatorV2.connect(consumer).requestRandomWords(
        kh, // keyhash
        subId, // subId
        1, // minReqConf
        1000, // callbackGasLimit
        1, // numWords
      )
      const reqReceipt = await tx.wait()
      // We give it the right proof length and a valid preSeed
      // but an invalid commitment
      const preSeed = reqReceipt.events[0].args['preSeed']
      const proof = [
        [BigNumber.from('1'), BigNumber.from('2')],
        [BigNumber.from('1'), BigNumber.from('2')],
        BigNumber.from('1'),
        BigNumber.from('1'),
        preSeed,
        randomAddress,
        [BigNumber.from('1'), BigNumber.from('2')],
        [BigNumber.from('1'), BigNumber.from('2')],
        BigNumber.from('1'),
      ]
      const rc = [
        reqReceipt.blockNumber + 1, // Wrong blocknumber
        subId,
        1000,
        1,
        await consumer.getAddress(),
      ]
      await expect(
        vrfCoordinatorV2.connect(oracle).fulfillRandomWords(proof, rc),
      ).to.be.revertedWith(`IncorrectCommitment()`)
    })
  })

  describe('#getFeeTier', async function () {
    beforeEach(async () => {
      await expect(
        vrfCoordinatorV2
          .connect(owner)
          .setConfig(
            c.minimumRequestBlockConfirmations,
            c.maxGasLimit,
            c.stalenessSeconds,
            c.gasAfterPaymentCalculation,
            c.weiPerUnitLink,
            [10000, 1000, 100, 10, 1, 10, 20, 30, 40],
          ),
      )
    })
    it('tier1', async function () {
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(0)) == 10000)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(5)) == 10000)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(10)) == 10000)
    })
    it('tier2', async function () {
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(11)) == 1000)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(12)) == 1000)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(20)) == 1000)
    })
    it('tier3', async function () {
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(21)) == 100)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(30)) == 100)
    })
    it('tier4', async function () {
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(31)) == 10)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(40)) == 10)
    })
    it('tier5', async function () {
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(41)) == 1)
      assert((await vrfCoordinatorV2.connect(random).getFeeTier(123102)) == 1)
    })
  })

  /*
    Note that all the fulfillment happy path testing is done in Go, to make use of the existing go code to produce
    proofs offchain.
   */
})
