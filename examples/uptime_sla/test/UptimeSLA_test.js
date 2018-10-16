import cbor from 'cbor'
import util from 'ethereumjs-util'
import {
  days,
  fastForwardTo,
  getLatestTimestamp,
} from './support/helpers'
import {
  assertActionThrows,
  eth,
  functionSelector,
  getEvents,
  getLatestEvent,
  oracleNode,
  requestDataBytes,
  requestDataFrom
} from '../../../solidity/test/support/helpers'
import { assertBigNum } from '../../../solidity/test/support/matchers'

contract('UptimeSLA', () => {
  let Link = artifacts.require('LinkToken')
  let Oracle = artifacts.require('Oracle')
  let SLA = artifacts.require('UptimeSLA')
  let specId = '4c7b7ffb66b344fbaa64995af81e355a'
  let deposit = 1000000000
  let link, oc, sla, client, serviceProvider, startAt

  beforeEach(async () => {
    client = "0xf000000000000000000000000000000000000001"
    serviceProvider = "0xf000000000000000000000000000000000000002"
    link = await Link.new()
    oc = await Oracle.new(link.address, {from: oracleNode})
    sla = await SLA.new(client, serviceProvider, link.address, oc.address, specId, {
      value: deposit
    })
    link.transfer(sla.address, web3.toWei(1, 'ether'))
    startAt = await getLatestTimestamp()
  })

  describe("before updates", () => {
    it("does not release money to anyone", async () => {
      assert.equal(await eth.getBalance(sla.address), deposit)
      assert.equal(await eth.getBalance(client), 0)
      assert.equal(await eth.getBalance(serviceProvider), 0)
    })
  })

  describe("#updateUptime", () => {
    it("triggers a log event in the Oracle contract", async () => {
      const tx = await sla.updateUptime("0")

      const events = await getEvents(oc)
      assert.equal(1, events.length)

      const event = events[0]
      assert.equal(web3.toUtf8(event.args.specId), specId)

      const decoded = cbor.decodeFirstSync(util.toBuffer(event.args.data))
      assert.deepEqual(
        decoded,
        {'url': 'https://status.heroku.com/api/ui/availabilities', 'path': ['data', '0', 'attributes', 'calculation']}
      )
    })
  })

  describe('#fulfillData', () => {
    const response = '0x00000000000000000000000000000000000000000000000000000000000f8c4c'
    let requestId

    beforeEach(async () => {
      await sla.updateUptime('0')
      const event = await getLatestEvent(oc)
      requestId = event.args.internalId
    })

    context("when the value is below 9999", async () => {
      const response = "0x000000000000000000000000000000000000000000000000000000000000270e"

      it('sends the deposit to the client', async () => {
        await oc.fulfillData(requestId, response, {from: oracleNode})

        assert.equal(await eth.getBalance(sla.address), 0)
        assert.equal(await eth.getBalance(client), deposit)
        assert.equal(await eth.getBalance(serviceProvider), 0)
      })
    })

    context("when the value is 9999 or above", () => {
      const response = "0x000000000000000000000000000000000000000000000000000000000000270f"
      let originalClientBalance

      beforeEach(async () => {
        originalClientBalance = await eth.getBalance(client)
      })

      it('does not move the money', async () => {
        await oc.fulfillData(requestId, response, {from: oracleNode})

        assertBigNum(await eth.getBalance(sla.address), deposit)
        assertBigNum(await eth.getBalance(client), originalClientBalance)
        assertBigNum(await eth.getBalance(serviceProvider), 0)
      })

      context('and a month has passed', () => {
        beforeEach(async () => {
          await fastForwardTo(startAt + days(30))
        })

        it('gives the money back to the service provider', async () => {
          await oc.fulfillData(requestId, response, {from: oracleNode})

          assertBigNum(await eth.getBalance(sla.address), 0)
          assertBigNum(await eth.getBalance(client), originalClientBalance)
          assertBigNum(await eth.getBalance(serviceProvider), deposit)
        })
      })
    })

    context('when the consumer does not recognize the request ID', () => {
      beforeEach(async () => {
        let fid = functionSelector("report(uint256,bytes32)")
        let args = requestDataBytes(specId, sla.address, fid, "xid", "")
        await requestDataFrom(oc, link, 0, args)
        let event = await getLatestEvent(oc)
        requestId = event.args.internalId
      })

      it("does not accept the data provided", async () => {
        let originalUptime = await sla.uptime()
        await oc.fulfillData(requestId, response, {from: oracleNode})
        let newUptime = await sla.uptime()

        assert.isTrue(originalUptime.equals(newUptime))
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await assertActionThrows(async () => {
          await sla.report(requestId, response, {from: oracleNode})
        })
      })
    })
  })
})
