import { ethers } from 'hardhat'
import { expect } from 'chai'
import { BigNumber, Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let functionsOracleFactory: ContractFactory
let clientTestHelperFactory: ContractFactory
let functionsBillingRegistryFactory: ContractFactory
let linkTokenFactory: ContractFactory
let mockAggregatorV3Factory: ContractFactory
let roles: Roles

const stringToHex = (s: string) => {
  return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(s))
}

const anyValue = () => true

const encodeReport = (requestId: string, result: string, err: string) => {
  const abi = ethers.utils.defaultAbiCoder
  return abi.encode(
    ['bytes32[]', 'bytes[]', 'bytes[]'],
    [[requestId], [result], [err]],
  )
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
    roles.defaultAccount,
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

describe('FunctionsOracle', () => {
  let subscriptionId: number
  const donPublicKey =
    '0x3804a19f2437f7bba4fcfbc194379e43e514aa98073db3528ccdbdb642e24011'
  let client: Contract
  let oracle: Contract
  let registry: Contract
  let linkToken: Contract
  let mockLinkEth: Contract
  let transmitters: string[]

  beforeEach(async () => {
    // Deploy contracts
    linkToken = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    mockLinkEth = await mockAggregatorV3Factory.deploy(
      0,
      ethers.BigNumber.from(5021530000000000),
    )
    oracle = await functionsOracleFactory.connect(roles.defaultAccount).deploy()
    registry = await functionsBillingRegistryFactory
      .connect(roles.defaultAccount)
      .deploy(linkToken.address, mockLinkEth.address, oracle.address)

    // Setup contracts
    await oracle.setRegistry(registry.address)
    await oracle.deactivateAuthorizedReceiver()
    client = await clientTestHelperFactory
      .connect(roles.defaultAccount)
      .deploy(oracle.address)
    await registry.setAuthorizedSenders([oracle.address])

    await registry.setConfig(
      1_000_000, // maxGasLimit
      86_400, // stalenessSeconds
      39_173 /* gasAfterPaymentCalculation
      gathered by taking the difference from gasleft() directly after payment calculation, and then again after the BillingEnd event, using a hardhat console log
      */,
      ethers.BigNumber.from('5000000000000000'), // fallbackWeiPerUnitLink
      519_719 /* gasOverhead
      gathered by taking the difference from initialGas and gasleft() directly after payment calculation, adding back the user's callback gas usage, using a hardhat console log
      NOTE: this number can vary slightly by number of nodes on the DON
      */,
      300, // requestTimeoutSeconds
    )

    // Setup accounts
    const createSubTx = await registry
      .connect(roles.defaultAccount)
      .createSubscription()
    const receipt = await createSubTx.wait()
    subscriptionId = receipt.events[0].args['subscriptionId'].toNumber()

    await registry
      .connect(roles.defaultAccount)
      .addConsumer(subscriptionId, await roles.defaultAccount.getAddress())

    await registry
      .connect(roles.defaultAccount)
      .addConsumer(subscriptionId, client.address)

    await linkToken
      .connect(roles.defaultAccount)
      .transferAndCall(
        registry.address,
        ethers.BigNumber.from('300938394174049741'),
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subscriptionId]),
      )

    const signers = Array.from(
      [0, 0, 0, 0],
      (_) => ethers.Wallet.createRandom().address,
    )
    transmitters = [
      await roles.oracleNode1.getAddress(),
      await roles.oracleNode2.getAddress(),
      await roles.oracleNode3.getAddress(),
      await roles.oracleNode4.getAddress(),
    ]
    await oracle.setConfig(signers, transmitters, 1, [], 1, [])
  })

  describe('General', () => {
    it('#typeAndVersion', async () => {
      expect(await oracle.callStatic.typeAndVersion()).to.be.equal(
        'FunctionsOracle 0.0.0',
      )
    })

    it('returns DON public key set on this Oracle', async () => {
      await expect(oracle.setDONPublicKey(donPublicKey)).not.to.be.reverted
      expect(await oracle.callStatic.getDONPublicKey()).to.be.equal(
        donPublicKey,
      )
    })

    it('reverts setDONPublicKey for empty data', async () => {
      const emptyPublicKey = stringToHex('')
      await expect(oracle.setDONPublicKey(emptyPublicKey)).to.be.revertedWith(
        'EmptyPublicKey',
      )
    })

    async function validatePubKeys(
      expectedNodes: string[],
      expectedKeys: string[],
    ) {
      const allNodesAndKeys = await oracle.getAllNodePublicKeys()
      for (let i = 0; i < expectedNodes.length; i++) {
        expect(allNodesAndKeys[0][i]).to.be.equal(expectedNodes[i])
        expect(allNodesAndKeys[1][i]).to.be.equal(expectedKeys[i])
      }
    }

    it('set/delete/get node public keys', async () => {
      const emptyKey = stringToHex('')
      const publicKey2 = stringToHex('key420')
      const publicKey3 = stringToHex('key666')

      await oracle.setNodePublicKey(roles.oracleNode2.getAddress(), publicKey2)
      await oracle.setNodePublicKey(roles.oracleNode3.getAddress(), publicKey3)
      validatePubKeys(transmitters, [
        emptyKey,
        publicKey2,
        publicKey3,
        emptyKey,
      ])

      await oracle.deleteNodePublicKey(roles.oracleNode1.getAddress())
      await oracle.deleteNodePublicKey(roles.oracleNode2.getAddress())
      validatePubKeys(transmitters, [emptyKey, emptyKey, publicKey3, emptyKey])
    })

    it('reverts setNodePublicKey for unauthorized callers', async () => {
      const pubKey = stringToHex('abcd')

      await expect(
        oracle
          .connect(roles.defaultAccount)
          .setNodePublicKey(roles.oracleNode2.getAddress(), pubKey),
      ).not.to.be.reverted

      await expect(
        oracle
          .connect(roles.consumer)
          .setNodePublicKey(roles.oracleNode2.getAddress(), pubKey),
      ).to.be.revertedWith('UnauthorizedPublicKeyChange')

      await expect(
        oracle
          .connect(roles.consumer)
          .setNodePublicKey(roles.consumer.getAddress(), pubKey),
      ).to.be.revertedWith('UnauthorizedPublicKeyChange')

      await expect(
        oracle
          .connect(roles.oracleNode2)
          .setNodePublicKey(roles.oracleNode3.getAddress(), pubKey),
      ).to.be.revertedWith('UnauthorizedPublicKeyChange')

      await expect(
        oracle
          .connect(roles.oracleNode2)
          .setNodePublicKey(roles.oracleNode2.getAddress(), pubKey),
      ).not.to.be.reverted
    })

    it('reverts deleteNodePublicKey for unauthorized callers', async () => {
      await expect(
        oracle
          .connect(roles.defaultAccount)
          .deleteNodePublicKey(roles.oracleNode2.getAddress()),
      ).not.to.be.reverted

      await expect(
        oracle
          .connect(roles.consumer)
          .deleteNodePublicKey(roles.oracleNode2.getAddress()),
      ).to.be.revertedWith('UnauthorizedPublicKeyChange')

      await expect(
        oracle
          .connect(roles.consumer)
          .deleteNodePublicKey(roles.consumer.getAddress()),
      ).not.to.be.reverted
    })
  })

  describe('Sending requests', () => {
    it('#sendRequest emits OracleRequest event', async () => {
      const data = stringToHex('some data')
      await expect(oracle.sendRequest(subscriptionId, data, 0))
        .to.emit(oracle, 'OracleRequest')
        .withArgs(
          anyValue,
          await roles.defaultAccount.getAddress(),
          await roles.defaultAccount.getAddress(),
          subscriptionId,
          await roles.defaultAccount.getAddress(),
          data,
        )
    })

    it('#sendRequest reverts for empty data', async () => {
      const data = stringToHex('')
      await expect(
        oracle.sendRequest(subscriptionId, data, 0),
      ).to.be.revertedWith('EmptyRequestData')
    })

    it('#sendRequest returns non-empty requestId', async () => {
      const data = stringToHex('test data')
      const requestId = await oracle.callStatic.sendRequest(
        subscriptionId,
        data,
        0,
      )
      expect(requestId).not.to.be.empty
    })

    it('#sendRequest returns different requestIds', async () => {
      const data = stringToHex('test data')
      const requestId1 = await oracle.callStatic.sendRequest(
        subscriptionId,
        data,
        0,
      )
      await expect(oracle.sendRequest(subscriptionId, data, 0))
        .to.emit(oracle, 'OracleRequest')
        .withArgs(
          anyValue,
          await roles.defaultAccount.getAddress(),
          await roles.defaultAccount.getAddress(),
          subscriptionId,
          await roles.defaultAccount.getAddress(),
          data,
        )
      const requestId2 = await oracle.callStatic.sendRequest(
        subscriptionId,
        data,
        0,
      )
      expect(requestId1).not.to.be.equal(requestId2)
    })
  })

  describe('Fulfilling requests', () => {
    const placeTestRequest = async () => {
      const requestId = await client
        .connect(roles.oracleNode)
        .callStatic.sendSimpleRequestWithJavaScript(
          'function(){}',
          subscriptionId,
        )
      await expect(
        client
          .connect(roles.oracleNode)
          .sendSimpleRequestWithJavaScript('function(){}', subscriptionId),
      )
        .to.emit(client, 'RequestSent')
        .withArgs(requestId)
      return requestId
    }

    it('#fulfillRequest emits an error for unknown requestId', async () => {
      const requestId =
        '0x67c6a2e151d4352a55021b5d0028c18121cfc24c7d73b179d22b17daff069c6e'

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex(''),
      )

      await expect(oracle.callReport(report)).to.emit(
        oracle,
        'UserCallbackRawError',
      )
    })

    it('#fulfillRequest emits OracleResponse', async () => {
      const requestId = await placeTestRequest()

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex(''),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
    })

    it('#estimateCost correctly estimates cost [ @skip-coverage ]', async () => {
      const [subscriptionBalanceBefore] = await registry.getSubscription(
        subscriptionId,
      )

      const request = await client
        .connect(roles.oracleNode)
        .sendSimpleRequestWithJavaScript('function(){}', subscriptionId)
      const receipt = await request.wait()
      const requestId = receipt.events[3].args[0]

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex(''),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
        .to.emit(registry, 'BillingEnd')

      const [subscriptionBalanceAfter] = await registry.getSubscription(
        subscriptionId,
      )

      const feeData = await ethers.provider.getFeeData()
      const estimatedCost = await client.estimateJuelCost(
        'function(){}',
        subscriptionId,
        feeData.gasPrice ?? BigNumber.from(0),
      )
      // Expect charged amount to be +-0.01%
      expect(
        subscriptionBalanceBefore.sub(subscriptionBalanceAfter),
      ).to.be.below(estimatedCost.add(estimatedCost.div(100)))
      expect(
        subscriptionBalanceBefore.sub(subscriptionBalanceAfter),
      ).to.be.above(estimatedCost.sub(estimatedCost.div(100)))
    })

    it('#fulfillRequest emits UserCallbackError if callback reverts', async () => {
      const requestId = await placeTestRequest()

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex(''),
      )

      await client.setRevertFulfillRequest(true)

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'UserCallbackError')
        .withArgs(requestId, anyValue)
    })

    it('#fulfillRequest emits UserCallbackError if callback does invalid op', async () => {
      const requestId = await placeTestRequest()

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex(''),
      )

      await client.setDoInvalidOperation(true)

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'UserCallbackError')
        .withArgs(requestId, anyValue)
    })

    it('#fulfillRequest invokes client fulfillRequest', async () => {
      const requestId = await placeTestRequest()

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex('err'),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, stringToHex('response'), stringToHex('err'))
    })

    it('#fulfillRequest invalidates requestId', async () => {
      const requestId = await placeTestRequest()

      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('response'),
        stringToHex('err'),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, stringToHex('response'), stringToHex('err'))

      // for second fulfill the requestId becomes invalid
      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'UserCallbackRawError')
        .withArgs(requestId, '0xda7aa3e1')
    })

    it('#_report reverts for inconsistent encoding', async () => {
      const requestId = await placeTestRequest()

      const abi = ethers.utils.defaultAbiCoder
      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [[requestId], [], []],
      )

      await expect(
        oracle.connect(roles.oracleNode).callReport(report),
      ).to.be.revertedWith('ReportInvalid()')
    })

    it('#_report handles multiple reports', async () => {
      const requestId1 = await placeTestRequest()
      const requestId2 = await placeTestRequest()
      const result1 = stringToHex('result1')
      const result2 = stringToHex('result2')
      const err = stringToHex('')

      const abi = ethers.utils.defaultAbiCoder
      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [
          [requestId1, requestId2],
          [result1, result2],
          [err, err],
        ],
      )

      await expect(
        oracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 300_000 }),
      )
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId1, result1, err)
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId2, result2, err)
    })

    it('#_report handles multiple failures', async () => {
      const requestId1 = await placeTestRequest()
      const requestId2 = await placeTestRequest()
      const result1 = stringToHex('result1')
      const result2 = stringToHex('result2')
      const err = stringToHex('')

      const abi = ethers.utils.defaultAbiCoder
      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [
          [requestId1, requestId2],
          [result1, result2],
          [err, err],
        ],
      )

      await client.setRevertFulfillRequest(true)

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'UserCallbackError')
        .withArgs(requestId1, anyValue)
        .to.emit(oracle, 'UserCallbackError')
        .withArgs(requestId2, anyValue)
    })
  })
})
