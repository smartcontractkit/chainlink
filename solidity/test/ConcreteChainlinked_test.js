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
  let cc, gs, oc, link

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    oc = await deploy('Oracle.sol', link.address)
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
      let tx = await cc.publicRequestRun(specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)

      let events = await getEvents(cc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'ChainlinkRequested')
    })
  })

  describe('#cancelChainlinkRequest', () => {
    let requestId

    beforeEach(async () => {
      await cc.publicRequestRun(specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)
      requestId = (await getLatestEvent(cc)).args.id
    })

    it('emits an event from the contract showing the run was cancelled', async () => {
      await increaseTime5Minutes();
      let tx = await cc.publicCancelRequest(requestId)

      let events = await getEvents(cc)
      assert.equal(1, events.length)
      let event = events[0]
      assert.equal(event.event, 'ChainlinkCancelled')
      assert.equal(requestId, event.args.id)
    })

    context('when the request ID is no longer unfulfilled', () => {
      beforeEach(async () => {
        await increaseTime5Minutes();
        await cc.publicCancelRequest(requestId)
      })

      it('throws an error', async () => {
        await assertActionThrows(async () => {
          await cc.publicCancelRequest(requestId)
        })
      })
    })
  })

  describe('#checkChainlinkFulfillment(modifier)', () => {
    let internalId, requestId

    beforeEach(async () => {
      await cc.publicRequestRun(specId, cc.address, 'fulfillRequest(bytes32,bytes32)', 0)
      requestId = (await getLatestEvent(cc)).args.id
      internalId = (await getLatestEvent(oc)).args.internalId
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
})
