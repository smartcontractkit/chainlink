import { ethers } from 'hardhat'
import { expect } from 'chai'
import { BigNumber, Contract, ContractFactory, Signer } from 'ethers'
import { Roles, getUsers } from '../../../test-helpers/setup'
import { randomAddressString } from 'hardhat/internal/hardhat-network/provider/utils/random'
import { stringToBytes } from '../../../test-helpers/helpers'

let functionsRouterFactory: ContractFactory
let functionsCoordinatorFactory: ContractFactory
let clientTestHelperFactory: ContractFactory
let linkTokenFactory: ContractFactory
let mockAggregatorV3Factory: ContractFactory
let roles: Roles
let subOwner: Signer
let subOwnerAddress: string
let consumer: Signer
let consumerAddress: string
let stranger: Signer
let strangerAddress: string

const anyValue = () => true

const stringToHex = (s: string) => {
  return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(s))
}

const encodeReport = (requestId: string, result: string, err: string) => {
  const abi = ethers.utils.defaultAbiCoder
  return abi.encode(
    ['bytes32[]', 'bytes[]', 'bytes[]'],
    [[requestId], [result], [err]],
  )
}

const linkEth = BigNumber.from(5021530000000000)

type routerConfig = {
  maxCallbackGasLimit: number
  feedStalenessSeconds: number
  gasOverheadBeforeCallback: number
  gasOverheadAfterCallback: number
  requestTimeoutSeconds: number
  donFee: number
  fallbackNativePerUnitLink: BigNumber
  maxSupportedRequestDataVersion: number
}
const config: routerConfig = {
  maxCallbackGasLimit: 1_000_000,
  feedStalenessSeconds: 86_400,
  gasOverheadBeforeCallback:
    21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
  gasOverheadAfterCallback:
    21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
  requestTimeoutSeconds: 300,
  donFee: 0,
  fallbackNativePerUnitLink: BigNumber.from(5000000000000000),
  maxSupportedRequestDataVersion: 1,
}

before(async () => {
  roles = (await getUsers()).roles

  functionsRouterFactory = await ethers.getContractFactory(
    'src/v0.8/functions/dev/1_0_0/FunctionsRouter.sol:FunctionsRouter',
    roles.defaultAccount,
  )

  functionsCoordinatorFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/1_0_0/testhelpers/FunctionsCoordinatorTestHelper.sol:FunctionsCoordinatorTestHelper',
    roles.consumer,
  )

  clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/1_0_0/testhelpers/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
  )

  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    roles.consumer,
  )

  mockAggregatorV3Factory = await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
    roles.consumer,
  )
})

