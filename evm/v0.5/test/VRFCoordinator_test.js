import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const VRFCoordinatorFactory = artifacts.require('dev/VRFCoordinator.sol')
const VRFConsumerFactory = artifacts.require('tests/VRFConsumer.sol')

const generator = [
  // This is the public key for "secret" key 1.
  '0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798',
  '0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8',
].map(h.bigNum)

const key = generator // No need to multiply, since secretKey = 1
const fee = 100
const seed = 1 // Never do this

// Taken from the proofBlob in VRF_test.js
const proofBlob =
  '0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798483ada77' +
  '26a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8530fddd863609aa120' +
  '30a07c5fdb323bb392a88343cea123b7f074883d2654c46fd4ee394bf2a3de542c0e5f3c86' +
  'fc8f75b278a017701a59d69bdf5134dd6b70409cfc2326bdefab042b03b630d2dde19c9591' +
  '19e3b32c9060150ea3f951500cbf6303dcd9421054fbd4fc49cf2d221d1e194bcccb9573ab' +
  '5fbd4fe8d6e4f1360000000000000000000000000000000000000000000000000000000000' +
  '0000010000000000000000000000007e5f4552091a69125d5dfcb7b8c2659029395bdfa2e0' +
  '3a05b089db7b79cd0f6655d6af3e2d06bd0129f87f9f2155612b4e2a41d80a1dadcabf900b' +
  'dfb6484e9a4390bffa6ccd666a565a991f061faf868cc9fce8f82b4f9161ab41ae7c11e7de' +
  'b628024ef9f5e9a0bca029f0ccb5cb534c70be31f26e7c0b4f039ca54cfa100b3457b301ac' +
  'b3e0b6c690d7ea5a86f8e1c481057e39ac5202b29f0aa8a42507af9fd4de3266e804f217e8' +
  '430158254b4a8e96f746'

// Taken from expectedOutput in VRF_test.js
const vRFOutput = h.bigNum(
  '0x7c8bf11f27437d2ced1f68cb3a9125a45a5046e22ab062af2a31fb676cdd161b',
  16,
)

const jobID = '0x1234567890abcdef1234567890abcdef'
const jobIDASCIIHex =
  '0x' + Buffer.from(jobID.slice(2), 'ascii').toString('hex')

contract('VRFCoordinator', () => {
  let VRFCoordinator, VRFConsumer, link, keyHash, requestId
  beforeEach(async () => {
    link = await h.linkContract()
    VRFCoordinator = await VRFCoordinatorFactory.new(link.address)
    VRFConsumer = await VRFConsumerFactory.new(
      VRFCoordinator.address,
      link.address,
    )
    await VRFCoordinator.registerProvingKey(fee, key, jobIDASCIIHex, {
      from: h.personas.Neil,
    })
    keyHash = await VRFCoordinator.hashOfKey(key)
    requestId = await VRFCoordinator.makeRequestId(keyHash, seed)
    assert(
      link.transfer(VRFConsumer.address, 1000),
      'failed to fund consumer contract with LINK',
    )
  })
  it('has a limited public interface', async () => {
    h.checkPublicABI(VRFCoordinator, [
      'callbacks',
      'serviceAgreements',
      'registerProvingKey',
      'fulfillRandomnessRequest',
      'withdraw',
      'withdrawableTokens',
      'onTokenTransfer',
      'makeRequestId',
      'randomValueFromVRFProof',
      'hashOfKey',
    ])
  })
  describe('#registerProvingKey', async () => {
    it('correctly represents the key registered in beforeEach', async () => {
      const keyRepresentation = await VRFCoordinator.serviceAgreements(keyHash)
      assert.equal(
        h.personas.Neil,
        keyRepresentation.vRFOracle,
        "sender should be registered as key's oracle",
      )
      assertBigNum(fee, keyRepresentation.fee, 'fee not recorded')
      const events = await h.getEvents(VRFCoordinator)
      const nSAEvents = events.filter(e => e.event == 'NewServiceAgreement')
      assert.equal(1, nSAEvents.length)
      const event = nSAEvents[0].args
      assert.equal(keyHash, event.keyHash)
      assert.equal(fee.toString(), event.fee)
    })
    it('rejects attempts to re-register a key', async () => {
      h.assertActionThrows(async () => {
        await VRFCoordinator.registerProvingKey(fee, key, jobIDASCIIHex, {
          from: h.personas.Neil,
        })
      })
    })
  })
  describe('#randomnessRequest', async () => {
    beforeEach(async () => {
      await VRFConsumer.requestRandomness(keyHash, fee, seed, {
        from: h.personas.Carol,
      })
    })
    it('emits an event signaling the request, and stores it correctly', async () => {
      const events = await h.getEvents(VRFCoordinator)
      const rREvents = events.filter(e => e.event == 'RandomnessRequest')
      assert.equal(1, rREvents.length)
      const event = rREvents[0].args
      console.log(rREvents[0])
      assert.equal(keyHash, event.keyHash)
      assert.equal(seed, event.seed)
      assert.equal(jobIDASCIIHex, event.jobID)
      assert.equal(VRFConsumer.address, event.sender)
      assert.equal(fee, event.fee)
      const callback = await VRFCoordinator.callbacks(requestId)
      assert(h.personas.Carol, callback.callbackContract)
      assertBigNum(fee, callback.randomnessFee)
      assertBigNum(seed, callback.seed)
    })
    it('rejects attempts to request a previously observed seed', async () => {
      h.assertActionThrows(async () => {
        await VRFConsumer.requestRandomness(keyHash, fee, seed)
      })
    })
  })
  describe('#fulfillRandomnessRequest', async () => {
    beforeEach(async () => {
      await VRFConsumer.requestRandomness(keyHash, fee, seed)
      const resp = await VRFCoordinator.fulfillRandomnessRequest(proofBlob, {
        from: h.personas.Neil,
      })
      assert(resp, 'fulfillment failed')
    })
    it('pays the oracle and reports the randomness', async () => {
      const neilBalance = await VRFCoordinator.withdrawableTokens(
        h.personas.Neil,
      )
      assertBigNum(fee, neilBalance)
      assertBigNum(vRFOutput, await VRFConsumer.randomnessOutput())
      assert.equal(requestId, await VRFConsumer.requestId())
    })
  })
  describe('#withdraw', async () => {
    beforeEach(async () => {
      await VRFConsumer.requestRandomness(keyHash, fee, seed)
      const resp = await VRFCoordinator.fulfillRandomnessRequest(proofBlob, {
        from: h.personas.Neil,
      })
      assert(resp, 'fulfillment failed')
    })
    it('allows the oracle to withdraw the funds it earned', async () => {
      const payment = 5
      VRFCoordinator.withdraw(h.personas.Nelly, payment, {
        from: h.personas.Neil,
      })
      assertBigNum(
        fee - payment,
        await VRFCoordinator.withdrawableTokens(h.personas.Neil),
      )
      assertBigNum( payment, await link.balanceOf(h.personas.Nelly))
    })
    it("doesn't allow people to withdraw more than their balance", async () => {
      h.assertActionThrows(async () =>
        VRFCoordinator.withdraw(h.personas.Neil, fee + 1, {
          from: h.personas.Neil,
        }),
      )
      h.assertActionThrows(async () =>
        VRFCoordinator.withdraw(h.personas.Nelly, fee + 1, {
          from: h.personas.Nelly,
        }),
      )
    })
  })
})
