import * as h from '../src/helpersV2'
import ganache from 'ganache-core'
import { AbstractContract } from '../src/contract'
import { ethers } from 'ethers'
import { assert } from 'chai'
import { GetterSetter } from 'contracts/GetterSetter'

const GetterSetterContract = AbstractContract.fromArtifactName('GetterSetter')
const ganacheProvider: any = ganache.provider()

let roles: h.Roles

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(ganacheProvider)

  roles = rolesAndPersonas.roles
})

describe('GetterSetter', () => {
  const requestId =
    '0x3bd198932d9cc01e2950ffc518fd38a303812200000000000000000000000000'
  const bytes32 = ethers.utils.formatBytes32String('Hi Mom!')
  const uint256 = ethers.utils.bigNumberify(645746535432)
  let gs: GetterSetter

  beforeEach(async () => {
    gs = await GetterSetterContract.deploy(roles.defaultAccount)
  })

  describe('#setBytes32Val', () => {
    it('updates the bytes32 value', async () => {
      await gs.connect(roles.stranger).functions.setBytes32(bytes32)

      const currentBytes32 = await gs.functions.getBytes32()
      assert.equal(
        ethers.utils.toUtf8String(currentBytes32),
        ethers.utils.toUtf8String(bytes32),
      )
    })

    it('logs an event', async () => {
      const tx = await gs.connect(roles.stranger).functions.setBytes32(bytes32)

      const receipt = await tx.wait()
      const args: any = receipt.events![0].args

      assert.equal(1, receipt.events!.length)
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
      await gs
        .connect(roles.stranger)
        .functions.requestedBytes32(requestId, bytes32)

      const currentRequestId = await gs.functions.requestId()
      assert.equal(currentRequestId, requestId)

      const currentBytes32 = await gs.functions.getBytes32()
      assert.equal(
        ethers.utils.toUtf8String(currentBytes32),
        ethers.utils.toUtf8String(bytes32),
      )
    })
  })

  describe('#setUint256', () => {
    it('updates uint256 value', async () => {
      await gs.connect(roles.stranger).functions.setUint256(uint256)

      const currentUint256 = await gs.functions.getUint256()
      assert.isTrue(currentUint256.eq(uint256))
    })

    it('logs an event', async () => {
      const tx = await gs.connect(roles.stranger).functions.setUint256(uint256)
      const receipt = await tx.wait()
      const args: any = receipt.events![0].args

      assert.equal(1, receipt.events!.length)
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

      const currentRequestId = await gs.functions.requestId()
      assert.equal(currentRequestId, requestId)

      const currentUint256 = await gs.functions.getUint256()
      assert.isTrue(currentUint256.eq(uint256))
    })
  })
})
