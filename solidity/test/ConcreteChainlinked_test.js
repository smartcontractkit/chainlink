import cbor from 'cbor'
import {
  assertActionThrows,
  bigNum,
  decodeRunABI,
  deploy,
  getEvents,
  getLatestEvent,
  toHexWithoutPrefix,
  toWei,
  increaseTime5Minutes
} from './support/helpers'

contract('ConcreteChainlinked', () => {
  const sourcePath = 'examples/ConcreteChainlinked.sol'
  let specId = '4c7b7ffb66b344fbaa64995af81e355a'
  let cc, gs, oc, newoc, link

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    oc = await deploy('Oracle.sol', link.address)
    newoc = await deploy('Oracle.sol', link.address)
    gs = await deploy('examples/GetterSetter.sol')
    cc = await deploy(sourcePath, link.address, oc.address)
  })

  describe('#newRun', () => {
    it('forwards the information to the oracle contract through the link token', async () => {
      let tx = await cc.publicNewRun(
        specId,
        gs.address,
        'requestedBytes32(bytes32,bytes32)')

      assert.equal(1, tx.receipt.logs.length)
      let [jId, cbAddr, cbFId, cborData] = decodeRunABI(tx.receipt.logs[0])
      let params = await cbor.decodeFirst(cborData)

      assert.equal(specId, jId)
      assert.equal(gs.address, `0x${cbAddr}`)
      assert.equal('ed53e511', toHexWithoutPrefix(cbFId))
      assert.deepEqual({}, params)
    })
  })

  describe('#chainlinkRequest(Run)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      await cc.publicRequestRun(specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)

      let events = await getEvents(cc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'ChainlinkRequested')
    })
  })

  describe('#chainlinkRequestFrom(Run)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      await cc.publicRequestRunFrom(newoc.address, specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)

      let events = await getEvents(cc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'ChainlinkRequested')
    })

    it('emits an event on the target oracle contract', async () => {
      await cc.publicRequestRunFrom(newoc.address, specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)

      let events = await getEvents(newoc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'RunRequest')
    })

    it('does not modify the stored oracle address', async () => {
      await cc.publicRequestRunFrom(newoc.address, specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)

      const actualOracleAddress = await cc.publicOracleAddress()
      assert.equal(oc.address, actualOracleAddress)
    })
  })

  describe('#cancelChainlinkRequest', () => {
    let requestId

    beforeEach(async () => {
      oc = await deploy('examples/EmptyOracle.sol')
      cc = await deploy(sourcePath, link.address, oc.address)

      await cc.publicRequestRun(specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)
      requestId = (await getLatestEvent(cc)).args.id
    })

    it('emits an event from the contract showing the run was cancelled', async () => {
      const tx = await cc.publicCancelRequest(requestId)

      const events = await getEvents(cc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'ChainlinkCancelled')
      assert.equal(requestId, event.args.id)
    })

    it('throws if given a bogus event ID', async () => {
      await assertActionThrows(async () => {
        await cc.publicCancelRequest("bogus ID")
      })
    })
  })

  describe('#checkChainlinkFulfillment(modifier)', () => {
    let internalId, requestId

    beforeEach(async () => {
      await cc.publicRequestRun(specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)
      requestId = (await getLatestEvent(cc)).args.id
      internalId = (await getLatestEvent(oc)).args.requestId
    })

    it('emits an event marking the request cancelled', async () => {
      await oc.fulfillData(internalId, 'hi mom!')

      let events = await getEvents(cc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'ChainlinkFulfilled')
      assert.equal(requestId, event.args.id)
    })
  })

  describe('#chainlinkToken', () => {
    it('returns the Link Token address', async () => {
      const addr = await cc.publicChainlinkToken.call();
      assert.equal(addr, link.address)
    })
  })

  describe('#LINK', () => {
    it('multiplies the value by a trillion', async () => {
      await cc.publicLINK(1)
      const event = await getLatestEvent(cc)
      assert.isTrue(event.args.amount.equals(toWei(1)))
    })

    it('throws an error if overflowing', async () => {
      let overflowAmount = bigNum('1157920892373161954235709850086879078532699846656405640394575')
      await assertActionThrows(async () => {
        await cc.publicLINK(overflowAmount)
      })
      const events = await getEvents(cc)
      assert.equal(0, events.length)
    })
  })

  describe('#addExternalRequest', () => {
    let mock, requestId

    beforeEach(async () => {
      mock = await deploy(sourcePath, link.address, oc.address)
      await cc.publicRequestRun(specId, mock.address, 'fulfillRequest(bytes32,bytes32)', 0)
      requestId = (await getLatestEvent(cc)).args.id
      await mock.publicAddExternalRequest(oc.address, requestId)
    })

    it('allows the external request to be fulfilled', async () => {
      await oc.fulfillData(requestId, 'hi mom!')
    })

    it('does not allow the same requestId to be used', async () => {
      await assertActionThrows(async () => {
        await cc.publicAddExternalRequest(newoc.address, requestId)
      })
    })
  })
})
