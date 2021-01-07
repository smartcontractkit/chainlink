import { contract, setup } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { GetterSetter__factory } from '../../ethers/v0.4/factories/GetterSetter__factory'

const GetterSetterContract = new GetterSetter__factory()
const provider = setup.provider()

let roles: setup.Roles

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('GetterSetter', () => {
  const requestId =
    '0x3bd198932d9cc01e2950ffc518fd38a303812200000000000000000000000000'
  const bytes32 = ethers.utils.formatBytes32String('Hi Mom!')
  const uint256 = ethers.utils.bigNumberify(645746535432)
  let gs: contract.Instance<GetterSetter__factory>
  const deployment = setup.snapshot(provider, async () => {
    gs = await GetterSetterContract.connect(roles.defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#setBytes32Val', () => {
    it('updates the bytes32 value', async () => {
      await gs.connect(roles.stranger).setBytes32(bytes32)

      const currentBytes32 = await gs.getBytes32()
      assert.equal(
        ethers.utils.toUtf8String(currentBytes32),
        ethers.utils.toUtf8String(bytes32),
      )
    })

    it('logs an event', async () => {
      const tx = await gs.connect(roles.stranger).setBytes32(bytes32)

      const receipt = await tx.wait()
      const args: any = receipt.events?.[0].args

      assert.equal(1, receipt.events?.length)
      assert.equal(
        roles.stranger.address.toLowerCase(),
        args.from.toLowerCase(),
      )
      assert.equal(
        ethers.utils.toUtf8String(bytes32),
        ethers.utils.toUtf8String(args.value),
      )
    })
  })

  describe('#requestedBytes32', () => {
    it('updates the request ID and value', async () => {
      await gs.connect(roles.stranger).requestedBytes32(requestId, bytes32)

      const currentRequestId = await gs.requestId()
      assert.equal(currentRequestId, requestId)

      const currentBytes32 = await gs.getBytes32()
      assert.equal(
        ethers.utils.toUtf8String(currentBytes32),
        ethers.utils.toUtf8String(bytes32),
      )
    })
  })

  describe('#setUint256', () => {
    it('updates uint256 value', async () => {
      await gs.connect(roles.stranger).setUint256(uint256)

      const currentUint256 = await gs.getUint256()
      assert.isTrue(currentUint256.eq(uint256))
    })

    it('logs an event', async () => {
      const tx = await gs.connect(roles.stranger).setUint256(uint256)
      const receipt = await tx.wait()
      const args: any = receipt.events?.[0].args

      assert.equal(1, receipt.events?.length)
      assert.equal(
        roles.stranger.address.toLowerCase(),
        args.from.toLowerCase(),
      )

      assert.isTrue(uint256.eq(args.value))
    })
  })

  describe('#requestedUint256', () => {
    it('updates the request ID and value', async () => {
      await gs.connect(roles.stranger).requestedUint256(requestId, uint256)

      const currentRequestId = await gs.requestId()
      assert.equal(currentRequestId, requestId)

      const currentUint256 = await gs.getUint256()
      assert.isTrue(currentUint256.eq(uint256))
    })
  })
})
