import { assert, expect } from 'chai'
import { BigNumber, Contract, Signer } from 'ethers'
import { ethers } from 'hardhat'

describe('VRFCoordinatorV2Mock', () => {
  let vrfCoordinatorV2Mock: Contract
  let vrfConsumerV2: Contract
  let linkToken: Contract
  let subOwner: Signer
  let random: Signer
  let subOwnerAddress: string
  let pointOneLink = BigNumber.from('100000000000000000')
  let oneLink = BigNumber.from('1000000000000000000')
  let keyhash =
    '0xe90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0'
  let testConsumerAddress = '0x1111000000000000000000000000000000001111'
  let testConsumerAddress2 = '0x1111000000000000000000000000000000001110'

  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    subOwner = accounts[1]
    subOwnerAddress = await subOwner.getAddress()
    random = accounts[2]

    const vrfCoordinatorV2MockFactory = await ethers.getContractFactory(
      'src/v0.8/mocks/VRFCoordinatorV2Mock.sol:VRFCoordinatorV2Mock',
      accounts[0],
    )
    vrfCoordinatorV2Mock = await vrfCoordinatorV2MockFactory.deploy(
      pointOneLink,
      1e9, // 0.000000001 LINK per gas
    )

    const ltFactory = await ethers.getContractFactory(
      'src/v0.4/LinkToken.sol:LinkToken',
      accounts[0],
    )
    linkToken = await ltFactory.deploy()

    const vrfConsumerV2Factory = await ethers.getContractFactory(
      'src/v0.8/tests/VRFConsumerV2.sol:VRFConsumerV2',
      accounts[0],
    )
    vrfConsumerV2 = await vrfConsumerV2Factory.deploy(
      vrfCoordinatorV2Mock.address,
      linkToken.address,
    )
  })

  async function createSubscription(): Promise<number> {
    const tx = await vrfCoordinatorV2Mock.connect(subOwner).createSubscription()
    const receipt = await tx.wait()
    return receipt.events[0].args['subId']
  }

  describe('#createSubscription', async function () {
    it('can create a subscription', async function () {
      await expect(vrfCoordinatorV2Mock.connect(subOwner).createSubscription())
        .to.emit(vrfCoordinatorV2Mock, 'SubscriptionCreated')
        .withArgs(1, subOwnerAddress)
      const s = await vrfCoordinatorV2Mock.getSubscription(1)
      assert(s.balance.toString() == '0', 'invalid balance')
      assert(s.owner == subOwnerAddress, 'invalid address')
    })
    it('subscription id increments', async function () {
      await expect(vrfCoordinatorV2Mock.connect(subOwner).createSubscription())
        .to.emit(vrfCoordinatorV2Mock, 'SubscriptionCreated')
        .withArgs(1, subOwnerAddress)
      await expect(vrfCoordinatorV2Mock.connect(subOwner).createSubscription())
        .to.emit(vrfCoordinatorV2Mock, 'SubscriptionCreated')
        .withArgs(2, subOwnerAddress)
    })
  })
  describe('#addConsumer', async function () {
    it('can add a consumer to a subscription', async function () {
      let subId = await createSubscription()
      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .addConsumer(subId, testConsumerAddress),
      )
        .to.emit(vrfCoordinatorV2Mock, 'ConsumerAdded')
        .withArgs(subId, testConsumerAddress)
      let sub = await vrfCoordinatorV2Mock
        .connect(subOwner)
        .getSubscription(subId)
      expect(sub.consumers).to.eql([testConsumerAddress])
    })
    it('cannot add a consumer to a nonexistent subscription', async function () {
      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .addConsumer(4, testConsumerAddress),
      ).to.be.revertedWith('InvalidSubscription')
    })
    it('cannot add more than the consumer maximum', async function () {
      let subId = await createSubscription()
      for (let i = 0; i < 100; i++) {
        const testIncrementingAddress = BigNumber.from(i)
          .add('0x1000000000000000000000000000000000000000')
          .toHexString()
        await expect(
          vrfCoordinatorV2Mock
            .connect(subOwner)
            .addConsumer(subId, testIncrementingAddress),
        ).to.emit(vrfCoordinatorV2Mock, 'ConsumerAdded')
      }
      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .addConsumer(subId, testConsumerAddress),
      ).to.be.revertedWith('TooManyConsumers')
    })
  })
  describe('#removeConsumer', async function () {
    it('can remove a consumer from a subscription', async function () {
      let subId = await createSubscription()
      for (const addr of [testConsumerAddress, testConsumerAddress2]) {
        await expect(
          vrfCoordinatorV2Mock.connect(subOwner).addConsumer(subId, addr),
        ).to.emit(vrfCoordinatorV2Mock, 'ConsumerAdded')
      }

      let sub = await vrfCoordinatorV2Mock
        .connect(subOwner)
        .getSubscription(subId)
      expect(sub.consumers).to.eql([testConsumerAddress, testConsumerAddress2])

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .removeConsumer(subId, testConsumerAddress),
      )
        .to.emit(vrfCoordinatorV2Mock, 'ConsumerRemoved')
        .withArgs(subId, testConsumerAddress)

      sub = await vrfCoordinatorV2Mock.connect(subOwner).getSubscription(subId)
      expect(sub.consumers).to.eql([testConsumerAddress2])
    })
    it('cannot remove a consumer from a nonexistent subscription', async function () {
      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .removeConsumer(4, testConsumerAddress),
      ).to.be.revertedWith('InvalidSubscription')
    })
    it('cannot remove a consumer after it is already removed', async function () {
      let subId = await createSubscription()

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .addConsumer(subId, testConsumerAddress),
      ).to.emit(vrfCoordinatorV2Mock, 'ConsumerAdded')

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .removeConsumer(subId, testConsumerAddress),
      ).to.emit(vrfCoordinatorV2Mock, 'ConsumerRemoved')

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .removeConsumer(subId, testConsumerAddress),
      ).to.be.revertedWith('InvalidConsumer')
    })
  })
  describe('#fundSubscription', async function () {
    it('can fund a subscription', async function () {
      let subId = await createSubscription()
      await expect(
        vrfCoordinatorV2Mock.connect(subOwner).fundSubscription(subId, oneLink),
      )
        .to.emit(vrfCoordinatorV2Mock, 'SubscriptionFunded')
        .withArgs(subId, 0, oneLink)
      let sub = await vrfCoordinatorV2Mock
        .connect(subOwner)
        .getSubscription(subId)
      expect(sub.balance).to.equal(oneLink)
    })
    it('cannot fund a nonexistent subscription', async function () {
      await expect(
        vrfCoordinatorV2Mock.connect(subOwner).fundSubscription(4, oneLink),
      ).to.be.revertedWith('InvalidSubscription')
    })
  })
  describe('#cancelSubscription', async function () {
    it('can cancel a subscription', async function () {
      let subId = await createSubscription()
      await expect(
        vrfCoordinatorV2Mock.connect(subOwner).getSubscription(subId),
      ).to.not.be.reverted

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .cancelSubscription(subId, subOwner.getAddress()),
      ).to.emit(vrfCoordinatorV2Mock, 'SubscriptionCanceled')

      await expect(
        vrfCoordinatorV2Mock.connect(subOwner).getSubscription(subId),
      ).to.be.revertedWith('InvalidSubscription')
    })
  })
  describe('#fulfillRandomWords', async function () {
    it('fails to fulfill without being a valid consumer', async function () {
      let subId = await createSubscription()

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .requestRandomWords(keyhash, subId, 3, 500_000, 2),
      ).to.be.revertedWith('InvalidConsumer')
    })
    it('fails to fulfill with insufficient funds', async function () {
      let subId = await createSubscription()
      await vrfCoordinatorV2Mock
        .connect(subOwner)
        .addConsumer(subId, await subOwner.getAddress())

      await expect(
        vrfCoordinatorV2Mock
          .connect(subOwner)
          .requestRandomWords(keyhash, subId, 3, 500_000, 2),
      )
        .to.emit(vrfCoordinatorV2Mock, 'RandomWordsRequested')
        .withArgs(keyhash, 1, 100, subId, 3, 500_000, 2, subOwnerAddress)

      await expect(
        vrfCoordinatorV2Mock
          .connect(random)
          .fulfillRandomWords(1, vrfConsumerV2.address),
      ).to.be.revertedWith('InsufficientBalance')
    })
    it('can request and fulfill [ @skip-coverage ]', async function () {
      let subId = await createSubscription()
      await vrfCoordinatorV2Mock
        .connect(subOwner)
        .addConsumer(subId, vrfConsumerV2.address)
      await expect(
        vrfCoordinatorV2Mock.connect(subOwner).fundSubscription(subId, oneLink),
      ).to.not.be.reverted

      // Call requestRandomWords from the consumer contract so that the requestId
      // member variable on the consumer is appropriately set.
      expect(
        await vrfConsumerV2
          .connect(subOwner)
          .testRequestRandomness(keyhash, subId, 3, 500_000, 2),
      )
        .to.emit(vrfCoordinatorV2Mock, 'RandomWordsRequested')
        .withArgs(keyhash, 1, 100, subId, 3, 500_000, 2, vrfConsumerV2.address)

      let tx = await vrfCoordinatorV2Mock
        .connect(random)
        .fulfillRandomWords(1, vrfConsumerV2.address)
      let receipt = await tx.wait()
      expect(receipt.events[0].event).to.equal('RandomWordsFulfilled')
      expect(receipt.events[0].args['requestId']).to.equal(1)
      expect(receipt.events[0].args['outputSeed']).to.equal(1)
      expect(receipt.events[0].args['success']).to.equal(true)
      assert(
        receipt.events[0].args['payment']
          .sub(BigNumber.from('100119017000000000'))
          .lt(BigNumber.from('10000000000')),
      )

      // Check that balance was subtracted
      let sub = await vrfCoordinatorV2Mock
        .connect(random)
        .getSubscription(subId)
      expect(sub.balance).to.equal(
        oneLink.sub(receipt.events[0].args['payment']),
      )
    })
  })
})
