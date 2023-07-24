import { ethers } from 'hardhat'
import { expect } from 'chai'
import { BigNumber } from 'ethers'
import { randomAddressString } from 'hardhat/internal/hardhat-network/provider/utils/random'
import {
  getSetupFactory,
  FunctionsContracts,
  FunctionsRoles,
  createSubscription,
  acceptTermsOfService,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let roles: FunctionsRoles

const donLabel = ethers.utils.formatBytes32String('1')

beforeEach(async () => {
  ;({ contracts, roles } = setup())
})

describe('Functions Router - Subscriptions', () => {
  describe('Subscription management', () => {
    describe('#createSubscription', async function () {
      it('can create a subscription', async function () {
        await acceptTermsOfService(
          contracts.accessControl,
          roles.subOwner,
          roles.subOwnerAddress,
        )
        await expect(
          contracts.router.connect(roles.subOwner).createSubscription(),
        )
          .to.emit(contracts.router, 'SubscriptionCreated')
          .withArgs(1, roles.subOwnerAddress)
        const s = await contracts.router.getSubscription(1)
        expect(s.balance.toString()).to.equal('0')
        expect(s.owner).to.equal(roles.subOwnerAddress)
      })
      it('subscription id increments', async function () {
        await acceptTermsOfService(
          contracts.accessControl,
          roles.subOwner,
          roles.subOwnerAddress,
        )
        await expect(
          contracts.router.connect(roles.subOwner).createSubscription(),
        )
          .to.emit(contracts.router, 'SubscriptionCreated')
          .withArgs(1, roles.subOwnerAddress)
        await expect(
          contracts.router.connect(roles.subOwner).createSubscription(),
        )
          .to.emit(contracts.router, 'SubscriptionCreated')
          .withArgs(2, roles.subOwnerAddress)
      })
      it('cannot create more than the max', async function () {
        const subId = createSubscription(
          roles.subOwner,
          [],
          contracts.router,
          contracts.accessControl,
        )
        for (let i = 0; i < 100; i++) {
          await contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, randomAddressString())
        }
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, randomAddressString()),
        ).to.be.revertedWith(`TooManyConsumers()`)
      })
    })

    describe('#requestSubscriptionOwnerTransfer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
      })
      it('rejects non-owner', async function () {
        await expect(
          contracts.router
            .connect(roles.stranger)
            .requestSubscriptionOwnerTransfer(subId, roles.strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${roles.subOwnerAddress}")`)
      })
      it('owner can request transfer', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .requestSubscriptionOwnerTransfer(subId, roles.strangerAddress),
        )
          .to.emit(contracts.router, 'SubscriptionOwnerTransferRequested')
          .withArgs(subId, roles.subOwnerAddress, roles.strangerAddress)
        // Same request is a noop
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .requestSubscriptionOwnerTransfer(subId, roles.strangerAddress),
        ).to.not.emit(contracts.router, 'SubscriptionOwnerTransferRequested')
      })
    })

    describe('#acceptSubscriptionOwnerTransfer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
      })
      it('subscription must exist', async function () {
        // 0x0 is requested owner
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .acceptSubscriptionOwnerTransfer(1203123123),
        ).to.be.revertedWith(`MustBeRequestedOwner`)
      })
      it('must be requested owner to accept', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .requestSubscriptionOwnerTransfer(subId, roles.strangerAddress),
        )
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .acceptSubscriptionOwnerTransfer(subId),
        ).to.be.revertedWith(`MustBeRequestedOwner("${roles.strangerAddress}")`)
      })
      it('requested owner can accept', async function () {
        await acceptTermsOfService(
          contracts.accessControl,
          roles.stranger,
          roles.strangerAddress,
        )
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .requestSubscriptionOwnerTransfer(subId, roles.strangerAddress),
        )
          .to.emit(contracts.router, 'SubscriptionOwnerTransferRequested')
          .withArgs(subId, roles.subOwnerAddress, roles.strangerAddress)
        await expect(
          contracts.router
            .connect(roles.stranger)
            .acceptSubscriptionOwnerTransfer(subId),
        )
          .to.emit(contracts.router, 'SubscriptionOwnerTransferred')
          .withArgs(subId, roles.subOwnerAddress, roles.strangerAddress)
      })
    })

    describe('#addConsumer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
      })
      it('subscription must exist', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(1203123123, roles.strangerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          contracts.router
            .connect(roles.stranger)
            .addConsumer(subId, roles.strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${roles.subOwnerAddress}")`)
      })
      it('add is idempotent', async function () {
        await contracts.router
          .connect(roles.subOwner)
          .addConsumer(subId, roles.strangerAddress)
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, roles.strangerAddress),
        ).to.not.be.reverted
      })
      it('cannot add more than maximum', async function () {
        // There is one consumer, add another 99 to hit the max
        for (let i = 0; i < 99; i++) {
          await contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, randomAddressString())
        }
        // Adding one more should fail
        // await contracts.router.connect(roles.subOwner).addConsumer(subId, roles.strangerAddress);
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, roles.strangerAddress),
        ).to.be.revertedWith(`TooManyConsumers()`)
        // Same is true if we first create with the maximum
        const consumers: string[] = []
        for (let i = 0; i < 100; i++) {
          consumers.push(randomAddressString())
        }
        subId = await createSubscription(
          roles.subOwner,
          consumers,
          contracts.router,
          contracts.accessControl,
        )
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, roles.strangerAddress),
        ).to.be.revertedWith(`TooManyConsumers()`)
      })
      it('owner can update', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, roles.strangerAddress),
        )
          .to.emit(contracts.router, 'SubscriptionConsumerAdded')
          .withArgs(subId, roles.strangerAddress)
      })
    })

    describe('#removeConsumer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
      })
      it('subscription must exist', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .removeConsumer(1203123123, roles.strangerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          contracts.router
            .connect(roles.stranger)
            .removeConsumer(subId, roles.strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${roles.subOwnerAddress}")`)
      })
      it('owner can update', async function () {
        const subBefore = await contracts.router.getSubscription(subId)
        await contracts.router
          .connect(roles.subOwner)
          .addConsumer(subId, roles.strangerAddress)
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .removeConsumer(subId, roles.strangerAddress),
        )
          .to.emit(contracts.router, 'SubscriptionConsumerRemoved')
          .withArgs(subId, roles.strangerAddress)
        const subAfter = await contracts.router.getSubscription(subId)
        // Subscription should NOT contain the removed consumer
        expect(subBefore.consumers).to.deep.equal(subAfter.consumers)
      })
      it('can remove all consumers', async function () {
        // Testing the handling of zero.
        await contracts.router
          .connect(roles.subOwner)
          .addConsumer(subId, roles.strangerAddress)
        await contracts.router
          .connect(roles.subOwner)
          .removeConsumer(subId, roles.strangerAddress)
        await contracts.router
          .connect(roles.subOwner)
          .removeConsumer(subId, roles.consumerAddress)
        // Should be empty
        const subAfter = await contracts.router.getSubscription(subId)
        expect(subAfter.consumers).to.deep.equal([])
      })
    })

    describe('#pendingRequestExists', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )

        await contracts.linkToken
          .connect(roles.subOwner)
          .transferAndCall(
            contracts.router.address,
            BigNumber.from('130790416713017745'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await contracts.router
          .connect(roles.subOwner)
          .addConsumer(subId, contracts.client.address)
      })
      it('returns false when there is no latest pending request', async function () {
        expect(
          await contracts.router
            .connect(roles.subOwner)
            .pendingRequestExists(subId),
        ).to.be.false
      })
      it('returns true when the latest request is pending', async function () {
        await contracts.client
          .connect(roles.consumer)
          .sendSimpleRequestWithJavaScript(
            `return 'hello world'`,
            subId,
            donLabel,
          )
        expect(
          await contracts.router
            .connect(roles.subOwner)
            .pendingRequestExists(subId),
        ).to.be.true
      })
    })

    describe('#cancelSubscription', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
      })
      it('subscription must exist', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .cancelSubscription(1203123123, roles.subOwnerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          contracts.router
            .connect(roles.stranger)
            .cancelSubscription(subId, roles.subOwnerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${roles.subOwnerAddress}")`)
      })
      it('can cancel', async function () {
        await contracts.linkToken
          .connect(roles.subOwner)
          .transferAndCall(
            contracts.router.address,
            BigNumber.from('1000'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .cancelSubscription(subId, roles.strangerAddress),
        )
          .to.emit(contracts.router, 'SubscriptionCanceled')
          .withArgs(subId, roles.strangerAddress, BigNumber.from('1000'))
        const strangerBalance = await contracts.linkToken.balanceOf(
          roles.strangerAddress,
        )
        expect(strangerBalance.toString()).to.equal('1000000000000001000')
        await expect(
          contracts.router.connect(roles.subOwner).getSubscription(subId),
        ).to.be.revertedWith('InvalidSubscription')
      })
      it('can add same consumer after canceling', async function () {
        await contracts.linkToken
          .connect(roles.subOwner)
          .transferAndCall(
            contracts.router.address,
            BigNumber.from('1000'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await contracts.router
          .connect(roles.subOwner)
          .addConsumer(subId, roles.strangerAddress)
        await contracts.router
          .connect(roles.subOwner)
          .cancelSubscription(subId, roles.strangerAddress)
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
        // The cancel should have removed this consumer, so we can add it again.
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .addConsumer(subId, roles.strangerAddress),
        ).to.not.be.reverted
      })
      it('cannot cancel with pending request', async function () {
        await contracts.linkToken
          .connect(roles.subOwner)
          .transferAndCall(
            contracts.router.address,
            BigNumber.from('130790416713017745'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await contracts.router
          .connect(roles.subOwner)
          .addConsumer(subId, contracts.client.address)
        await contracts.client
          .connect(roles.consumer)
          .sendSimpleRequestWithJavaScript(
            `return 'hello world'`,
            subId,
            donLabel,
          )
        // Should revert with outstanding requests
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .cancelSubscription(subId, roles.strangerAddress),
        ).to.be.revertedWith('PendingRequestExists()')
        // However the owner is able to cancel
        // funds go to the sub owner.
        await expect(
          contracts.router
            .connect(roles.defaultAccount)
            .ownerCancelSubscription(subId),
        )
          .to.emit(contracts.router, 'SubscriptionCanceled')
          .withArgs(
            subId,
            roles.subOwnerAddress,
            BigNumber.from('130790416713017745'),
          )
      })
    })

    describe('#recoverFunds', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(
          roles.subOwner,
          [roles.consumerAddress],
          contracts.router,
          contracts.accessControl,
        )
      })

      it('function that should change internal balance do', async function () {
        type bf = [() => Promise<any>, BigNumber]
        const balanceChangingFns: Array<bf> = [
          [
            async function () {
              const s = ethers.utils.defaultAbiCoder.encode(['uint64'], [subId])
              await contracts.linkToken
                .connect(roles.subOwner)
                .transferAndCall(
                  contracts.router.address,
                  BigNumber.from('1000'),
                  s,
                )
            },
            BigNumber.from('1000'),
          ],
          [
            async function () {
              await contracts.router
                .connect(roles.subOwner)
                .cancelSubscription(subId, roles.strangerAddress)
            },
            BigNumber.from('-1000'),
          ],
        ]
        for (const [fn, expectedBalanceChange] of balanceChangingFns) {
          const startingBalance = await contracts.router.getTotalBalance()
          await fn()
          const endingBalance = await contracts.router.getTotalBalance()
          expect(endingBalance.sub(startingBalance.toString())).to.equal(
            expectedBalanceChange.toString(),
          )
        }
      })
      it('only owner can recover', async function () {
        await expect(
          contracts.router
            .connect(roles.subOwner)
            .recoverFunds(roles.strangerAddress),
        ).to.be.revertedWith('Only callable by owner')
      })

      it('owner can recover link transferred', async function () {
        // Set the internal balance
        expect(
          await contracts.linkToken.balanceOf(roles.strangerAddress),
        ).to.equal(BigNumber.from('1000000000000000000'))
        const subscription = ethers.utils.defaultAbiCoder.encode(
          ['uint64'],
          [subId],
        )
        await contracts.linkToken
          .connect(roles.subOwner)
          .transferAndCall(
            contracts.router.address,
            BigNumber.from('1000'),
            subscription,
          )
        // Circumvent internal balance
        await contracts.linkToken
          .connect(roles.subOwner)
          .transfer(contracts.router.address, BigNumber.from('1000'))
        // Should recover this 1000
        await expect(
          contracts.router
            .connect(roles.defaultAccount)
            .recoverFunds(roles.strangerAddress),
        )
          .to.emit(contracts.router, 'FundsRecovered')
          .withArgs(roles.strangerAddress, BigNumber.from('1000'))
        expect(
          await contracts.linkToken.balanceOf(roles.strangerAddress),
        ).to.equal(BigNumber.from('1000000000000001000'))
      })
    })
  })

  describe('#oracleWithdraw', async function () {
    it('cannot withdraw with no balance', async function () {
      await expect(
        contracts.router
          .connect(roles.oracleNode)
          .oracleWithdraw(randomAddressString(), BigNumber.from('100')),
      ).to.be.revertedWith(`InsufficientBalance`)
    })
  })
})