describe('Functionsrouter', () => {
  const routerLabel = ethers.utils.formatBytes32String('')
  let router: Contract
  let coordinator: Contract
  const donLabel = ethers.utils.formatBytes32String('1')

  let client: Contract
  let linkToken: Contract
  let mockLinkEth: Contract

  beforeEach(async () => {
    const { roles } = await getUsers()
    subOwner = roles.consumer
    subOwnerAddress = await subOwner.getAddress()
    consumer = roles.consumer2
    consumerAddress = await consumer.getAddress()
    stranger = roles.stranger
    strangerAddress = await stranger.getAddress()

    // Deploy
    linkToken = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth)
    router = await functionsRouterFactory
      .connect(roles.defaultAccount)
      .deploy(
        0,
        0,
        linkToken.address,
        ethers.utils.defaultAbiCoder.encode(
          ['uint96', 'bytes4'],
          [0, 0x0ca76175],
        ),
      )
    coordinator = await functionsCoordinatorFactory
      .connect(roles.defaultAccount)
      .deploy(
        router.address,
        ethers.utils.defaultAbiCoder.encode(
          [
            'uint32',
            'uint32',
            'uint32',
            'uint32',
            'int256',
            'uint32',
            'uint96',
            'uint16',
          ],
          [...Object.values(config)],
        ),
        mockLinkEth.address,
      )
    client = await clientTestHelperFactory
      .connect(roles.consumer)
      .deploy(router.address)

    // Setup accounts
    await linkToken.transfer(
      subOwnerAddress,
      BigNumber.from('1000000000000000000'), // 1 LINK
    )
    await linkToken.transfer(
      strangerAddress,
      BigNumber.from('1000000000000000000'), // 1 LINK
    )
  })

  describe('Config', () => {
    it('non-owner is unable set config', async () => {
      await expect(
        router
          .connect(roles.stranger)
          .proposeContractsUpdate(
            [donLabel],
            [
              ethers.constants.AddressZero,
              ethers.constants.AddressZero,
              ethers.constants.AddressZero,
            ],
            [coordinator.address],
          ),
      ).to.be.revertedWith('Only callable by owner')
    })

    it('Owner can update config on the Router', async () => {
      const [
        beforeSystemVersionMajor,
        beforeSystemVersionMinor,
        beforeSystemVersionPatch,
      ] = await router.version()
      const beforeConfigHash = await router.getConfigHash()

      await expect(
        router.proposeConfigUpdate(
          routerLabel,
          ethers.utils.defaultAbiCoder.encode(
            ['uint96', 'bytes4'],
            [1, 0x0ca76175],
          ),
        ),
      ).to.emit(router, 'ConfigProposed')
      await expect(router.updateConfig(routerLabel)).to.emit(
        router,
        'ConfigUpdated',
      )
      const [
        afterSystemVersionMajor,
        afterSystemVersionMinor,
        afterSystemVersionPatch,
      ] = await router.version()
      const afterConfigHash = await router.getConfigHash()
      expect(afterSystemVersionMajor).to.equal(beforeSystemVersionMajor)
      expect(afterSystemVersionMinor).to.equal(beforeSystemVersionMinor)
      expect(afterSystemVersionPatch).to.equal(beforeSystemVersionPatch + 1)
      expect(beforeConfigHash).to.not.equal(afterConfigHash)
    })

    it('Config of a contract on a route can be updated', async () => {
      await router.proposeContractsUpdate(
        [donLabel],
        [ethers.constants.AddressZero],
        [coordinator.address],
      )
      await router.updateContracts()
      const [
        beforeSystemVersionMajor,
        beforeSystemVersionMinor,
        beforeSystemVersionPatch,
      ] = await router.version()
      const beforeConfigHash = await coordinator.getConfigHash()

      await expect(
        router.proposeConfigUpdate(
          donLabel,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
            ],
            [
              ...Object.values({
                ...config,
                maxSupportedRequestDataVersion: 2,
              }),
            ],
          ),
        ),
      ).to.emit(router, 'ConfigProposed')
      await expect(router.updateConfig(donLabel)).to.emit(
        router,
        'ConfigUpdated',
      )
      const [
        afterSystemVersionMajor,
        afterSystemVersionMinor,
        afterSystemVersionPatch,
      ] = await router.version()
      const afterConfigHash = await router.getConfigHash()
      expect(afterSystemVersionMajor).to.equal(beforeSystemVersionMajor)
      expect(afterSystemVersionMinor).to.equal(beforeSystemVersionMinor)
      expect(afterSystemVersionPatch).to.equal(beforeSystemVersionPatch + 1)
      expect(beforeConfigHash).to.not.equal(afterConfigHash)
    })

    it('returns the config set on the Router', async () => {
      await router.connect(roles.stranger).getAdminFee()
    })
  })

  describe('Updates', () => {
    it('One or more contracts on a route can be updated', async () => {
      const coordinator2 = await functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          router.address,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
            ],
            [...Object.values(config)],
          ),
          mockLinkEth.address,
        )
      const donLabel2 = ethers.utils.formatBytes32String('2')
      const coordinator3 = await functionsCoordinatorFactory
        .connect(roles.defaultAccount)
        .deploy(
          router.address,
          ethers.utils.defaultAbiCoder.encode(
            [
              'uint32',
              'uint32',
              'uint32',
              'uint32',
              'int256',
              'uint32',
              'uint96',
              'uint16',
            ],
            [...Object.values(config)],
          ),
          mockLinkEth.address,
        )
      const donLabel3 = ethers.utils.formatBytes32String('3')

      const [
        beforeSystemVersionMajor,
        beforeSystemVersionMinor,
        beforeSystemVersionPatch,
      ] = await router.version()
      await expect(
        router['getContractById(bytes32)'](donLabel),
      ).to.be.revertedWith('RouteNotFound')
      await expect(
        router.proposeContractsUpdate(
          [donLabel, donLabel2, donLabel3],
          [
            ethers.constants.AddressZero,
            ethers.constants.AddressZero,
            ethers.constants.AddressZero,
          ],
          [coordinator.address, coordinator2.address, coordinator3.address],
        ),
      ).to.emit(router, `ContractProposed`)
      await expect(router.updateContracts()).to.emit(router, 'ContractUpdated')
      expect(await router['getContractById(bytes32)'](donLabel)).to.equal(
        coordinator.address,
      )
      const [
        afterSystemVersionMajor,
        afterSystemVersionMinor,
        afterSystemVersionPatch,
      ] = await router.version()
      expect(afterSystemVersionMajor).to.equal(beforeSystemVersionMajor)
      expect(afterSystemVersionMinor).to.equal(beforeSystemVersionMinor + 1)
      expect(afterSystemVersionPatch).to.equal(beforeSystemVersionPatch)
    })

    it('non-owner is unable to register a DON', async () => {
      await expect(
        router
          .connect(roles.stranger)
          .proposeContractsUpdate(
            [donLabel],
            [ethers.constants.AddressZero],
            [coordinator.address],
          ),
      ).to.be.revertedWith('Only callable by owner')
    })
  })

  async function createSubscription(
    owner: Signer,
    consumers: string[],
  ): Promise<number> {
    const tx = await router.connect(owner).createSubscription()
    const receipt = await tx.wait()
    const subId = receipt.events[0].args['subscriptionId'].toNumber()
    for (let i = 0; i < consumers.length; i++) {
      await router.connect(owner).addConsumer(subId, consumers[i])
    }
    return subId
  }

  describe('Subscription management', () => {
    describe('#createSubscription', async function () {
      it('can create a subscription', async function () {
        await expect(router.connect(subOwner).createSubscription())
          .to.emit(router, 'SubscriptionCreated')
          .withArgs(1, subOwnerAddress)
        const s = await router.getSubscription(1)
        expect(s.balance.toString() == '0', 'invalid balance')
        expect(s.owner == subOwnerAddress, 'invalid address')
      })
      it('subscription id increments', async function () {
        await expect(router.connect(subOwner).createSubscription())
          .to.emit(router, 'SubscriptionCreated')
          .withArgs(1, subOwnerAddress)
        await expect(router.connect(subOwner).createSubscription())
          .to.emit(router, 'SubscriptionCreated')
          .withArgs(2, subOwnerAddress)
      })
      it('cannot create more than the max', async function () {
        const subId = createSubscription(subOwner, [])
        for (let i = 0; i < 100; i++) {
          await router
            .connect(subOwner)
            .addConsumer(subId, randomAddressString())
        }
        await expect(
          router.connect(subOwner).addConsumer(subId, randomAddressString()),
        ).to.be.revertedWith(`TooManyConsumers()`)
      })
    })

    describe('#requestSubscriptionOwnerTransfer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('rejects non-owner', async function () {
        await expect(
          router
            .connect(roles.stranger)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('owner can request transfer', async function () {
        await expect(
          router
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        )
          .to.emit(router, 'SubscriptionOwnerTransferRequested')
          .withArgs(subId, subOwnerAddress, strangerAddress)
        // Same request is a noop
        await expect(
          router
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        ).to.not.emit(router, 'SubscriptionOwnerTransferRequested')
      })
    })

    describe('#acceptSubscriptionOwnerTransfer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('subscription must exist', async function () {
        await expect(
          router.connect(subOwner).acceptSubscriptionOwnerTransfer(1203123123),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be requested owner to accept', async function () {
        await expect(
          router
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        )
        await expect(
          router.connect(subOwner).acceptSubscriptionOwnerTransfer(subId),
        ).to.be.revertedWith(`MustBeRequestedOwner("${strangerAddress}")`)
      })
      it('requested owner can accept', async function () {
        await expect(
          router
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        )
          .to.emit(router, 'SubscriptionOwnerTransferRequested')
          .withArgs(subId, subOwnerAddress, strangerAddress)
        await expect(
          router.connect(stranger).acceptSubscriptionOwnerTransfer(subId),
        )
          .to.emit(router, 'SubscriptionOwnerTransferred')
          .withArgs(subId, subOwnerAddress, strangerAddress)
      })
    })

    describe('#addConsumer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('subscription must exist', async function () {
        await expect(
          router.connect(subOwner).addConsumer(1203123123, strangerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          router.connect(stranger).addConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('add is idempotent', async function () {
        await router.connect(subOwner).addConsumer(subId, strangerAddress)
        await router.connect(subOwner).addConsumer(subId, strangerAddress)
      })
      it('cannot add more than maximum', async function () {
        // There is one consumer, add another 99 to hit the max
        for (let i = 0; i < 99; i++) {
          await router
            .connect(subOwner)
            .addConsumer(subId, randomAddressString())
        }
        // Adding one more should fail
        // await router.connect(subOwner).addConsumer(subId, strangerAddress);
        await expect(
          router.connect(subOwner).addConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`TooManyConsumers()`)
        // Same is true if we first create with the maximum
        const consumers: string[] = []
        for (let i = 0; i < 100; i++) {
          consumers.push(randomAddressString())
        }
        subId = await createSubscription(subOwner, consumers)
        await expect(
          router.connect(subOwner).addConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`TooManyConsumers()`)
      })
      it('owner can update', async function () {
        await expect(
          router.connect(subOwner).addConsumer(subId, strangerAddress),
        )
          .to.emit(router, 'SubscriptionConsumerAdded')
          .withArgs(subId, strangerAddress)
      })
    })

    describe('#removeConsumer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('subscription must exist', async function () {
        await expect(
          router.connect(subOwner).removeConsumer(1203123123, strangerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          router.connect(stranger).removeConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('owner can update', async function () {
        const subBefore = await router.getSubscription(subId)
        await router.connect(subOwner).addConsumer(subId, strangerAddress)
        await expect(
          router.connect(subOwner).removeConsumer(subId, strangerAddress),
        )
          .to.emit(router, 'SubscriptionConsumerRemoved')
          .withArgs(subId, strangerAddress)
        const subAfter = await router.getSubscription(subId)
        // Subscription should NOT contain the removed consumer
        expect(subBefore.consumers).to.deep.equal(subAfter.consumers)
      })
      it('can remove all consumers', async function () {
        // Testing the handling of zero.
        await router.connect(subOwner).addConsumer(subId, strangerAddress)
        await router.connect(subOwner).removeConsumer(subId, strangerAddress)
        await router.connect(subOwner).removeConsumer(subId, consumerAddress)
        // Should be empty
        const subAfter = await router.getSubscription(subId)
        expect(subAfter.consumers).to.deep.equal([])
      })
    })

    describe('#pendingRequestExists', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])

        await linkToken
          .connect(subOwner)
          .transferAndCall(
            router.address,
            BigNumber.from('130790416713017745'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await router.connect(subOwner).addConsumer(subId, client.address)
      })
      it('returns false when there is no latest pending request', async function () {
        expect(await router.connect(subOwner).pendingRequestExists(subId)).to.be
          .false
      })
      it('returns true when the latest request is pending', async function () {
        await router.proposeContractsUpdate(
          [donLabel],
          [ethers.constants.AddressZero],
          [coordinator.address],
        )
        await router.updateContracts()
        await client
          .connect(consumer)
          .sendSimpleRequestWithJavaScript(
            `return 'hello world'`,
            subId,
            donLabel,
          )
        expect(await router.connect(subOwner).pendingRequestExists(subId)).to.be
          .true
      })
    })

    describe('#cancelSubscription', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('subscription must exist', async function () {
        await expect(
          router
            .connect(subOwner)
            .cancelSubscription(1203123123, subOwnerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          router.connect(stranger).cancelSubscription(subId, subOwnerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('can cancel', async function () {
        await linkToken
          .connect(subOwner)
          .transferAndCall(
            router.address,
            BigNumber.from('1000'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await expect(
          router.connect(subOwner).cancelSubscription(subId, strangerAddress),
        )
          .to.emit(router, 'SubscriptionCanceled')
          .withArgs(subId, strangerAddress, BigNumber.from('1000'))
        const strangerBalance = await linkToken.balanceOf(strangerAddress)
        expect(strangerBalance.toString()).to.equal('1000000000000001000')
        await expect(
          router.connect(subOwner).getSubscription(subId),
        ).to.be.revertedWith('InvalidSubscription')
      })
      it('can add same consumer after canceling', async function () {
        await linkToken
          .connect(subOwner)
          .transferAndCall(
            router.address,
            BigNumber.from('1000'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await router.connect(subOwner).addConsumer(subId, strangerAddress)
        await router
          .connect(subOwner)
          .cancelSubscription(subId, strangerAddress)
        subId = await createSubscription(subOwner, [consumerAddress])
        // The cancel should have removed this consumer, so we can add it again.
        await router.connect(subOwner).addConsumer(subId, strangerAddress)
      })
      it('cannot cancel with pending request', async function () {
        await router.proposeContractsUpdate(
          [donLabel],
          [ethers.constants.AddressZero],
          [coordinator.address],
        )
        await router.updateContracts()
        await linkToken
          .connect(subOwner)
          .transferAndCall(
            router.address,
            BigNumber.from('130790416713017745'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await router.connect(subOwner).addConsumer(subId, client.address)
        await client
          .connect(consumer)
          .sendSimpleRequestWithJavaScript(
            `return 'hello world'`,
            subId,
            donLabel,
          )
        // Should revert with outstanding requests
        await expect(
          router.connect(subOwner).cancelSubscription(subId, strangerAddress),
        ).to.be.revertedWith('PendingRequestExists()')
        // However the owner is able to cancel
        // funds go to the sub owner.
        await expect(
          router.connect(roles.defaultAccount).ownerCancelSubscription(subId),
        )
          .to.emit(router, 'SubscriptionCanceled')
          .withArgs(
            subId,
            subOwnerAddress,
            BigNumber.from('130790416713017745'),
          )
      })
    })

    describe('#recoverFunds', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })

      it('function that should change internal balance do', async function () {
        type bf = [() => Promise<any>, BigNumber]
        const balanceChangingFns: Array<bf> = [
          [
            async function () {
              const s = ethers.utils.defaultAbiCoder.encode(['uint64'], [subId])
              await linkToken
                .connect(subOwner)
                .transferAndCall(router.address, BigNumber.from('1000'), s)
            },
            BigNumber.from('1000'),
          ],
          [
            async function () {
              await router
                .connect(subOwner)
                .cancelSubscription(subId, strangerAddress)
            },
            BigNumber.from('-1000'),
          ],
        ]
        for (const [fn, expectedBalanceChange] of balanceChangingFns) {
          const startingBalance = await router.getTotalBalance()
          await fn()
          const endingBalance = await router.getTotalBalance()
          expect(
            endingBalance.sub(startingBalance).toString() ==
              expectedBalanceChange.toString(),
          )
        }
      })
      it('only owner can recover', async function () {
        await expect(
          router.connect(subOwner).recoverFunds(strangerAddress),
        ).to.be.revertedWith('Only callable by owner')
      })

      it('owner can recover link transferred', async function () {
        // Set the internal balance
        expect(BigNumber.from('0'), linkToken.balanceOf(strangerAddress))
        const s = ethers.utils.defaultAbiCoder.encode(['uint64'], [subId])
        await linkToken
          .connect(subOwner)
          .transferAndCall(router.address, BigNumber.from('1000'), s)
        // Circumvent internal balance
        await linkToken
          .connect(subOwner)
          .transfer(router.address, BigNumber.from('1000'))
        // Should recover this 1000
        await expect(
          router.connect(roles.defaultAccount).recoverFunds(strangerAddress),
        )
          .to.emit(router, 'FundsRecovered')
          .withArgs(strangerAddress, BigNumber.from('1000'))
        expect(BigNumber.from('1000'), linkToken.balanceOf(strangerAddress))
      })
    })
  })

  describe('#oracleWithdraw', async function () {
    it('cannot withdraw with no balance', async function () {
      await expect(
        router
          .connect(roles.oracleNode)
          .oracleWithdraw(randomAddressString(), BigNumber.from('100')),
      ).to.be.revertedWith(`InsufficientSubscriptionBalance`)
    })
  })
})
