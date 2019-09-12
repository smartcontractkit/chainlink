import { deploy, stranger, toHex, toUtf8 } from './support/helpers'
import { assertBigNum } from './support/matchers'
const GetterSetter = artifacts.require('GetterSetter.sol')

contract('GetterSetter', () => {
  const requestId =
    '0x3bd198932d9cc01e2950ffc518fd38a303812200000000000000000000000000'
  const bytes32 = toHex('Hi Mom!')
  const uint256 = 645746535432
  let gs

  beforeEach(async () => {
    gs = await GetterSetter.new()
  })

  describe('#setBytes32Val', () => {
    it('updates the bytes32 value', async () => {
      await gs.setBytes32(bytes32, { from: stranger })

      const currentBytes32 = await gs.getBytes32.call()
      assert.equal(toUtf8(currentBytes32), toUtf8(bytes32))
    })

    it('logs an event', async () => {
      const tx = await gs.setBytes32(bytes32, { from: stranger })

      assert.equal(1, tx.logs.length)
      assert.equal(stranger.toLowerCase(), tx.logs[0].args.from.toLowerCase())
      assert.equal(toUtf8(bytes32), toUtf8(tx.logs[0].args.value))
    })
  })

  describe('#requestedBytes32', () => {
    it('updates the request ID and value', async () => {
      await gs.requestedBytes32(requestId, bytes32, { from: stranger })

      const currentRequestId = await gs.requestId.call()
      assert.equal(currentRequestId, requestId)

      const currentBytes32 = await gs.getBytes32.call()
      assert.equal(toUtf8(currentBytes32), toUtf8(bytes32))
    })
  })

  describe('#setUint256', () => {
    it('updates uint256 value', async () => {
      await gs.setUint256(uint256, { from: stranger })

      const currentUint256 = await gs.getUint256.call()
      assert.equal(currentUint256, uint256)
    })

    it('logs an event', async () => {
      const tx = await gs.setUint256(uint256, { from: stranger })

      assert.equal(1, tx.logs.length)
      assert.equal(stranger.toLowerCase(), tx.logs[0].args.from.toLowerCase())
      assertBigNum(uint256, tx.logs[0].args.value)
    })
  })

  describe('#requestedUint256', () => {
    it('updates the request ID and value', async () => {
      await gs.requestedUint256(requestId, uint256, { from: stranger })

      const currentRequestId = await gs.requestId.call()
      assert.equal(currentRequestId, requestId)

      const currentUint256 = await gs.getUint256.call()
      assert.equal(currentUint256, uint256)
    })
  })
})
