import { ethers } from 'hardhat'
import { expect } from 'chai'
import { BigNumber, Contract, ContractFactory, Signer } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'
import { randomAddressString } from 'hardhat/internal/hardhat-network/provider/utils/random'
import { stringToBytes } from '../../test-helpers/helpers'

let functionsOracleFactory: ContractFactory
let clientTestHelperFactory: ContractFactory
let functionsBillingRegistryFactory: ContractFactory
let linkTokenFactory: ContractFactory
let mockAggregatorV3Factory: ContractFactory
let roles: Roles
let subOwner: Signer
let subOwnerAddress: string
let consumer: Signer
let consumerAddress: string
let stranger: Signer
let strangerAddress: string

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

type RegistryConfig = {
  maxGasLimit: number
  stalenessSeconds: number
  gasAfterPaymentCalculation: number
  weiPerUnitLink: BigNumber
  gasOverhead: number
  requestTimeoutSeconds: number
}
const config: RegistryConfig = {
  maxGasLimit: 1_000_000,
  stalenessSeconds: 86_400,
  gasAfterPaymentCalculation:
    21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
  weiPerUnitLink: BigNumber.from('5000000000000000'),
  gasOverhead: 100_000,
  requestTimeoutSeconds: 300,
}

before(async () => {
  roles = (await getUsers()).roles

  functionsOracleFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsOracleHelper.sol:FunctionsOracleHelper',
    roles.defaultAccount,
  )

  clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
  )

  functionsBillingRegistryFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsBillingRegistryWithInit.sol:FunctionsBillingRegistryWithInit',
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

