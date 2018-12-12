import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

contract('Oracle', () => {
  const sourcePath = 'Oracle.sol'
  const fHash = h.functionSelector('requestedBytes32(bytes32,bytes32)')
  const specId = '4c7b7ffb66b344fbaa64995af81e355a'
  const to = '0x80e29acb842498fe6591f020bd82766dce619d43'
  let link, oc

  beforeEach(async () => {
    link = await h.deploy('link_token/contracts/LinkToken.sol')
    oc = await h.deploy(sourcePath, link.address)
    await oc.setFulfillmentPermission(h.oracleNode, true, { from: h.defaultAccount })
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(sourcePath), [
      'EXPIRY_TIME',
      'cancel',
      'fulfillData',
      'getAuthorizationStatus',
      'onTokenTransfer',
      'owner',
      'renounceOwnership',
      'requestData',
      'setFulfillmentPermission',
      'transferOwnership',
      'withdraw',
      'withdrawable'
    ])
  })

  describe('#setFulfillmentPermission', () => {
    context('when called by the owner', () => {
      beforeEach(async () => {
        await oc.setFulfillmentPermission(h.stranger, true, { from: h.defaultAccount })
      })

      it('adds an authorized node', async () => {
        let authorized = await oc.getAuthorizationStatus(h.stranger)
        assert.equal(true, authorized)
      })

      it('removes an authorized node', async () => {
        await oc.setFulfillmentPermission(h.stranger, false, { from: h.defaultAccount })
        let authorized = await oc.getAuthorizationStatus(h.stranger)
        assert.equal(false, authorized)
      })
    })

    context('when called by a non-owner', () => {
      it('cannot add an authorized node', async () => {
        await h.assertActionThrows(async () => {
          await oc.setFulfillmentPermission(h.stranger, true, { from: h.stranger })
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

        let tx = await link.transferAndCall(oc.address, 0, callData)
        assert.equal(3, tx.receipt.logs.length)
      })

      context('with no data', () => {
        it('reverts', async () => {
          await h.assertActionThrows(async () => {
            await link.transferAndCall(oc.address, 0, '')
          })
        })
      })
    })

    context('malicious requester', () => {
      let mock, requester
      const paymentAmount = web3.toWei('1', 'ether')

      beforeEach(async () => {
        mock = await h.deploy('examples/MaliciousRequester.sol', link.address, oc.address)
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

        assert.isTrue(ocOriginalBalance.equals(ocNewBalance))
        assert.isTrue(mockNewBalance.equals(mockOriginalBalance))
      })

      context('if the requester tries to create a requestId for another contract', () => {
        let specId = h.newHash('0x4c7b7ffb66b344fbaa64995af81e355a')

        it('the requesters ID will not match with the oracle contract', async () => {
          const tx = await mock.maliciousTargetConsumer(to)
          let events = await h.getEvents(oc)
          const mockRequestId = tx.receipt.logs[0].data
          const requestId = events[0].args.requestId
          assert.notEqual(mockRequestId, requestId)
        })

        it('the target requester can still create valid requests', async () => {
          requester = await h.deploy('examples/BasicConsumer.sol', link.address, oc.address, h.toHex(specId))
          await link.transfer(requester.address, paymentAmount)
          await mock.maliciousTargetConsumer(requester.address)
          await requester.requestEthereumPrice('USD')
        })
      })
    })

    it('does not allow recursive calls of onTokenTransfer', async () => {
      const requestPayload = h.requestDataBytes(specId, to, fHash, 'id', '')

      const ottSelector = h.functionSelector('onTokenTransfer(address,uint256,bytes)')
      const header = '000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef' + // to
        '0000000000000000000000000000000000000000000000000000000000000539' + // amount
        '0000000000000000000000000000000000000000000000000000000000000060' + // offset
        '0000000000000000000000000000000000000000000000000000000000000136' //   length

      const maliciousPayload = ottSelector + header + requestPayload.slice(2)

      await h.assertActionThrows(async () => {
        await link.transferAndCall(oc.address, 0, maliciousPayload)
      })
    })
  })

  describe('#requestData', () => {
    context('when called through the LINK token', () => {
      const paid = 100
      let log, tx

      beforeEach(async () => {
        let args = h.requestDataBytes(specId, to, fHash, 1, '')
        tx = await h.requestDataFrom(oc, link, paid, args)
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2]
      })

      it('logs an event', async () => {
        assert.equal(oc.address, log.address)

        assert.equal(specId, web3.toUtf8(log.topics[1]))
        assert.equal(h.defaultAccount, web3.toDecimal(log.topics[2]))
        assert.equal(paid, web3.toDecimal(log.topics[3]))
      })

      it('uses the expected event signature', async () => {
        // If updating this test, be sure to update services.RunLogTopic.
        let eventSignature = '0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529'
        assert.equal(eventSignature, log.topics[0])
      })

      it('does not allow the same requestId to be used twice', async () => {
        let args2 = h.requestDataBytes(specId, to, fHash, 1, '')
        await h.assertActionThrows(async () => {
          await h.requestDataFrom(oc, link, paid, args2)
        })
      })
    })

    context('when not called through the LINK token', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await oc.requestData(0, 0, 1, specId, to, fHash, 1, '', { from: h.oracleNode })
        })
      })
    })
  })

  describe('#fulfillData', () => {
    let mock, requestId
    const paymentAmount = web3.toWei(1)
    let currency = 'USD'

    context('cooperative consumer', () => {
      beforeEach(async () => {
        mock = await h.deploy('examples/BasicConsumer.sol', link.address, oc.address, specId)
        await link.transfer(mock.address, web3.toWei('1', 'ether'))
        await mock.requestEthereumPrice(currency)
        let event = await h.getLatestEvent(oc)
        requestId = event.args.requestId
      })

      context('when called by an unauthorized node', () => {
        it('raises an error', async () => {
          let unauthorized = await oc.getAuthorizationStatus(h.stranger)
          assert.equal(false, unauthorized)
          await h.assertActionThrows(async () => {
            await oc.fulfillData(requestId, 'Hello World!', { from: h.stranger })
          })
        })
      })

      context('when called by an authorized node', () => {
        it('raises an error if the request ID does not exist', async () => {
          await h.assertActionThrows(async () => {
            await oc.fulfillData(0xdeadbeef, 'Hello World!', { from: h.oracleNode })
          })
        })

        it('sets the value on the requested contract', async () => {
          await oc.fulfillData(requestId, 'Hello World!', { from: h.oracleNode })

          let currentValue = await mock.currentPrice.call()
          assert.equal('Hello World!', web3.toUtf8(currentValue))
        })

        it('does not allow a request to be fulfilled twice', async () => {
          await oc.fulfillData(requestId, 'First message!', { from: h.oracleNode })
          await h.assertActionThrows(async () => {
            await oc.fulfillData(requestId, 'Second message!!', { from: h.oracleNode })
          })
        })
      })

      context('when the oracle does not provide enough gas', () => {
        // if updating this defaultGasLimit, be sure it matches with the
        // defaultGasLimit specified in store/tx_manager.go
        const defaultGasLimit = 500000

        beforeEach(async () => {
          assertBigNum(h.bigNum(0), await oc.withdrawable.call())
        })

        it('does not allow the oracle to withdraw the payment', async () => {
          await h.assertActionThrows(async () => {
            await oc.fulfillData(requestId, 'Hello World!', {
              from: h.oracleNode,
              gas: 70000
            })
          })

          assertBigNum(h.bigNum(0), await oc.withdrawable.call())
        })

        it(`${defaultGasLimit} is enough to pass the gas requirement`, async () => {
          assertBigNum(h.bigNum(0), await oc.withdrawable.call())

          await oc.fulfillData(requestId, 'Hello World!', {
            from: h.oracleNode,
            gas: defaultGasLimit
          })

          assertBigNum(h.bigNum(paymentAmount), await oc.withdrawable.call())
        })
      })
    })

    context('with a malicious requester', () => {
      const paymentAmount = h.toWei(1)

      beforeEach(async () => {
        mock = await h.deploy('examples/MaliciousRequester.sol', link.address, oc.address)
        await link.transfer(mock.address, paymentAmount)
      })

      it('cannot cancel before the expiration', async () => {
        await h.assertActionThrows(async () => {
          await mock.maliciousRequestCancel(specId, 'doesNothing(bytes32,bytes32)')
        })
      })

      it('cannot call functions on the LINK token through callbacks', async () => {
        await h.assertActionThrows(async () => {
          await mock.request(specId, link.address, 'transfer(address,uint256)')
        })
      })

      context('requester lies about amount of LINK sent', () => {
        it('the oracle uses the amount of LINK actually paid', async () => {
          const req = await mock.maliciousPrice(specId)
          const log = req.receipt.logs[3]

          assert.equal(web3.toWei(1), web3.toDecimal(log.topics[3]))
        })
      })
    })

    context('with a malicious consumer', () => {
      const paymentAmount = h.toWei(1)

      beforeEach(async () => {
        mock = await h.deploy('examples/MaliciousConsumer.sol', link.address, oc.address)
        await link.transfer(mock.address, paymentAmount)
      })

      context('fails during fulfillment', () => {
        beforeEach(async () => {
          await mock.requestData(specId, 'assertFail(bytes32,bytes32)')
          let events = await h.getEvents(oc)
          requestId = events[0].args.requestId
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })

          const balance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(h.oracleNode, paymentAmount, { from: h.defaultAccount })
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })
          await h.assertActionThrows(async () => {
            await oc.fulfillData(requestId, 'hack the planet 102', { from: h.oracleNode })
          })
        })
      })

      context('calls selfdestruct', () => {
        beforeEach(async () => {
          await mock.requestData(specId, 'doesNothing(bytes32,bytes32)')
          let events = await h.getEvents(oc)
          requestId = events[0].args.requestId
          await mock.remove()
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })

          const balance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(h.oracleNode, paymentAmount, { from: h.defaultAccount })
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })
      })

      context('request is canceled during fulfillment', () => {
        beforeEach(async () => {
          const req = await mock.requestData(specId, 'cancelRequestOnFulfill(bytes32,bytes32)')
          let event = await h.getLatestEvent(oc)
          requestId = event.args.requestId

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))

          const balance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(h.oracleNode, paymentAmount, { from: h.defaultAccount })
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })
          await h.assertActionThrows(async () => {
            await oc.fulfillData(requestId, 'hack the planet 102', { from: h.oracleNode })
          })
        })
      })

      context('tries to steal funds from node', () => {
        it('is not successful with call', async () => {
          await mock.requestData(specId, 'stealEthCall(bytes32,bytes32)')
          let events = await h.getEvents(oc)
          requestId = events[0].args.requestId

          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })
          const mockBalance = web3.fromWei(web3.eth.getBalance(mock.address))
          assert.equal(mockBalance, 0)
        })

        it('is not successful with send', async () => {
          await mock.requestData(specId, 'stealEthSend(bytes32,bytes32)')
          let events = await h.getEvents(oc)
          requestId = events[0].args.requestId

          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })
          const mockBalance = web3.fromWei(web3.eth.getBalance(mock.address))
          assert.equal(mockBalance, 0)
        })

        it('is not successful with transfer', async () => {
          await mock.requestData(specId, 'stealEthTransfer(bytes32,bytes32)')
          let events = await h.getEvents(oc)
          requestId = events[0].args.requestId

          await oc.fulfillData(requestId, 'hack the planet 101', { from: h.oracleNode })
          const mockBalance = web3.fromWei(web3.eth.getBalance(mock.address))
          assert.equal(mockBalance, 0)
        })
      })
    })
  })

  describe('#withdraw', () => {
    context('without reserving funds via requestData', () => {
      it('does nothing', async () => {
        let balance = await link.balanceOf(h.oracleNode)
        assert.equal(0, balance)
        await h.assertActionThrows(async () => {
          await oc.withdraw(h.oracleNode, h.toWei(1), { from: h.defaultAccount })
        })
        balance = await link.balanceOf(h.oracleNode)
        assert.equal(0, balance)
      })
    })

    context('reserving funds via requestData', () => {
      let requestId, amount

      beforeEach(async () => {
        amount = 15
        const mock = await h.deploy('examples/GetterSetter.sol')
        const args = h.requestDataBytes(specId, mock.address, fHash, 'id', '')
        const tx = await h.requestDataFrom(oc, link, amount, args)
        assert.equal(3, tx.receipt.logs.length)
        requestId = h.runRequestId(tx.receipt.logs[2])
      })

      context('but not freeing funds w fulfillData', () => {
        it('does not transfer funds', async () => {
          await h.assertActionThrows(async () => {
            await oc.withdraw(h.oracleNode, amount, { from: h.defaultAccount })
          })
          let balance = await link.balanceOf(h.oracleNode)
          assert.equal(0, balance)
        })
      })

      context('and freeing funds', () => {
        beforeEach(async () => {
          await oc.fulfillData(requestId, 'Hello World!', { from: h.oracleNode })
        })

        it('does not allow input greater than the balance', async () => {
          let originalOracleBalance = await link.balanceOf(oc.address)
          let originalStrangerBalance = await link.balanceOf(h.stranger)
          let withdrawAmount = amount + 1

          assert.isAbove(withdrawAmount, originalOracleBalance.toNumber())
          await h.assertActionThrows(async () => {
            await oc.withdraw(h.stranger, withdrawAmount, { from: h.defaultAccount })
          })

          let newOracleBalance = await link.balanceOf(oc.address)
          let newStrangerBalance = await link.balanceOf(h.stranger)

          assert.equal(originalOracleBalance.toNumber(), newOracleBalance.toNumber())
          assert.equal(originalStrangerBalance.toNumber(), newStrangerBalance.toNumber())
        })

        it('allows transfer of partial balance by owner to specified address', async () => {
          let partialAmount = 6
          let difference = amount - partialAmount
          await oc.withdraw(h.stranger, partialAmount, { from: h.defaultAccount })
          let strangerBalance = await link.balanceOf(h.stranger)
          let oracleBalance = await link.balanceOf(oc.address)
          assert.equal(partialAmount, strangerBalance)
          assert.equal(difference, oracleBalance)
        })

        it('allows transfer of entire balance by owner to specified address', async () => {
          await oc.withdraw(h.stranger, amount, { from: h.defaultAccount })
          let balance = await link.balanceOf(h.stranger)
          assert.equal(amount, balance)
        })

        it('does not allow a transfer of funds by non-owner', async () => {
          await h.assertActionThrows(async () => {
            await oc.withdraw(h.stranger, amount, { from: h.stranger })
          })
          let balance = await link.balanceOf(h.stranger)
          assert.equal(0, balance)
        })
      })
    })
  })

  describe('#withdrawable', () => {
    let internalId, amount

    beforeEach(async () => {
      amount = web3.toWei('1', 'ether')
      const mock = await h.deploy('examples/GetterSetter.sol')
      const args = h.requestDataBytes(specId, mock.address, fHash, 'id', '')
      const tx = await h.requestDataFrom(oc, link, amount, args)
      assert.equal(3, tx.receipt.logs.length)
      internalId = h.runRequestId(tx.receipt.logs[2])
      await oc.fulfillData(internalId, 'Hello World!', { from: h.oracleNode })
    })

    it('returns the correct value', async () => {
      const withdrawAmount = await oc.withdrawable.call()
      assert.equal(withdrawAmount.toNumber(), amount)
    })
  })

  describe('#cancel', () => {
    context('with no pending requests', () => {
      it('fails', async () => {
        await h.increaseTime5Minutes()
        await h.assertActionThrows(async () => {
          await oc.cancel(1337, { from: h.stranger })
        })
      })
    })

    context('with a pending request', () => {
      let requestId, tx, mock, requestAmount, startingBalance
      assert(mock === undefined, 'silence linter')
      beforeEach(async () => {
        startingBalance = 100
        requestAmount = 20

        mock = await h.deploy('examples/GetterSetter.sol')
        await link.transfer(h.consumer, startingBalance)

        let args = h.requestDataBytes(specId, h.consumer, fHash, 1, '')
        tx = await link.transferAndCall(oc.address, requestAmount, args, { from: h.consumer })
        assert.equal(3, tx.receipt.logs.length)
        requestId = h.runRequestId(tx.receipt.logs[2])
      })

      it('has correct initial balances', async () => {
        let oracleBalance = await link.balanceOf(oc.address)
        assert.equal(requestAmount, oracleBalance)

        let consumerAmount = await link.balanceOf(h.consumer)
        assert.equal(startingBalance - requestAmount, consumerAmount)
      })

      context('from a stranger', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await oc.cancel(requestId, { from: h.stranger })
          })
        })
      })

      context('from the requester', () => {
        it('refunds the correct amount', async () => {
          await h.increaseTime5Minutes()
          await oc.cancel(requestId, { from: h.consumer })
          let balance = await link.balanceOf(h.consumer)
          assert.equal(startingBalance, balance) // 100
        })

        it('triggers a cancellation event', async () => {
          await h.increaseTime5Minutes()
          const tx = await oc.cancel(requestId, { from: h.consumer })

          assert.equal(tx.receipt.logs.length, 2)
          assert.equal(requestId, tx.receipt.logs[1].data)
        })

        context('canceling twice', () => {
          it('fails', async () => {
            await h.increaseTime5Minutes()
            await oc.cancel(requestId, { from: h.consumer })
            await h.assertActionThrows(async () => {
              await oc.cancel(requestId, { from: h.consumer })
            })
          })
        })
      })
    })
  })
})
