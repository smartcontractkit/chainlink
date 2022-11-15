import { ethers } from 'hardhat'
import { expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let ocr2drOracleFactory: ContractFactory
let clientTestHelperFactory: ContractFactory
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

  ocr2drOracleFactory = await ethers.getContractFactory(
    'src/v0.8/tests/OCR2DROracleHelper.sol:OCR2DROracleHelper',
    roles.defaultAccount,
  )

  clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/tests/OCR2DRClientTestHelper.sol:OCR2DRClientTestHelper',
    roles.consumer,
  )
})

describe('OCR2DROracle', () => {
  const donPublicKey =
    '0x3804a19f2437f7bba4fcfbc194379e43e514aa98073db3528ccdbdb642e24011'
  let oracle: Contract
  let client: Contract

  beforeEach(async () => {
    const { roles } = await getUsers()
    oracle = await ocr2drOracleFactory.connect(roles.defaultAccount).deploy()
    client = await clientTestHelperFactory
      .connect(roles.consumer)
      .deploy(oracle.address)
  })

  describe('General', () => {
    it('#typeAndVersion', async () => {
      expect(await oracle.callStatic.typeAndVersion()).to.be.equal(
        'OCR2DROracle 0.0.0',
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
  })

  describe('Sending requests', () => {
    it('#sendRequest emits OracleRequest event', async () => {
      const data = stringToHex('some data')
      await expect(oracle.sendRequest(0, data))
        .to.emit(oracle, 'OracleRequest')
        .withArgs(anyValue, data)
    })

    it('#sendRequest reverts for empty data', async () => {
      const data = stringToHex('')
      await expect(oracle.sendRequest(0, data)).to.be.revertedWith(
        'EmptyRequestData',
      )
    })

    it('#sendRequest returns non-empty requestId', async () => {
      const data = stringToHex('test data')
      const requestId = await oracle.callStatic.sendRequest(0, data)
      expect(requestId).not.to.be.empty
    })

    it('#sendRequest returns different requestIds', async () => {
      const data = stringToHex('test data')
      const requestId1 = await oracle.callStatic.sendRequest(0, data)
      await expect(oracle.sendRequest(0, data))
        .to.emit(oracle, 'OracleRequest')
        .withArgs(anyValue, data)
      const requestId2 = await oracle.callStatic.sendRequest(0, data)
      expect(requestId1).not.to.be.equal(requestId2)
    })
  })

  describe('Fulfilling requests', () => {
    const placeTestRequest = async () => {
      const requestId = await client.callStatic.sendSimpleRequestWithJavaScript(
        'function(){}',
        0,
      )
      await expect(client.sendSimpleRequestWithJavaScript('function(){}', 0))
        .to.emit(client, 'RequestSent')
        .withArgs(requestId)
      return requestId
    }

    it('#fulfillRequest reverts for unknown requestId', async () => {
      const requestId =
        '0x67c6a2e151d4352a55021b5d0028c18121cfc24c7d73b179d22b17daff069c6e'

      const report = encodeReport(
        requestId,
        stringToHex('response'),
        stringToHex(''),
      )

      await expect(oracle.callReport(report)).to.be.revertedWith(
        'InvalidRequestID',
      )
    })

    it('#fulfillRequest emits OracleResponse', async () => {
      const { roles } = await getUsers()
      const requestId = await placeTestRequest()

      const report = encodeReport(
        requestId,
        stringToHex('response'),
        stringToHex(''),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
    })

    it('#fulfillRequest emits UserCallbackError if callback reverts', async () => {
      const { roles } = await getUsers()
      const requestId = await placeTestRequest()

      const report = encodeReport(
        requestId,
        stringToHex('response'),
        stringToHex(''),
      )

      await client.setRevertFulfillRequest(true)

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'UserCallbackError')
        .withArgs(requestId, 'asked to revert')
    })

    it('#fulfillRequest emits UserCallbackRawError if callback does invalid op', async () => {
      const { roles } = await getUsers()
      const requestId = await placeTestRequest()

      const report = encodeReport(
        requestId,
        stringToHex('response'),
        stringToHex(''),
      )

      await client.setDoInvalidOperation(true)

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(oracle, 'UserCallbackRawError')
        .withArgs(requestId, anyValue)
    })

    it('#fulfillRequest invokes client fulfillRequest', async () => {
      const { roles } = await getUsers()
      const requestId = await placeTestRequest()

      const report = encodeReport(
        requestId,
        stringToHex('response'),
        stringToHex('err'),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, stringToHex('response'), stringToHex('err'))
    })

    it('#fulfillRequest invalidates requestId', async () => {
      const { roles } = await getUsers()
      const requestId = await placeTestRequest()

      const report = encodeReport(
        requestId,
        stringToHex('response'),
        stringToHex('err'),
      )

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, stringToHex('response'), stringToHex('err'))

      // for second fulfill the requestId becomes invalid
      await expect(
        oracle.connect(roles.oracleNode).callReport(report),
      ).to.be.revertedWith('InvalidRequestID')
    })

    it('#_report reverts for inconsistent encoding', async () => {
      const { roles } = await getUsers()
      const requestId = await placeTestRequest()

      const abi = ethers.utils.defaultAbiCoder
      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [[requestId], [], []],
      )

      await expect(
        oracle.connect(roles.oracleNode).callReport(report),
      ).to.be.revertedWith('InconsistentReportData')
    })

    it('#_report handles multiple reports', async () => {
      const { roles } = await getUsers()
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

      await expect(oracle.connect(roles.oracleNode).callReport(report))
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId1, result1, err)
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId2, result2, err)
    })

    it('#_report handles multiple failures', async () => {
      const { roles } = await getUsers()
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
        .withArgs(requestId1, 'asked to revert')
        .to.emit(oracle, 'UserCallbackError')
        .withArgs(requestId2, 'asked to revert')
    })
  })
})
