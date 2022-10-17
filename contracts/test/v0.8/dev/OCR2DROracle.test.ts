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

before(async () => {
  roles = (await getUsers()).roles

  ocr2drOracleFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ocr2dr/OCR2DROracle.sol:OCR2DROracle',
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
    oracle = await ocr2drOracleFactory
      .connect(roles.defaultAccount)
      .deploy(await roles.defaultAccount.getAddress(), donPublicKey)
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
      expect(await oracle.callStatic.getDONPublicKey()).to.be.equal(
        donPublicKey,
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

    const setAuthorizedSender = async (sender: string) => {
      await expect(oracle.setAuthorizedSenders([sender])).not.to.be.reverted
    }

    it('#fulfillRequest reverts for unknown requestId', async () => {
      const requestId =
        '0x67c6a2e151d4352a55021b5d0028c18121cfc24c7d73b179d22b17daff069c6e'

      await expect(
        oracle.fulfillRequest(
          requestId,
          stringToHex('response'),
          stringToHex(''),
        ),
      ).to.be.revertedWith('InvalidRequestID')
    })

    it('#fulfillRequest reverts for unauthorized sender', async () => {
      const requestId = await placeTestRequest()

      await expect(
        oracle.fulfillRequest(
          requestId,
          stringToHex('response'),
          stringToHex(''),
        ),
      ).to.be.revertedWith('UnauthorizedSender')
    })

    it('#fulfillRequest reverts on low consumer gas', async () => {
      const { roles } = await getUsers()
      const sender = await roles.oracleNode.getAddress()
      const requestId = await placeTestRequest()

      await setAuthorizedSender(sender)

      await expect(
        oracle
          .connect(roles.oracleNode)
          .fulfillRequest(requestId, stringToHex('response'), stringToHex(''), {
            gasLimit: 300000,
          }),
      ).to.be.revertedWith('LowGasForConsumer')
    })

    it('#fulfillRequest emits OracleResponse', async () => {
      const { roles } = await getUsers()
      const sender = await roles.oracleNode.getAddress()
      const requestId = await placeTestRequest()

      await setAuthorizedSender(sender)

      await expect(
        oracle
          .connect(roles.oracleNode)
          .fulfillRequest(requestId, stringToHex('response'), stringToHex('')),
      )
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
    })

    it('#fulfillRequest invokes client fulfillRequest', async () => {
      const { roles } = await getUsers()
      const sender = await roles.oracleNode.getAddress()
      const requestId = await placeTestRequest()

      await setAuthorizedSender(sender)

      await expect(
        oracle
          .connect(roles.oracleNode)
          .fulfillRequest(
            requestId,
            stringToHex('response'),
            stringToHex('err'),
          ),
      )
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, stringToHex('response'), stringToHex('err'))
    })

    it('#fulfillRequest invalidates requestId', async () => {
      const { roles } = await getUsers()
      const sender = await roles.oracleNode.getAddress()
      const requestId = await placeTestRequest()

      await setAuthorizedSender(sender)

      await expect(
        oracle
          .connect(roles.oracleNode)
          .fulfillRequest(
            requestId,
            stringToHex('response'),
            stringToHex('err'),
          ),
      )
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, stringToHex('response'), stringToHex('err'))

      // for second fulfill the requestId becomes invalid
      await expect(
        oracle
          .connect(roles.oracleNode)
          .fulfillRequest(
            requestId,
            stringToHex('response'),
            stringToHex('err'),
          ),
      ).to.be.revertedWith('InvalidRequestID')
    })
  })
})
