import { ethers } from 'hardhat'
import { expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let ocr2drOracleFactory: ContractFactory
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
})

describe('OCR2DROracle', () => {
  const donPublicKey =
    '0x3804a19f2437f7bba4fcfbc194379e43e514aa98073db3528ccdbdb642e24011'
  let oracle: Contract

  beforeEach(async () => {
    const { roles } = await getUsers()
    oracle = await ocr2drOracleFactory.connect(roles.defaultAccount).deploy()
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
})
