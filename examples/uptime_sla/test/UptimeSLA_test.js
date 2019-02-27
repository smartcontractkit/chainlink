import { toBuffer } from 'ethereumjs-util'
import { days, fastForwardTo, getLatestTimestamp } from './support/helpers'
import {
  assertActionThrows,
  decodeDietCBOR,
  decodeRunRequest,
  functionSelector,
  fulfillOracleRequest,
  requestDataBytes
} from './support/moreHelpers'
import { assertBigNum } from './support/matchers'

contract('UptimeSLA', accounts => {
  const Oracle = artifacts.require('Oracle')
  const SLA = artifacts.require('UptimeSLA')
  const LinkToken = artifacts.require('LinkToken')
  const specId = '0x4c7b7ffb66b344fbaa64995af81e355a'
  const deposit = 1000000000
  const oracleNode = accounts[1]
  let link, oc, sla, client, serviceProvider, startAt

  beforeEach(async () => {
    client = '0xf000000000000000000000000000000000000001'
    serviceProvider = '0xf000000000000000000000000000000000000002'
    link = await LinkToken.new()
    oc = await Oracle.new(link.address, { from: oracleNode })
    sla = await SLA.new(
      client,
      serviceProvider,
      link.address,
      oc.address,
      specId,
      {
        value: deposit
      }
    )
    link.transfer(sla.address, web3.utils.toWei('1', 'ether'))
    startAt = await getLatestTimestamp()
  })

  describe('before updates', () => {
    it('does not release money to anyone', async () => {
      assert.equal(await web3.eth.getBalance(sla.address), deposit)
      assert.equal(await web3.eth.getBalance(client), 0)
      assert.equal(await web3.eth.getBalance(serviceProvider), 0)
    })
  })

  describe('#updateUptime', () => {
    it('triggers a log event in the Oracle contract', async () => {
      await sla.updateUptime('0')

      const events = await oc.getPastEvents('allEvents')
      assert.equal(1, events.length)
      assert.equal(
        events[0].args.specId,
        specId + '00000000000000000000000000000000'
      )

      const decoded = await decodeDietCBOR(toBuffer(events[0].args.data))
      assert.deepEqual(decoded, {
        url: 'https://status.heroku.com/api/ui/availabilities',
        path: ['data', '0', 'attributes', 'calculation']
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response =
      '0x00000000000000000000000000000000000000000000000000000000000f8c4c'
    let request

    beforeEach(async () => {
      const tx = await sla.updateUptime('0')
      request = decodeRunRequest(tx.receipt.rawLogs[3])
    })

    context('when the value is below 9999', async () => {
      const response =
        '0x000000000000000000000000000000000000000000000000000000000000270e'

      it('sends the deposit to the client', async () => {
        await fulfillOracleRequest(oc, request, response, { from: oracleNode })

        assert.equal(await web3.eth.getBalance(sla.address), 0)
        assert.equal(await web3.eth.getBalance(client), deposit)
        assert.equal(await web3.eth.getBalance(serviceProvider), 0)
      })
    })

    context('when the value is 9999 or above', () => {
      const response =
        '0x000000000000000000000000000000000000000000000000000000000000270f'
      let originalClientBalance

      beforeEach(async () => {
        originalClientBalance = await web3.eth.getBalance(client)
      })

      it('does not move the money', async () => {
        await fulfillOracleRequest(oc, request, response, { from: oracleNode })

        assertBigNum(await web3.eth.getBalance(sla.address), deposit)
        assertBigNum(await web3.eth.getBalance(client), originalClientBalance)
        assertBigNum(await web3.eth.getBalance(serviceProvider), 0)
      })

      context('and a month has passed', () => {
        beforeEach(async () => {
          await fastForwardTo(startAt + days(30))
        })

        it('gives the money back to the service provider', async () => {
          await fulfillOracleRequest(oc, request, response, {
            from: oracleNode
          })

          assertBigNum(await web3.eth.getBalance(sla.address), 0)
          assertBigNum(await web3.eth.getBalance(client), originalClientBalance)
          assertBigNum(await web3.eth.getBalance(serviceProvider), deposit)
        })
      })
    })

    context('when the consumer does not recognize the request ID', () => {
      beforeEach(async () => {
        let fid = functionSelector('report(uint256,bytes32)')
        let args = requestDataBytes(specId, sla.address, fid, 'xid', 'foo')
        const tx = await link.transferAndCall(oc.address, 0, args)
        request = decodeRunRequest(tx.receipt.rawLogs[2])
      })

      it('does not accept the data provided', async () => {
        let originalUptime = await sla.uptime()
        await fulfillOracleRequest(oc, request, response, { from: oracleNode })
        let newUptime = await sla.uptime()

        assertBigNum(originalUptime, newUptime)
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await assertActionThrows(async () => {
          await sla.report(request.id, response, { from: oracleNode })
        })
      })
    })
  })
})