describe('FunctionsRegistry', () => {
  let registry: Contract
  let oracle: Contract
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
    oracle = await functionsOracleFactory.connect(roles.defaultAccount).deploy()
    registry = await functionsBillingRegistryFactory
      .connect(roles.defaultAccount)
      .deploy(linkToken.address, mockLinkEth.address, oracle.address)
    client = await clientTestHelperFactory
      .connect(roles.consumer)
      .deploy(oracle.address)

    // Setup contracts
    await oracle.setRegistry(registry.address)
    await oracle.deactivateAuthorizedReceiver()

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

  // NOTE: Temporarily disabled until contract size can be reduced in another way
  // describe('General', () => {
  //   it('#typeAndVersion', async () => {
  //     expect(await registry.callStatic.typeAndVersion()).to.be.equal(
  //       'FunctionsBillingRegistry 0.0.0',
  //     )
  //   })
  // })

  describe('Config', () => {
    it('non-owner is unable set config', async () => {
      await expect(
        registry
          .connect(roles.stranger)
          .setConfig(
            config.maxGasLimit,
            config.stalenessSeconds,
            config.gasAfterPaymentCalculation,
            config.weiPerUnitLink,
            config.gasOverhead,
            config.requestTimeoutSeconds,
          ),
      ).to.be.revertedWith('OnlyCallableByOwner()')
    })

    it('owner can set config', async () => {
      await expect(
        registry
          .connect(roles.defaultAccount)
          .setConfig(
            config.maxGasLimit,
            config.stalenessSeconds,
            config.gasAfterPaymentCalculation,
            config.weiPerUnitLink,
            config.gasOverhead,
            config.requestTimeoutSeconds,
          ),
      ).not.to.be.reverted
    })

    it('returns the config set on the registry', async () => {
      await registry
        .connect(roles.defaultAccount)
        .setConfig(
          config.maxGasLimit,
          config.stalenessSeconds,
          config.gasAfterPaymentCalculation,
          config.weiPerUnitLink,
          config.gasOverhead,
          config.requestTimeoutSeconds,
        )

      const [
        maxGasLimit,
        stalenessSeconds,
        gasAfterPaymentCalculation,
        weiPerUnitLink,
        gasOverhead,
      ] = await registry.connect(roles.stranger).getConfig()

      await expect(config.maxGasLimit).to.equal(maxGasLimit)
      await expect(config.stalenessSeconds).to.equal(stalenessSeconds)
      await expect(config.gasAfterPaymentCalculation).to.equal(
        gasAfterPaymentCalculation,
      )
      await expect(config.weiPerUnitLink).to.equal(weiPerUnitLink)
      await expect(config.gasOverhead).to.equal(gasOverhead)
    })
  })

  describe('DON registration', () => {
    it('non-owner is unable to register a DON', async () => {
      await expect(
        registry.connect(roles.stranger).setAuthorizedSenders([oracle.address]),
      ).to.be.revertedWith('OnlyCallableByOwner()')
    })

    it('owner can register a DON', async () => {
      await expect(
        registry
          .connect(roles.defaultAccount)
          .setAuthorizedSenders([oracle.address]),
      ).not.to.be.reverted
    })
  })

  async function createSubscription(
    owner: Signer,
    consumers: string[],
  ): Promise<number> {
    const tx = await registry.connect(owner).createSubscription()
    const receipt = await tx.wait()
    const subId = receipt.events[0].args['subscriptionId'].toNumber()
    for (let i = 0; i < consumers.length; i++) {
      await registry.connect(owner).addConsumer(subId, consumers[i])
    }
    return subId
  }

  describe('Subscription management', () => {
    describe('#createSubscription', async function () {
      it('can create a subscription', async function () {
        await expect(registry.connect(subOwner).createSubscription())
          .to.emit(registry, 'SubscriptionCreated')
          .withArgs(1, subOwnerAddress)
        const s = await registry.getSubscription(1)
        expect(s.balance.toString() == '0', 'invalid balance')
        expect(s.owner == subOwnerAddress, 'invalid address')
      })
      it('subscription id increments', async function () {
        await expect(registry.connect(subOwner).createSubscription())
          .to.emit(registry, 'SubscriptionCreated')
          .withArgs(1, subOwnerAddress)
        await expect(registry.connect(subOwner).createSubscription())
          .to.emit(registry, 'SubscriptionCreated')
          .withArgs(2, subOwnerAddress)
      })
      it('cannot create more than the max', async function () {
        const subId = createSubscription(subOwner, [])
        for (let i = 0; i < 100; i++) {
          await registry
            .connect(subOwner)
            .addConsumer(subId, randomAddressString())
        }
        await expect(
          registry.connect(subOwner).addConsumer(subId, randomAddressString()),
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
          registry
            .connect(roles.stranger)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('owner can request transfer', async function () {
        await expect(
          registry
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        )
          .to.emit(registry, 'SubscriptionOwnerTransferRequested')
          .withArgs(subId, subOwnerAddress, strangerAddress)
        // Same request is a noop
        await expect(
          registry
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        ).to.not.emit(registry, 'SubscriptionOwnerTransferRequested')
      })
    })

    describe('#acceptSubscriptionOwnerTransfer', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('subscription must exist', async function () {
        await expect(
          registry
            .connect(subOwner)
            .acceptSubscriptionOwnerTransfer(1203123123),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be requested owner to accept', async function () {
        await expect(
          registry
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        )
        await expect(
          registry.connect(subOwner).acceptSubscriptionOwnerTransfer(subId),
        ).to.be.revertedWith(`MustBeRequestedOwner("${strangerAddress}")`)
      })
      it('requested owner can accept', async function () {
        await expect(
          registry
            .connect(subOwner)
            .requestSubscriptionOwnerTransfer(subId, strangerAddress),
        )
          .to.emit(registry, 'SubscriptionOwnerTransferRequested')
          .withArgs(subId, subOwnerAddress, strangerAddress)
        await expect(
          registry.connect(stranger).acceptSubscriptionOwnerTransfer(subId),
        )
          .to.emit(registry, 'SubscriptionOwnerTransferred')
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
          registry.connect(subOwner).addConsumer(1203123123, strangerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          registry.connect(stranger).addConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('add is idempotent', async function () {
        await registry.connect(subOwner).addConsumer(subId, strangerAddress)
        await registry.connect(subOwner).addConsumer(subId, strangerAddress)
      })
      it('cannot add more than maximum', async function () {
        // There is one consumer, add another 99 to hit the max
        for (let i = 0; i < 99; i++) {
          await registry
            .connect(subOwner)
            .addConsumer(subId, randomAddressString())
        }
        // Adding one more should fail
        // await registry.connect(subOwner).addConsumer(subId, strangerAddress);
        await expect(
          registry.connect(subOwner).addConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`TooManyConsumers()`)
        // Same is true if we first create with the maximum
        const consumers: string[] = []
        for (let i = 0; i < 100; i++) {
          consumers.push(randomAddressString())
        }
        subId = await createSubscription(subOwner, consumers)
        await expect(
          registry.connect(subOwner).addConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`TooManyConsumers()`)
      })
      it('owner can update', async function () {
        await expect(
          registry.connect(subOwner).addConsumer(subId, strangerAddress),
        )
          .to.emit(registry, 'SubscriptionConsumerAdded')
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
          registry
            .connect(subOwner)
            .removeConsumer(1203123123, strangerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          registry.connect(stranger).removeConsumer(subId, strangerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('owner can update', async function () {
        const subBefore = await registry.getSubscription(subId)
        await registry.connect(subOwner).addConsumer(subId, strangerAddress)
        await expect(
          registry.connect(subOwner).removeConsumer(subId, strangerAddress),
        )
          .to.emit(registry, 'SubscriptionConsumerRemoved')
          .withArgs(subId, strangerAddress)
        const subAfter = await registry.getSubscription(subId)
        // Subscription should NOT contain the removed consumer
        expect(subBefore.consumers).to.deep.equal(subAfter.consumers)
      })
      it('can remove all consumers', async function () {
        // Testing the handling of zero.
        await registry.connect(subOwner).addConsumer(subId, strangerAddress)
        await registry.connect(subOwner).removeConsumer(subId, strangerAddress)
        await registry.connect(subOwner).removeConsumer(subId, consumerAddress)
        // Should be empty
        const subAfter = await registry.getSubscription(subId)
        expect(subAfter.consumers).to.deep.equal([])
      })
    })

    describe('#pendingRequestExists', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
        await registry.setConfig(
          config.maxGasLimit,
          config.stalenessSeconds,
          config.gasAfterPaymentCalculation,
          config.weiPerUnitLink,
          config.gasOverhead,
          config.requestTimeoutSeconds,
        )

        await linkToken
          .connect(subOwner)
          .transferAndCall(
            registry.address,
            BigNumber.from('130790416713017745'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await registry.connect(subOwner).addConsumer(subId, client.address)
        await registry.connect(roles.defaultAccount).reg
        await registry.setAuthorizedSenders([oracle.address])
      })
      it('returns false when there is no latest pending request', async function () {
        expect(await registry.connect(subOwner).pendingRequestExists(subId)).to
          .be.false
      })
      it('returns true when the latest request is pending', async function () {
        await client
          .connect(consumer)
          .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId)
        expect(await registry.connect(subOwner).pendingRequestExists(subId)).to
          .be.true
      })
    })

    describe('#cancelSubscription', async function () {
      let subId: number
      beforeEach(async () => {
        subId = await createSubscription(subOwner, [consumerAddress])
      })
      it('subscription must exist', async function () {
        await expect(
          registry
            .connect(subOwner)
            .cancelSubscription(1203123123, subOwnerAddress),
        ).to.be.revertedWith(`InvalidSubscription`)
      })
      it('must be owner', async function () {
        await expect(
          registry.connect(stranger).cancelSubscription(subId, subOwnerAddress),
        ).to.be.revertedWith(`MustBeSubOwner("${subOwnerAddress}")`)
      })
      it('can cancel', async function () {
        await linkToken
          .connect(subOwner)
          .transferAndCall(
            registry.address,
            BigNumber.from('1000'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await expect(
          registry.connect(subOwner).cancelSubscription(subId, strangerAddress),
        )
          .to.emit(registry, 'SubscriptionCanceled')
          .withArgs(subId, strangerAddress, BigNumber.from('1000'))
        const strangerBalance = await linkToken.balanceOf(strangerAddress)
        expect(strangerBalance.toString()).to.equal('1000000000000001000')
        await expect(
          registry.connect(subOwner).getSubscription(subId),
        ).to.be.revertedWith('InvalidSubscription')
      })
      it('can add same consumer after canceling', async function () {
        await linkToken
          .connect(subOwner)
          .transferAndCall(
            registry.address,
            BigNumber.from('1000'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await registry.connect(subOwner).addConsumer(subId, strangerAddress)
        await registry
          .connect(subOwner)
          .cancelSubscription(subId, strangerAddress)
        subId = await createSubscription(subOwner, [consumerAddress])
        // The cancel should have removed this consumer, so we can add it again.
        await registry.connect(subOwner).addConsumer(subId, strangerAddress)
      })
      it('cannot cancel with pending request', async function () {
        await registry.setConfig(
          config.maxGasLimit,
          config.stalenessSeconds,
          config.gasAfterPaymentCalculation,
          config.weiPerUnitLink,
          config.gasOverhead,
          config.requestTimeoutSeconds,
        )

        await linkToken
          .connect(subOwner)
          .transferAndCall(
            registry.address,
            BigNumber.from('130790416713017745'),
            ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
          )
        await registry.connect(subOwner).addConsumer(subId, client.address)
        await registry.connect(roles.defaultAccount).reg
        await registry.setAuthorizedSenders([oracle.address])
        await client
          .connect(consumer)
          .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId)
        // Should revert with outstanding requests
        await expect(
          registry.connect(subOwner).cancelSubscription(subId, strangerAddress),
        ).to.be.revertedWith('PendingRequestExists()')
        // However the owner is able to cancel
        // funds go to the sub owner.
        await expect(
          registry.connect(roles.defaultAccount).ownerCancelSubscription(subId),
        )
          .to.emit(registry, 'SubscriptionCanceled')
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
                .transferAndCall(registry.address, BigNumber.from('1000'), s)
            },
            BigNumber.from('1000'),
          ],
          [
            async function () {
              await registry
                .connect(subOwner)
                .cancelSubscription(subId, strangerAddress)
            },
            BigNumber.from('-1000'),
          ],
        ]
        for (const [fn, expectedBalanceChange] of balanceChangingFns) {
          const startingBalance = await registry.getTotalBalance()
          await fn()
          const endingBalance = await registry.getTotalBalance()
          expect(
            endingBalance.sub(startingBalance).toString() ==
              expectedBalanceChange.toString(),
          )
        }
      })
      it('only owner can recover', async function () {
        await expect(
          registry.connect(subOwner).recoverFunds(strangerAddress),
        ).to.be.revertedWith('OnlyCallableByOwner()')
      })

      it('owner can recover link transferred', async function () {
        // Set the internal balance
        expect(BigNumber.from('0'), linkToken.balanceOf(strangerAddress))
        const s = ethers.utils.defaultAbiCoder.encode(['uint64'], [subId])
        await linkToken
          .connect(subOwner)
          .transferAndCall(registry.address, BigNumber.from('1000'), s)
        // Circumvent internal balance
        await linkToken
          .connect(subOwner)
          .transfer(registry.address, BigNumber.from('1000'))
        // Should recover this 1000
        await expect(
          registry.connect(roles.defaultAccount).recoverFunds(strangerAddress),
        )
          .to.emit(registry, 'FundsRecovered')
          .withArgs(strangerAddress, BigNumber.from('1000'))
        expect(BigNumber.from('1000'), linkToken.balanceOf(strangerAddress))
      })
    })
  })

  describe('#startBilling', () => {
    let subId: number

    beforeEach(async () => {
      await registry.setAuthorizedSenders([oracle.address])

      await registry.setConfig(
        config.maxGasLimit,
        config.stalenessSeconds,
        config.gasAfterPaymentCalculation,
        config.weiPerUnitLink,
        config.gasOverhead,
        config.requestTimeoutSeconds,
      )

      subId = await createSubscription(subOwner, [consumerAddress])

      await linkToken
        .connect(subOwner)
        .transferAndCall(
          registry.address,
          BigNumber.from('54666805176129187'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        )
      await registry.connect(subOwner).addConsumer(subId, client.address)
      await registry.connect(roles.defaultAccount).reg
    })

    it('only callable by registered DONs', async () => {
      await expect(
        registry.connect(consumer).startBilling(stringToHex('some data'), {
          requester: consumerAddress,
          client: consumerAddress,
          subscriptionId: subId,
          gasPrice: 20_000,
          gasLimit: 20_000,
          confirmations: 50,
        }),
      ).to.be.revertedWith(`reverted with custom error 'UnauthorizedSender()'`)
    })

    it('a subscription can only be used by a subscription consumer', async () => {
      await expect(
        oracle
          .connect(stranger)
          .sendRequest(subId, stringToBytes('some data'), 0),
      ).to.be.revertedWith(
        `reverted with custom error 'InvalidConsumer(${subId}, "${strangerAddress}")`,
      )
      await expect(
        client
          .connect(consumer)
          .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId),
      ).to.not.be.reverted
    })

    it('fails if the subscription does not have the funds for the estimated cost', async () => {
      const subId = await createSubscription(subOwner, [subOwnerAddress])
      await registry.connect(subOwner).addConsumer(subId, client.address)

      await expect(
        client
          .connect(subOwner)
          .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId),
      ).to.be.revertedWith(`InsufficientBalance()`)
    })

    it('when successful, emits an event', async () => {
      await expect(
        client
          .connect(consumer)
          .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId),
      ).to.emit(registry, 'BillingStart')
    })

    it('fails multiple requests if the subscription does not have the funds for the estimated cost', async () => {
      client
        .connect(consumer)
        .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId, {
          gasPrice: 1000000008,
        })

      await expect(
        client
          .connect(subOwner)
          .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId, {
            gasPrice: 1000000008,
          }),
      ).to.be.revertedWith(`InsufficientBalance()`)
    })
  })

  describe('#fulfillAndBill', () => {
    let subId: number
    let requestId: string

    beforeEach(async () => {
      subId = await createSubscription(subOwner, [consumerAddress])

      await registry.setConfig(
        config.maxGasLimit,
        config.stalenessSeconds,
        config.gasAfterPaymentCalculation,
        config.weiPerUnitLink,
        config.gasOverhead,
        config.requestTimeoutSeconds,
      )

      await linkToken
        .connect(subOwner)
        .transferAndCall(
          registry.address,
          BigNumber.from('1000000000000000000'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        )
      await registry.connect(subOwner).addConsumer(subId, client.address)
      await registry.connect(roles.defaultAccount).reg
      await registry.setAuthorizedSenders([oracle.address])

      const request = await client
        .connect(consumer)
        .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId)
      requestId = (await request.wait()).events[3].args[0]
    })

    it('only callable by registered DONs', async () => {
      const someAddress = randomAddressString()
      const someSigners = Array(31).fill(ethers.constants.AddressZero)
      someSigners[0] = someAddress
      await expect(
        registry
          .connect(consumer)
          .fulfillAndBill(
            ethers.utils.hexZeroPad(requestId, 32),
            stringToHex('some data'),
            stringToHex('some data'),
            someAddress,
            someSigners,
            1,
            10,
            0,
          ),
      ).to.be.revertedWith('UnauthorizedSender()')
    })

    it('when successful, emits an event', async () => {
      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
        stringToHex(''),
      )
      await expect(
        oracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 500_000 }),
      ).to.emit(registry, 'BillingEnd')
    })

    it('pays the transmitter the expected amount', async () => {
      const oracleBalanceBefore = await linkToken.balanceOf(
        await roles.oracleNode.getAddress(),
      )
      const [subscriptionBalanceBefore] = await registry.getSubscription(subId)

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
        stringToHex(''),
      )

      await expect(
        oracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 500_000 }),
      )
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
        .to.emit(registry, 'BillingEnd')
        .to.emit(client, 'FulfillRequestInvoked')

      await registry
        .connect(roles.oracleNode)
        .oracleWithdraw(
          await roles.oracleNode.getAddress(),
          BigNumber.from('0'),
        )

      const oracleBalanceAfter = await linkToken.balanceOf(
        await roles.oracleNode.getAddress(),
      )
      const [subscriptionBalanceAfter] = await registry.getSubscription(subId)

      expect(subscriptionBalanceBefore.gt(subscriptionBalanceAfter)).to.be.true
      expect(oracleBalanceAfter.gt(oracleBalanceBefore)).to.be.true
      expect(subscriptionBalanceBefore.sub(subscriptionBalanceAfter)).to.equal(
        oracleBalanceAfter.sub(oracleBalanceBefore),
      )
    })
  })

  describe('#oracleWithdraw', async function () {
    it('cannot withdraw with no balance', async function () {
      await expect(
        registry
          .connect(roles.oracleNode)
          .oracleWithdraw(randomAddressString(), BigNumber.from('100')),
      ).to.be.revertedWith(`InsufficientBalance`)
    })
  })
})
