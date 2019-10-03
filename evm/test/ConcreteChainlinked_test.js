import * as h from '../src/helpers'
const ConcreteChainlinked = artifacts.require('ConcreteChainlinked.sol')
const EmptyOracle = artifacts.require('EmptyOracle.sol')
const GetterSetter = artifacts.require('GetterSetter.sol')
const Oracle = artifacts.require('Oracle.sol')

let roles

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas()

  roles = rolesAndPersonas.roles
})

contract('ConcreteChainlinked', () => {
  const specId =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  let cc, gs, oc, newoc, link

  beforeEach(async () => {
    link = await h.linkContract(roles.defaultAccount)
    oc = await Oracle.new(link.address)
    newoc = await Oracle.new(link.address)
    gs = await GetterSetter.new()
    cc = await ConcreteChainlinked.new(link.address, oc.address)
  })

  describe('#newRequest', () => {
    it('forwards the information to the oracle contract through the link token', async () => {
      const tx = await cc.publicNewRequest(
        specId,
        gs.address,
        h.toHex('requestedBytes32(bytes32,bytes32)'),
      )

      assert.equal(1, tx.receipt.rawLogs.length)
      const [jId, cbAddr, cbFId, cborData] = h.decodeRunABI(
        tx.receipt.rawLogs[0],
      )
      const params = await h.decodeDietCBOR(cborData)

      assert.equal(specId.toLowerCase(), h.toHex(jId))
      assert.equal(gs.address.toLowerCase(), `0x${cbAddr.toLowerCase()}`)
      assert.equal('ed53e511', h.toHexWithoutPrefix(cbFId))
      assert.deepEqual({}, params)
    })
  })

  describe('#chainlinkRequest(Request)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      await cc.publicRequest(
        specId,
        cc.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const events = await h.getEvents(cc)
      assert.equal(1, events.length)
      const event = events[0]
      assert.equal(event.event, 'ChainlinkRequested')
    })
  })

  describe('#chainlinkRequestTo(Request)', () => {
    it('emits an event from the contract showing the run ID', async () => {
      await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const events = await h.getEvents(cc)
      assert.equal(1, events.length)
      const event = events[0]
      assert.equal(event.event, 'ChainlinkRequested')
    })

    it('emits an event on the target oracle contract', async () => {
      await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const events = await h.getEvents(newoc)
      assert.equal(1, events.length)
      const event = events[0]
      assert.equal(event.event, 'OracleRequest')
    })

    it('does not modify the stored oracle address', async () => {
      await cc.publicRequestRunTo(
        newoc.address,
        specId,
        cc.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )

      const actualOracleAddress = await cc.publicOracleAddress()
      assert.equal(oc.address, actualOracleAddress)
    })
  })

  describe('#cancelChainlinkRequest', () => {
    let requestId

    beforeEach(async () => {
      oc = await EmptyOracle.new()
      cc = await ConcreteChainlinked.new(link.address, oc.address)
      await cc.publicRequest(
        specId,
        cc.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      requestId = (await h.getLatestEvent(cc)).args.id
    })

    it('emits an event from the contract showing the run was cancelled', async () => {
      await cc.publicCancelRequest(requestId, 0, h.toHex(0), 0)

      const events = await h.getEvents(cc)
      assert.equal(1, events.length)
      const event = events[0]
      assert.equal(event.event, 'ChainlinkCancelled')
      assert.equal(requestId, event.args.id)
    })

    it('throws if given a bogus event ID', async () => {
      await h.assertActionThrows(async () => {
        await cc.publicCancelRequest(h.toHex('bogusId'), 0, h.toHex(0), 0)
      })
    })
  })

  describe('#recordChainlinkFulfillment(modifier)', () => {
    let request

    beforeEach(async () => {
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
    })

    it('emits an event marking the request fulfilled', async () => {
      await h.fulfillOracleRequest(oc, request, h.toHex('hi mom!'))

      const events = await h.getEvents(cc)
      assert.equal(1, events.length)
      const event = events[0]
      assert.equal(event.event, 'ChainlinkFulfilled')
      assert.equal(request.id, event.args.id)
    })
  })

  describe('#fulfillChainlinkRequest(function)', () => {
    let request

    beforeEach(async () => {
      const tx = await cc.publicRequest(
        specId,
        cc.address,
        h.toHex('publicFulfillChainlinkRequest(bytes32,bytes32)'),
        0,
      )
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
    })

    it('emits an event marking the request fulfilled', async () => {
      await h.fulfillOracleRequest(oc, request, h.toHex('hi mom!'))

      const events = await h.getEvents(cc)
      assert.equal(1, events.length)
      const event = events[0]
      assert.equal(event.event, 'ChainlinkFulfilled')
      assert.equal(request.id, event.args.id)
    })
  })

  describe('#chainlinkToken', () => {
    it('returns the Link Token address', async () => {
      const addr = await cc.publicChainlinkToken.call()
      assert.equal(addr, link.address)
    })
  })

  describe('#addExternalRequest', () => {
    let mock, request

    beforeEach(async () => {
      mock = await ConcreteChainlinked.new(link.address, oc.address)
      const tx = await cc.publicRequest(
        specId,
        mock.address,
        h.toHex('fulfillRequest(bytes32,bytes32)'),
        0,
      )
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
      await mock.publicAddExternalRequest(oc.address, request.id)
    })

    it('allows the external request to be fulfilled', async () => {
      await h.fulfillOracleRequest(oc, request, h.toHex('hi mom!'))
    })

    it('does not allow the same requestId to be used', async () => {
      await h.assertActionThrows(async () => {
        await cc.publicAddExternalRequest(newoc.address, request.id)
      })
    })
  })
})
