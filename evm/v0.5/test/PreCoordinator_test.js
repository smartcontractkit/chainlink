import cbor from 'cbor'
import * as h from './support/helpers'
import { expectEvent, expectRevert, time } from 'openzeppelin-test-helpers'

contract('PreCoordinator', accounts => {
  const Oracle = artifacts.require('Oracle.sol')
  const PreCoordinator = artifacts.require('PreCoordinator.sol')
  const RequesterConsumer = artifacts.require('BasicConsumer.sol')

  const defaultAccount = accounts[0]
  const oracleNode1 = accounts[1]
  const oracleNode2 = accounts[2]
  const oracleNode3 = accounts[3]
  const oracleNode4 = accounts[4]
  const stranger = accounts[5]
  const consumer = accounts[6]

  // These parameters are used to validate the data was received
  // on the deployed oracle contract. The Job ID only represents
  // the type of data, but will not work on a public testnet.
  // For the latest JobIDs, visit our docs here:
  // https://docs.chain.link/docs/testnet-oracles
  const job1 = web3.utils.toHex('4c7b7ffb66b344fbaa64995af81e355a')
  const job2 = web3.utils.toHex('4c7b7ffb66b344fbaa64995af81e355b')
  const job3 = web3.utils.toHex('4c7b7ffb66b344fbaa64995af81e355c')
  const job4 = web3.utils.toHex('4c7b7ffb66b344fbaa64995af81e355d')
  const currency = 'USD'

  // Represents 1 LINK for testnet requests
  const payment = web3.utils.toWei('1')
  const totalPayment = web3.utils.toWei('4')

  const fulfilledEventSig = web3.utils.soliditySha3(
    'ChainlinkFulfilled(bytes32)'
  )

  let link, oc1, oc2, oc3, oc4, rc, pc

  beforeEach(async () => {
    link = await h.linkContract(defaultAccount)
    oc1 = await Oracle.new(link.address, { from: defaultAccount })
    oc2 = await Oracle.new(link.address, { from: defaultAccount })
    oc3 = await Oracle.new(link.address, { from: defaultAccount })
    oc4 = await Oracle.new(link.address, { from: defaultAccount })
    pc = await PreCoordinator.new(link.address, { from: defaultAccount })
    await oc1.setFulfillmentPermission(oracleNode1, true, {
      from: defaultAccount
    })
    await oc2.setFulfillmentPermission(oracleNode2, true, {
      from: defaultAccount
    })
    await oc3.setFulfillmentPermission(oracleNode3, true, {
      from: defaultAccount
    })
    await oc4.setFulfillmentPermission(oracleNode4, true, {
      from: defaultAccount
    })
  })

  describe('#createServiceAgreement', () => {
    context('when called by the owner', () => {
      it('emits the NewServiceAgreement log', async () => {
        const { logs } = await pc.createServiceAgreement(
          totalPayment,
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
          { from: defaultAccount }
        )
        expectEvent.inLogs(logs, 'NewServiceAgreement')
      })

      it('creates a service agreement', async () => {
        const tx = await pc.createServiceAgreement(
          totalPayment,
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
          { from: defaultAccount }
        )
        const saId = tx.receipt.rawLogs[0].topics[1]
        const sa = await pc.getServiceAgreement.call(saId)
        assert.equal(sa.totalPayment.toString(), totalPayment)
        assert.equal(sa.minResponses.toString(), 3)
        assert.deepEqual(sa.oracles, [
          oc1.address,
          oc2.address,
          oc3.address,
          oc4.address
        ])
        assert.deepEqual(sa.jobIds, [job1, job2, job3, job4])
        assert.equal(sa.payments.toString(), [
          payment,
          payment,
          payment,
          payment
        ])
      })
    })

    context('when called by a stranger', () => {
      it('reverts', async () => {
        await expectRevert.unspecified(
          pc.createServiceAgreement(
            totalPayment,
            3,
            [oc1.address, oc2.address, oc3.address, oc4.address],
            [job1, job2, job3, job4],
            [payment, payment, payment, payment],
            { from: stranger }
          )
        )
      })
    })

    context('when the array lengths are not equal', () => {
      it('reverts', async () => {
        await expectRevert(
          pc.createServiceAgreement(
            totalPayment,
            3,
            [oc1.address, oc2.address, oc3.address, oc4.address],
            [job1, job2, job3],
            [payment, payment, payment, payment],
            { from: defaultAccount }
          ),
          'Unmet length'
        )
      })
    })

    context('when the min responses is greater than the oracles', () => {
      it('reverts', async () => {
        await expectRevert(
          pc.createServiceAgreement(
            totalPayment,
            5,
            [oc1.address, oc2.address, oc3.address, oc4.address],
            [job1, job2, job3, job4],
            [payment, payment, payment, payment],
            { from: defaultAccount }
          ),
          'Invalid min responses'
        )
      })
    })
  })

  describe('#deleteServiceAgreement', () => {
    let saId

    beforeEach(async () => {
      const tx = await pc.createServiceAgreement(
        totalPayment,
        3,
        [oc1.address, oc2.address, oc3.address, oc4.address],
        [job1, job2, job3, job4],
        [payment, payment, payment, payment],
        { from: defaultAccount }
      )
      saId = tx.receipt.rawLogs[0].topics[1]
    })

    context('when called by a stranger', () => {
      it('reverts', async () => {
        await expectRevert.unspecified(
          pc.deleteServiceAgreement(saId, { from: stranger })
        )
      })
    })

    context('when called by the owner', () => {
      it('deletes the service agreement', async () => {
        await pc.deleteServiceAgreement(saId, { from: defaultAccount })
        const sa = await pc.getServiceAgreement(saId)
        assert.equal(sa.totalPayment, 0)
        assert.equal(sa.minResponses, 0)
        assert.deepEqual(sa.oracles, [])
        assert.deepEqual(sa.jobIds, [])
        assert.deepEqual(sa.payments, [])
      })
    })
  })

  describe('#onTokenTransfer', () => {
    context('when called by an address other than the LINK token', () => {
      it('reverts', async () => {
        let notLink = await h.linkContract(defaultAccount)
        let saId = await pc.createServiceAgreement.call(
          totalPayment,
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
          { from: defaultAccount }
        )
        let badRc = await RequesterConsumer.new(
          notLink.address,
          pc.address,
          saId,
          { from: consumer }
        )
        await notLink.transfer(badRc.address, totalPayment, {
          from: defaultAccount
        })
        await expectRevert.unspecified(
          badRc.requestEthereumPrice(currency, totalPayment, { from: consumer })
        )
      })
    })

    context('when called by the LINK token', () => {
      let saId
      beforeEach(async () => {
        const tx = await pc.createServiceAgreement(
          totalPayment,
          3,
          [oc1.address, oc2.address, oc3.address, oc4.address],
          [job1, job2, job3, job4],
          [payment, payment, payment, payment],
          { from: defaultAccount }
        )
        saId = tx.receipt.rawLogs[0].topics[1]
        rc = await RequesterConsumer.new(link.address, pc.address, saId, {
          from: consumer
        })
        await link.transfer(rc.address, totalPayment)
      })

      it('creates Chainlink requests', async () => {
        let tx = await rc.requestEthereumPrice(currency, totalPayment, {
          from: consumer
        })
        const log1 = tx.receipt.rawLogs[7]
        assert.equal(oc1.address, log1.address)
        const request1 = h.decodeRunRequest(log1)
        assert.equal(request1.requester, pc.address.toLowerCase())
        const log2 = tx.receipt.rawLogs[11]
        assert.equal(oc2.address, log2.address)
        const request2 = h.decodeRunRequest(log2)
        assert.equal(request2.requester, pc.address.toLowerCase())
        const log3 = tx.receipt.rawLogs[15]
        assert.equal(oc3.address, log3.address)
        const request3 = h.decodeRunRequest(log3)
        assert.equal(request3.requester, pc.address.toLowerCase())
        const log4 = tx.receipt.rawLogs[19]
        assert.equal(oc4.address, log4.address)
        const request4 = h.decodeRunRequest(log4)
        assert.equal(request4.requester, pc.address.toLowerCase())
        const expected = {
          path: ['USD'],
          get:
            'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY'
        }
        assert.deepEqual(expected, await cbor.decodeFirst(request1.data))
        assert.deepEqual(expected, await cbor.decodeFirst(request2.data))
        assert.deepEqual(expected, await cbor.decodeFirst(request3.data))
        assert.deepEqual(expected, await cbor.decodeFirst(request4.data))
      })

      context('when insufficient payment is supplied', () => {
        it('reverts', async () => {
          await expectRevert.unspecified(
            rc.requestEthereumPrice(currency, payment, { from: consumer })
          )
        })
      })

      context('when too much payment is supplied', () => {
        it('sends the extra back to the requester', async () => {
          await link.transfer(rc.address, payment)
          const extraPayment = web3.utils.toWei('5')
          const beforeBalance = await link.balanceOf(rc.address)
          assert.equal(beforeBalance, extraPayment)
          await rc.requestEthereumPrice(currency, extraPayment, {
            from: consumer
          })
          const afterBalance = await link.balanceOf(rc.address)
          assert.equal(afterBalance, payment)
        })
      })
    })
  })

  describe('#chainlinkCallback', () => {
    let saId, request1, request2, request3, request4
    const expected1 = 100
    const expected2 = 101
    const expected3 = 102
    const expected4 = 103
    const response1 = h.Ox(h.encodeInt256(expected1))
    const response2 = h.Ox(h.encodeInt256(expected2))
    const response3 = h.Ox(h.encodeInt256(expected3))
    const response4 = h.Ox(h.encodeInt256(expected4))

    beforeEach(async () => {
      const tx = await pc.createServiceAgreement(
        totalPayment,
        3,
        [oc1.address, oc2.address, oc3.address, oc4.address],
        [job1, job2, job3, job4],
        [payment, payment, payment, payment],
        { from: defaultAccount }
      )
      saId = tx.receipt.rawLogs[0].topics[1]
      rc = await RequesterConsumer.new(link.address, pc.address, saId, {
        from: consumer
      })
      await link.transfer(rc.address, totalPayment)

      const reqTx = await rc.requestEthereumPrice(currency, totalPayment, {
        from: consumer
      })
      const log1 = reqTx.receipt.rawLogs[7]
      request1 = h.decodeRunRequest(log1)
      const log2 = reqTx.receipt.rawLogs[11]
      request2 = h.decodeRunRequest(log2)
      const log3 = reqTx.receipt.rawLogs[15]
      request3 = h.decodeRunRequest(log3)
      const log4 = reqTx.receipt.rawLogs[19]
      request4 = h.decodeRunRequest(log4)
    })

    context('when called by a stranger', () => {
      it('reverts', async () => {
        await expectRevert(
          pc.chainlinkCallback(saId, response1),
          'Source must be the oracle of the request'
        )
      })
    })

    context('when called by the oracle contract', () => {
      it('records the answer', async () => {
        const tx = await oc1.fulfillOracleRequest(
          request1.id,
          request1.payment,
          request1.callbackAddr,
          request1.callbackFunc,
          request1.expiration,
          response1,
          { from: oracleNode1 }
        )
        assert.equal(tx.receipt.rawLogs[0].topics[0], fulfilledEventSig)
      })
    })

    context('when the minimum number of responses have returned', () => {
      beforeEach(async () => {
        await oc1.fulfillOracleRequest(
          request1.id,
          request1.payment,
          request1.callbackAddr,
          request1.callbackFunc,
          request1.expiration,
          response1,
          { from: oracleNode1 }
        )
        await oc2.fulfillOracleRequest(
          request2.id,
          request2.payment,
          request2.callbackAddr,
          request2.callbackFunc,
          request2.expiration,
          response2,
          { from: oracleNode2 }
        )
        await oc3.fulfillOracleRequest(
          request3.id,
          request3.payment,
          request3.callbackAddr,
          request3.callbackFunc,
          request3.expiration,
          response3,
          { from: oracleNode3 }
        )
      })

      it('returns the median to the requesting contract', async () => {
        const currentPrice = await rc.currentPrice.call()
        assert.equal(currentPrice.toString(), expected2)
      })

      context('when an oracle responds after aggregation', () => {
        it('does not update the requesting contract', async () => {
          await oc4.fulfillOracleRequest(
            request4.id,
            request4.payment,
            request4.callbackAddr,
            request4.callbackFunc,
            request4.expiration,
            response4,
            { from: oracleNode4 }
          )
          const currentPrice = await rc.currentPrice.call()
          assert.equal(currentPrice.toString(), expected2)
        })
      })
    })
  })

  describe('#withdrawLink', () => {
    beforeEach(async () => {
      await link.transfer(pc.address, payment)
      assert.equal(await link.balanceOf(pc.address), payment)
    })

    context('when called by a stranger', () => {
      it('reverts', async () => {
        await expectRevert.unspecified(pc.withdrawLink({ from: stranger }))
      })
    })

    context('when called by the owner', () => {
      it('allows the owner to withdraw LINK', async () => {
        await pc.withdrawLink({ from: defaultAccount })
        assert.equal(await link.balanceOf(pc.address), 0)
      })
    })
  })

  describe('#cancelOracleRequest', () => {
    let request

    beforeEach(async () => {
      const tx = await pc.createServiceAgreement(
        totalPayment,
        3,
        [oc1.address, oc2.address, oc3.address, oc4.address],
        [job1, job2, job3, job4],
        [payment, payment, payment, payment],
        { from: defaultAccount }
      )
      const saId = tx.receipt.rawLogs[0].topics[1]
      rc = await RequesterConsumer.new(link.address, pc.address, saId, {
        from: consumer
      })
      await link.transfer(rc.address, totalPayment)

      const reqTx = await rc.requestEthereumPrice(currency, totalPayment, {
        from: consumer
      })

      const log1 = reqTx.receipt.rawLogs[7]
      request = h.decodeRunRequest(log1)
    })

    context('before the minimum required time', () => {
      it('does not allow requests to be cancelled', async () => {
        await expectRevert(
          rc.cancelRequest(
            pc.address,
            request.id,
            request.payment,
            request.callbackFunc,
            request.expiration,
            { from: consumer }
          ),
          'Request is not expired'
        )
      })
    })

    context('after the minimum required time', () => {
      it('allows the requester to cancel', async () => {
        await time.increase(300)
        await rc.cancelRequest(
          pc.address,
          request.id,
          request.payment,
          request.callbackFunc,
          request.expiration,
          { from: consumer }
        )
        const balance = await link.balanceOf(rc.address)
        assert.equal(balance.toString(), payment)
      })

      it('allows others to call with the balance sent to the requester', async () => {
        await time.increase(300)
        await pc.cancelOracleRequest(
          request.id,
          request.payment,
          request.callbackFunc,
          request.expiration,
          { from: stranger }
        )
        const balance = await link.balanceOf(rc.address)
        assert.equal(balance.toString(), payment)
      })
    })
  })
})
