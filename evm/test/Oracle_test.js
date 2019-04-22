import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

contract('Oracle', () => {
  const sourcePath = 'Oracle.sol'
  const fHash = h.functionSelector('requestedBytes32(bytes32,bytes32)')
  const specId = 
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000'
  const to = '0x80e29acb842498fe6591f020bd82766dce619d43'
  let link, oc, withdraw

  beforeEach(async () => {
    link = await h.linkContract()
    oc = await h.deploy(sourcePath, link.address)
    await oc.setFulfillmentPermission(h.oracleNode, true, {
      from: h.defaultAccount
    })
    withdraw = async (address, amount, options) =>
      oc.withdraw(address, amount.toString(), options)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(sourcePath), [
      'EXPIRY_TIME',
      'cancelOracleRequest',
      'fulfillOracleRequest',
      'getAuthorizationStatus',
      'onTokenTransfer',
      'owner',
      'renounceOwnership',
      'oracleRequest',
      'setFulfillmentPermission',
      'transferOwnership',
      'withdraw',
      'withdrawable'
    ])
  })

  describe('#setFulfillmentPermission', () => {
    context('when called by the owner', () => {
      beforeEach(async () => {
        await oc.setFulfillmentPermission(h.stranger, true, {
          from: h.defaultAccount
        })
      })

      it('adds an authorized node', async () => {
        let authorized = await oc.getAuthorizationStatus(h.stranger)
        assert.equal(true, authorized)
      })

      it('removes an authorized node', async () => {
        await oc.setFulfillmentPermission(h.stranger, false, {
          from: h.defaultAccount
        })
        let authorized = await oc.getAuthorizationStatus(h.stranger)
        assert.equal(false, authorized)
      })
    })

    context('when called by a non-owner', () => {
      it('cannot add an authorized node', async () => {
        await h.assertActionThrows(async () => {
          await oc.setFulfillmentPermission(h.stranger, true, {
            from: h.stranger
          })
        })
      })
    })
  })

  describe('#onTokenTransfer', () => {
    context('when called from any address but the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = h.requestDataBytes(specId, to, fHash, 'id', '')

        await h.assertActionThrows(async () => {
          await oc.onTokenTransfer(h.oracleNode, 0, callData)
        })
      })
    })

    context('when called from the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = h.requestDataBytes(specId, to, fHash, 'id', '')

        let tx = await link.transferAndCall(oc.address, 0, callData, {
          value: 0
        })
        assert.equal(3, tx.receipt.rawLogs.length)
      })

      context('with no data', () => {
        it('reverts', async () => {
          await h.assertActionThrows(async () => {
            await link.transferAndCall(oc.address, 0, '0x', {
              value: 0
            })
          })
        })
      })
    })

    context('malicious requester', () => {
      let mock, requester
      const paymentAmount = h.toWei('1', 'ether')

      beforeEach(async () => {
        mock = await h.deploy(
          'examples/MaliciousRequester.sol',
          link.address,
          oc.address
        )
        await link.transfer(mock.address, paymentAmount)
      })

      it('cannot withdraw from oracle', async () => {
        const ocOriginalBalance = await link.balanceOf.call(oc.address)
        const mockOriginalBalance = await link.balanceOf.call(mock.address)

        await h.assertActionThrows(async () => {
          await mock.maliciousWithdraw()
        })

        const ocNewBalance = await link.balanceOf.call(oc.address)
        const mockNewBalance = await link.balanceOf.call(mock.address)

        assertBigNum(ocOriginalBalance, ocNewBalance)
        assertBigNum(mockNewBalance, mockOriginalBalance)
      })

      context(
        'if the requester tries to create a requestId for another contract',
        () => {
          it('the requesters ID will not match with the oracle contract', async () => {
            const tx = await mock.maliciousTargetConsumer(to)
            let events = await h.getEvents(oc)
            const mockRequestId = tx.receipt.rawLogs[0].data
            const requestId = events[0].args.requestId
            assert.notEqual(mockRequestId, requestId)
          })

          it('the target requester can still create valid requests', async () => {
            requester = await h.deploy(
              'examples/BasicConsumer.sol',
              link.address,
              oc.address,
              specId
            )
            await link.transfer(requester.address, paymentAmount)
            await mock.maliciousTargetConsumer(requester.address)
            await requester.requestEthereumPrice('USD')
          })
        }
      )
    })

    it('does not allow recursive calls of onTokenTransfer', async () => {
      const requestPayload = h.requestDataBytes(specId, to, fHash, 'id', '')

      const ottSelector = h.functionSelector(
        'onTokenTransfer(address,uint256,bytes)'
      )
      const header =
        '000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef' + // to
        '0000000000000000000000000000000000000000000000000000000000000539' + // amount
        '0000000000000000000000000000000000000000000000000000000000000060' + // offset
        '0000000000000000000000000000000000000000000000000000000000000136' //   length

      const maliciousPayload = ottSelector + header + requestPayload.slice(2)

      await h.assertActionThrows(async () => {
        await link.transferAndCall(oc.address, 0, maliciousPayload, {
          value: 0
        })
      })
    })
  })

  describe('#oracleRequest', () => {
    context('when called through the LINK token', () => {
      const paid = 100
      let log, tx

      beforeEach(async () => {
        let args = h.requestDataBytes(specId, to, fHash, 1, '')
        tx = await h.requestDataFrom(oc, link, paid, args)
        assert.equal(3, tx.receipt.rawLogs.length)

        log = tx.receipt.rawLogs[2]
      })

      it('logs an event', async () => {
        assert.equal(oc.address, log.address)

        assert.equal(specId, log.topics[1])
        const req = h.decodeRunRequest(tx.receipt.rawLogs[2])
        assert.equal(h.defaultAccount.toString().toLowerCase(), req.requester)
        assertBigNum(paid, req.payment)
      })

      it('uses the expected event signature', async () => {
        // If updating this test, be sure to update models.RunLogTopic.
        const eventSignature =
          '0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65'
        assert.equal(eventSignature, log.topics[0])
      })

      it('does not allow the same requestId to be used twice', async () => {
        let args2 = h.requestDataBytes(specId, to, fHash, 1, '')
        await h.assertActionThrows(async () => {
          await h.requestDataFrom(oc, link, paid, args2)
        })
      })

      context(
        'when called with a payload less than 2 EVM words + function selector',
        () => {
          const funcSelector = h.functionSelector(
            'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
          )
          const maliciousData =
            funcSelector +
            '0000000000000000000000000000000000000000000000000000000000000000000'

          it('throws an error', async () => {
            await h.assertActionThrows(async () => {
              await h.requestDataFrom(oc, link, paid, maliciousData)
            })
          })
        }
      )

      context('when called with a payload between 3 and 9 EVM words', () => {
        const funcSelector = h.functionSelector(
          'oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)'
        )
        const maliciousData =
          funcSelector +
          '000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001'

        it('throws an error', async () => {
          await h.assertActionThrows(async () => {
            await h.requestDataFrom(oc, link, paid, maliciousData)
          })
        })
      })
    })

    context('when not called through the LINK token', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await oc.oracleRequest(
            '0x0000000000000000000000000000000000000000',
            0,
            specId,
            to,
            fHash,
            1,
            1,
            '0x',
            { from: h.oracleNode }
          )
        })
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = 'Hi Mom!'
    let mock, request

    context('cooperative consumer', () => {
      beforeEach(async () => {
        mock = await h.deploy(
          'examples/BasicConsumer.sol',
          link.address,
          oc.address,
          specId
        )
        const paymentAmount = h.toWei(1)
        await link.transfer(mock.address, paymentAmount)
        const currency = 'USD'
        const tx = await mock.requestEthereumPrice(currency)
        request = h.decodeRunRequest(tx.receipt.rawLogs[3])
      })

      context('when called by an unauthorized node', () => {
        beforeEach(async () => {
          assert.equal(false, await oc.getAuthorizationStatus(h.stranger))
        })

        it('raises an error', async () => {
          await h.assertActionThrows(async () => {
            await h.fulfillOracleRequest(oc, request, response, {
              from: h.stranger
            })
          })
        })
      })

      context('when called by an authorized node', () => {
        it('raises an error if the request ID does not exist', async () => {
          request.id = '0xdeadbeef'
          await h.assertActionThrows(async () => {
            await h.fulfillOracleRequest(oc, request, response, {
              from: h.oracleNode
            })
          })
        })

        it('sets the value on the requested contract', async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          const currentValue = await mock.currentPrice.call()
          assert.equal(response, h.toUtf8(currentValue))
        })

        it('does not allow a request to be fulfilled twice', async () => {
          const response2 = response + ' && Hello World!!'

          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          await h.assertActionThrows(async () => {
            await h.fulfillOracleRequest(oc, request, response2, {
              from: h.oracleNode
            })
          })

          const currentValue = await mock.currentPrice.call()
          assert.equal(response, h.toUtf8(currentValue))
        })
      })

      context('when the oracle does not provide enough gas', () => {
        // if updating this defaultGasLimit, be sure it matches with the
        // defaultGasLimit specified in store/tx_manager.go
        const defaultGasLimit = 500000

        beforeEach(async () => {
          assertBigNum(0, await oc.withdrawable.call())
        })

        it('does not allow the oracle to withdraw the payment', async () => {
          await h.assertActionThrows(async () => {
            await h.fulfillOracleRequest(oc, request, response, {
              from: h.oracleNode,
              gas: 70000
            })
          })

          assertBigNum(0, await oc.withdrawable.call())
        })

        it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode,
            gas: defaultGasLimit
          })

          assertBigNum(request.payment, await oc.withdrawable.call())
        })
      })
    })

    context('with a malicious requester', () => {
      beforeEach(async () => {
        const paymentAmount = h.toWei(1)
        mock = await h.deploy(
          'examples/MaliciousRequester.sol',
          link.address,
          oc.address
        )
        await link.transfer(mock.address, paymentAmount)
      })

      it('cannot cancel before the expiration', async () => {
        await h.assertActionThrows(async () => {
          await mock.maliciousRequestCancel(
            specId,
            h.toHex('doesNothing(bytes32,bytes32)')
          )
        })
      })

      it('cannot call functions on the LINK token through callbacks', async () => {
        await h.assertActionThrows(async () => {
          await mock.request(
            specId,
            link.address,
            h.toHex('transfer(address,uint256)')
          )
        })
      })

      context('requester lies about amount of LINK sent', () => {
        it('the oracle uses the amount of LINK actually paid', async () => {
          const tx = await mock.maliciousPrice(specId)
          const req = h.decodeRunRequest(tx.receipt.rawLogs[3])

          assert(h.toWei(1).eq(h.bigNum(req.payment)))
        })
      })
    })

    context('with a malicious consumer', () => {
      const paymentAmount = h.toWei(1)

      beforeEach(async () => {
        mock = await h.deploy(
          'examples/MaliciousConsumer.sol',
          link.address,
          oc.address
        )
        await link.transfer(mock.address, paymentAmount)
      })

      context('fails during fulfillment', () => {
        beforeEach(async () => {
          const tx = await mock.requestData(
            specId,
            h.toHex('assertFail(bytes32,bytes32)')
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])
        })

        it('allows the oracle node to receive their payment', async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          const balance = await link.balanceOf.call(h.oracleNode)
          assertBigNum(balance, h.bigNum(0))

          await withdraw(h.oracleNode, paymentAmount, {
            from: h.defaultAccount
          })
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assertBigNum(paymentAmount, newBalance)
        })

        it("can't fulfill the data again", async () => {
          const response2 = 'hack the planet 102'

          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          await h.assertActionThrows(async () => {
            await h.fulfillOracleRequest(oc, request, response2, {
              from: h.oracleNode
            })
          })
        })
      })

      context('calls selfdestruct', () => {
        beforeEach(async () => {
          const tx = await mock.requestData(
            specId,
            h.toHex('doesNothing(bytes32,bytes32)')
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])
          await mock.remove()
        })

        it('allows the oracle node to receive their payment', async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          const balance = await link.balanceOf.call(h.oracleNode)
          assertBigNum(balance, h.bigNum(0))

          await withdraw(h.oracleNode, paymentAmount, {
            from: h.defaultAccount
          })
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assertBigNum(paymentAmount, newBalance)
        })
      })

      context('request is canceled during fulfillment', () => {
        beforeEach(async () => {
          const tx = await mock.requestData(
            specId,
            h.toHex('cancelRequestOnFulfill(bytes32,bytes32)')
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])

          assertBigNum(0, await link.balanceOf.call(mock.address))
        })

        it('allows the oracle node to receive their payment', async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          const mockBalance = await link.balanceOf.call(mock.address)
          assertBigNum(mockBalance, h.bigNum(0))

          const balance = await link.balanceOf.call(h.oracleNode)
          assertBigNum(balance, h.bigNum(0))

          await withdraw(h.oracleNode, paymentAmount, {
            from: h.defaultAccount
          })
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assertBigNum(paymentAmount, newBalance)
        })

        it("can't fulfill the data again", async () => {
          const response2 = 'hack the planet 102'

          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          await h.assertActionThrows(async () => {
            await h.fulfillOracleRequest(oc, request, response2, {
              from: h.oracleNode
            })
          })
        })
      })

      context('tries to steal funds from node', () => {
        it('is not successful with call', async () => {
          const tx = await mock.requestData(
            specId,
            h.toHex('stealEthCall(bytes32,bytes32)')
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])

          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          assertBigNum(0, await web3.eth.getBalance(mock.address))
        })

        it('is not successful with send', async () => {
          const tx = await mock.requestData(
            specId,
            h.toHex('stealEthSend(bytes32,bytes32)')
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])

          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })
          assertBigNum(0, await web3.eth.getBalance(mock.address))
        })

        it('is not successful with transfer', async () => {
          const tx = await mock.requestData(
            specId,
            h.toHex('stealEthTransfer(bytes32,bytes32)')
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])

          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })
          assertBigNum(0, await web3.eth.getBalance(mock.address))
        })
      })
    })
  })

  describe('#withdraw', () => {
    context('without reserving funds via oracleRequest', () => {
      it('does nothing', async () => {
        let balance = await link.balanceOf(h.oracleNode)
        assert.equal(0, balance)
        await h.assertActionThrows(async () => {
          await withdraw(h.oracleNode, h.toWei(1), { from: h.defaultAccount })
        })
        balance = await link.balanceOf(h.oracleNode)
        assert.equal(0, balance)
      })
    })

    context('reserving funds via oracleRequest', () => {
      const payment = 15
      let request

      beforeEach(async () => {
        const mock = await h.deploy('examples/GetterSetter.sol')
        const args = h.requestDataBytes(specId, mock.address, fHash, 'id', '')
        const tx = await h.requestDataFrom(oc, link, payment, args)
        assert.equal(3, tx.receipt.rawLogs.length)
        request = h.decodeRunRequest(tx.receipt.rawLogs[2])
      })

      context('but not freeing funds w fulfillOracleRequest', () => {
        it('does not transfer funds', async () => {
          await h.assertActionThrows(async () => {
            await withdraw(h.oracleNode, payment, { from: h.defaultAccount })
          })
          let balance = await link.balanceOf(h.oracleNode)
          assert.equal(0, balance)
        })
      })

      context('and freeing funds', () => {
        beforeEach(async () => {
          await h.fulfillOracleRequest(oc, request, 'Hello World!', {
            from: h.oracleNode
          })
        })

        it('does not allow input greater than the balance', async () => {
          const originalOracleBalance = await link.balanceOf(oc.address)
          const originalStrangerBalance = await link.balanceOf(h.stranger)
          const withdrawalAmount = payment + 1

          assert.isAbove(withdrawalAmount, originalOracleBalance.toNumber())
          await h.assertActionThrows(async () => {
            await withdraw(h.stranger, withdrawalAmount, {
              from: h.defaultAccount
            })
          })

          const newOracleBalance = await link.balanceOf(oc.address)
          const newStrangerBalance = await link.balanceOf(h.stranger)

          assert.equal(
            originalOracleBalance.toNumber(),
            newOracleBalance.toNumber()
          )
          assert.equal(
            originalStrangerBalance.toNumber(),
            newStrangerBalance.toNumber()
          )
        })

        it('allows transfer of partial balance by owner to specified address', async () => {
          const partialAmount = 6
          const difference = payment - partialAmount
          await withdraw(h.stranger, partialAmount, { from: h.defaultAccount })
          const strangerBalance = await link.balanceOf(h.stranger)
          const oracleBalance = await link.balanceOf(oc.address)
          assert.equal(partialAmount, strangerBalance)
          assert.equal(difference, oracleBalance)
        })

        it('allows transfer of entire balance by owner to specified address', async () => {
          await withdraw(h.stranger, payment, { from: h.defaultAccount })
          const balance = await link.balanceOf(h.stranger)
          assert.equal(payment, balance)
        })

        it('does not allow a transfer of funds by non-owner', async () => {
          await h.assertActionThrows(async () => {
            await withdraw(h.stranger, payment, { from: h.stranger })
          })
          const balance = await link.balanceOf(h.stranger)
          assert.equal(0, balance)
        })
      })
    })
  })

  describe('#withdrawable', () => {
    let request

    beforeEach(async () => {
      const amount = h.toWei(1, 'ether').toString()
      const mock = await h.deploy('examples/GetterSetter.sol')
      const args = h.requestDataBytes(specId, mock.address, fHash, 'id', '')
      const tx = await h.requestDataFrom(oc, link, amount, args)
      assert.equal(3, tx.receipt.rawLogs.length)
      request = h.decodeRunRequest(tx.receipt.rawLogs[2])
      await h.fulfillOracleRequest(oc, request, 'Hello World!', {
        from: h.oracleNode
      })
    })

    it('returns the correct value', async () => {
      const withdrawAmount = await oc.withdrawable.call()
      assertBigNum(withdrawAmount, request.payment)
    })
  })

  describe('#cancelOracleRequest', () => {
    context('with no pending requests', () => {
      it('fails', async () => {
        const fakeRequest = {
          id: h.toHex(1337),
          payment: 0,
          callbackFunc: h.functionSelector('requestedBytes32(bytes32,bytes32)'),
          expiration: 999999999999
        }
        await h.increaseTime5Minutes()

        await h.assertActionThrows(async () => {
          await h.cancelOracleRequest(oc, fakeRequest, { from: h.stranger })
        })
      })
    })

    context('with a pending request', () => {
      const startingBalance = 100
      let request, tx

      beforeEach(async () => {
        const requestAmount = 20

        await link.transfer(h.consumer, startingBalance)

        let args = h.requestDataBytes(specId, h.consumer, fHash, 1, '')
        tx = await link.transferAndCall(oc.address, requestAmount, args, {
          from: h.consumer
        })
        assert.equal(3, tx.receipt.rawLogs.length)
        request = h.decodeRunRequest(tx.receipt.rawLogs[2])
      })

      it('has correct initial balances', async () => {
        let oracleBalance = await link.balanceOf(oc.address)
        assertBigNum(request.payment, oracleBalance)

        let consumerAmount = await link.balanceOf(h.consumer)
        assert.equal(startingBalance - request.payment, consumerAmount)
      })

      context('from a stranger', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await h.cancelOracleRequest(oc, request, { from: h.consumer })
          })
        })
      })

      context('from the requester', () => {
        it('refunds the correct amount', async () => {
          await h.increaseTime5Minutes()
          await h.cancelOracleRequest(oc, request, { from: h.consumer })
          let balance = await link.balanceOf(h.consumer)
          assert.equal(startingBalance, balance) // 100
        })

        it('triggers a cancellation event', async () => {
          await h.increaseTime5Minutes()
          const tx = await h.cancelOracleRequest(oc, request, {
            from: h.consumer
          })

          assert.equal(tx.receipt.rawLogs.length, 2)
          assert.equal(request.id, tx.receipt.rawLogs[0].topics[1])
        })

        it('fails when called twice', async () => {
          await h.increaseTime5Minutes()
          await h.cancelOracleRequest(oc, request, { from: h.consumer })

          await h.assertActionThrows(async () => {
            await h.cancelOracleRequest(oc, request, { from: h.consumer })
          })
        })
      })
    })
  })
})
