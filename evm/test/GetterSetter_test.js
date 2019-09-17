import * as h from '../src/helpers'
import { assertBigNum } from '../src/matchers'
const GetterSetter = artifacts.require('GetterSetter.sol')

let roles

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas()

  roles = rolesAndPersonas.roles
})

contract('GetterSetter', () => {
  const requestId =
    '0x3bd198932d9cc01e2950ffc518fd38a303812200000000000000000000000000'
  const bytes32 = h.toHex('Hi Mom!')
  const uint256 = 645746535432
  let gs

  beforeEach(async () => {
    gs = await GetterSetter.new()
  })

  describe('#setBytes32Val', () => {
    it('updates the bytes32 value', async () => {
      await gs.setBytes32(bytes32, { from: roles.stranger })

      const currentBytes32 = await gs.getBytes32.call()
      assert.equal(h.toUtf8(currentBytes32), h.toUtf8(bytes32))
    })

    it('logs an event', async () => {
      const tx = await gs.setBytes32(bytes32, { from: roles.stranger })

      assert.equal(1, tx.logs.length)
      assert.equal(
        roles.stranger.toLowerCase(),
        tx.logs[0].args.from.toLowerCase(),
      )
      assert.equal(h.toUtf8(bytes32), h.toUtf8(tx.logs[0].args.value))
    })
  })

  describe('#requestedBytes32', () => {
    it('updates the request ID and value', async () => {
      await gs.requestedBytes32(requestId, bytes32, { from: roles.stranger })

      const currentRequestId = await gs.requestId.call()
      assert.equal(currentRequestId, requestId)

      const currentBytes32 = await gs.getBytes32.call()
      assert.equal(h.toUtf8(currentBytes32), h.toUtf8(bytes32))
    })
  })

  describe('#setUint256', () => {
    it('updates uint256 value', async () => {
      await gs.setUint256(uint256, { from: roles.stranger })

      const currentUint256 = await gs.getUint256.call()
      assert.equal(currentUint256, uint256)
    })

    it('logs an event', async () => {
      const tx = await gs.setUint256(uint256, { from: roles.stranger })

      assert.equal(1, tx.logs.length)
      assert.equal(
        roles.stranger.toLowerCase(),
        tx.logs[0].args.from.toLowerCase(),
      )
      assertBigNum(uint256, tx.logs[0].args.value)
    })
  })

  describe('#requestedUint256', () => {
    it('updates the request ID and value', async () => {
      await gs.requestedUint256(requestId, uint256, { from: roles.stranger })

      const currentRequestId = await gs.requestId.call()
      assert.equal(currentRequestId, requestId)

      const currentUint256 = await gs.getUint256.call()
      assert.equal(currentUint256, uint256)
    })
  })
})
