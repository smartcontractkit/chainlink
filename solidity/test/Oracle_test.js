import {
  assertActionThrows,
  consumer,
  checkPublicABI,
  defaultAccount,
  deploy,
  functionSelector,
  getLatestEvent,
  oracleNode,
  requestDataBytes,
  requestDataFrom,
  runRequestId,
  stranger,
  toHex,
  toWei,
  increaseTime5Minutes
} from './support/helpers'

contract('Oracle', () => {
  const sourcePath = 'Oracle.sol'
  const fHash = functionSelector('fulfill(bytes32,bytes32)')
  const specId = '4c7b7ffb66b344fbaa64995af81e355a'
  const to = '0x80e29acb842498fe6591f020bd82766dce619d43'
  let link, oc

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    oc = await deploy(sourcePath, link.address)
    await oc.setFulfillmentPermission(oracleNode, true, { from: defaultAccount })
  })

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [
      'cancel',
      'fulfillData',
      'getAuthorizationStatus',
      'onTokenTransfer',
      'owner',
      'renounceOwnership',
      'requestData',
      'setFulfillmentPermission',
      'transferOwnership',
      'withdraw'
    ])
  })

  describe('#setFulfillmentPermission', () => {
    context('when called by the owner', () => {
      beforeEach(async () => {
        await oc.setFulfillmentPermission(stranger, true, { from: defaultAccount })
      })

      it('adds an authorized node', async () => {
        let authorized = await oc.getAuthorizationStatus(stranger)
        assert.equal(true, authorized)
      })

      it('removes an authorized node', async () => {
        await oc.setFulfillmentPermission(stranger, false, { from: defaultAccount })
        let authorized = await oc.getAuthorizationStatus(stranger)
        assert.equal(false, authorized)
      })
    })

    context('when called by a non-owner', () => {
      it('cannot add an authorized node', async () => {
        await assertActionThrows(async () => {
          await oc.setFulfillmentPermission(stranger, true, { from: stranger })
        })
      })
    })
  })

  describe('#onTokenTransfer', () => {
    context('when called from any address but the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = requestDataBytes(to, specId, fHash, 'id', '')

        await assertActionThrows(async () => {
          await oc.onTokenTransfer(oracleNode, 0, callData)
        })
      })
    })

    context('when called from the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = requestDataBytes(to, specId, fHash, 'id', '')

        let tx = await link.transferAndCall(oc.address, 0, callData)
        assert.equal(3, tx.receipt.logs.length)
      })

      context('with no data', () => {
        it('reverts', async () => {
          await assertActionThrows(async () => {
            await link.transferAndCall(oc.address, 0, '')
          })
        })
      })
    })

    context('malicious requester', () => {
      let mock
      const paymentAmount = 1

      beforeEach(async () => {
        mock = await deploy('examples/MaliciousRequester.sol', link.address, oc.address)
        await link.transfer(mock.address, paymentAmount)
      })

      it('cannot withdraw from oracle', async () => {
        const ocOriginalBalance = await link.balanceOf.call(oc.address)
        const mockOriginalBalance = await link.balanceOf.call(mock.address)

        await assertActionThrows(async () => {
          await mock.maliciousWithdraw()
        })

        const ocNewBalance = await link.balanceOf.call(oc.address)
        const mockNewBalance = await link.balanceOf.call(mock.address)

        assert.isTrue(ocOriginalBalance.equals(ocNewBalance))
        assert.isTrue(mockNewBalance.equals(mockOriginalBalance))
      })
    })

    it('does not allow recursive calls of onTokenTransfer', async () => {
      const requestPayload = requestDataBytes(to, specId, fHash, 'id', '')

      const ottSelector = functionSelector('onTokenTransfer(address,uint256,bytes)')
      const header = '000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef' + // to
        '0000000000000000000000000000000000000000000000000000000000000539' + // amount
        '0000000000000000000000000000000000000000000000000000000000000060' + // offset
        '0000000000000000000000000000000000000000000000000000000000000136' //   length

      const maliciousPayload = ottSelector + header + requestPayload.slice(2)

      await assertActionThrows(async () => {
        await link.transferAndCall(oc.address, 0, maliciousPayload)
      })
    })
  })

  describe('#requestData', () => {
    context('when called through the LINK token', () => {
      const paid = 100
      let log, tx

      beforeEach(async () => {
        let args = requestDataBytes(to, specId, fHash, 'id', '')
        tx = await requestDataFrom(oc, link, paid, args)
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2]
      })

      it('logs an event', async () => {
        assert.equal(oc.address, log.address)

        assert.equal(specId, web3.toUtf8(log.topics[1]))
        assert.equal(defaultAccount, web3.toDecimal(log.topics[2]))
        assert.equal(paid, web3.toDecimal(log.topics[3]))
      })

      it('uses the expected event signature', async () => {
        // If updating this test, be sure to update services.RunLogTopic.
        let eventSignature = '0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529'
        assert.equal(eventSignature, log.topics[0])
      })

      it('does not allow the same requestId to be used twice', async () => {
        let args2 = requestDataBytes(to, specId, fHash, 'id', '')
        await assertActionThrows(async () => {
          await requestDataFrom(oc, link, paid, args2)
        })
      })
    })

    context('when not called through the LINK token', () => {
      it('reverts', async () => {
        await assertActionThrows(async () => {
          await oc.requestData(0, 0, 1, specId, fHash, 'id', '', { from: oracleNode })
        })
      })
    })
  })

  describe('#fulfillData', () => {
    let mock, internalId
    let currency = 'USD'

    context('cooperative consumer', () => {
      beforeEach(async () => {
        mock = await deploy('examples/BasicConsumer.sol', link.address, oc.address, toHex(specId))
        await link.transfer(mock.address, web3.toWei('1', 'ether'))
        await mock.requestEthereumPrice(currency)
        let event = await getLatestEvent(oc)
        internalId = event.args.internalId
      })

      context('when called by an unauthorized node', () => {
        it('raises an error', async () => {
          let unauthorized = await oc.getAuthorizationStatus(stranger)
          assert.equal(false, unauthorized)
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'Hello World!', { from: stranger })
          })
        })
      })

      context('when called by an authorized node', () => {
        it('raises an error if the request ID does not exist', async () => {
          await assertActionThrows(async () => {
            await oc.fulfillData(0xdeadbeef, 'Hello World!', { from: oracleNode })
          })
        })

        it('sets the value on the requested contract', async () => {
          await oc.fulfillData(internalId, 'Hello World!', { from: oracleNode })

          let currentPrice = await mock.currentPrice.call()
          assert.equal(web3.toUtf8(currentPrice), 'Hello World!')
        })

        it('does not allow a request to be fulfilled twice', async () => {
          await oc.fulfillData(internalId, 'First message!', { from: oracleNode })
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'Second message!!', { from: oracleNode })
          })
        })
      })
    })

    context('with a malicious requester', () => {
      const paymentAmount = toWei(1)

      it('cannot cancel before the expiration', async () => {
        mock = await deploy('examples/MaliciousRequester.sol', link.address, oc.address)
        await link.transfer(mock.address, paymentAmount)

        await assertActionThrows(async () => {
          await mock.maliciousRequestCancel()
        })
      })
    })

    context('with a malicious consumer', () => {
      const paymentAmount = toWei(1)

      beforeEach(async () => {
        mock = await deploy('examples/MaliciousConsumer.sol', link.address, oc.address)
        await link.transfer(mock.address, paymentAmount)
      })

      context('fails during fulfillment', () => {
        beforeEach(async () => {
          const req = await mock.requestData('assertFail(bytes32,bytes32)')
          internalId = runRequestId(req.receipt.logs[3])
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })

          const balance = await link.balanceOf.call(oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(oracleNode, paymentAmount, { from: defaultAccount })
          const newBalance = await link.balanceOf.call(oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'hack the planet 102', { from: oracleNode })
          })
        })
      })

      context('calls selfdestruct', () => {
        beforeEach(async () => {
          const req = await mock.requestData('doesNothing(bytes32,bytes32)')
          internalId = runRequestId(req.receipt.logs[3])
          await mock.remove()
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })

          const balance = await link.balanceOf.call(oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(oracleNode, paymentAmount, { from: defaultAccount })
          const newBalance = await link.balanceOf.call(oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })
      })

      context('request is canceled during fulfillment', () => {
        beforeEach(async () => {
          const req = await mock.requestData('cancelRequestOnFulfill(bytes32,bytes32)')
          internalId = runRequestId(req.receipt.logs[3])

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))

          const balance = await link.balanceOf.call(oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(oracleNode, paymentAmount, { from: defaultAccount })
          const newBalance = await link.balanceOf.call(oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'hack the planet 102', { from: oracleNode })
          })
        })
      })

      context('requester lies about amount of LINK sent', () => {
        it('the oracle uses the amount of LINK actually paid', async () => {
          const req = await mock.requestData('assertFail(bytes32,bytes32)')
          const log = req.receipt.logs[3]

          assert.equal(web3.toWei(1), web3.toDecimal(log.topics[3]))
        })
      })

      context('tries to steal funds from node', () => {
        it('is not successful with call', async () => {
          const req = await mock.requestData('stealEthCall(bytes32,bytes32)')
          internalId = runRequestId(req.receipt.logs[3])

          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })
          const mockBalance = web3.fromWei(web3.eth.getBalance(mock.address))
          assert.equal(mockBalance, 0)
        })

        it('is not successful with send', async () => {
          const req = await mock.requestData('stealEthSend(bytes32,bytes32)')
          internalId = runRequestId(req.receipt.logs[3])

          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })
          const mockBalance = web3.fromWei(web3.eth.getBalance(mock.address))
          assert.equal(mockBalance, 0)
        })

        it('is not successful with transfer', async () => {
          const req = await mock.requestData('stealEthTransfer(bytes32,bytes32)')
          internalId = runRequestId(req.receipt.logs[3])

          await oc.fulfillData(internalId, 'hack the planet 101', { from: oracleNode })
          const mockBalance = web3.fromWei(web3.eth.getBalance(mock.address))
          assert.equal(mockBalance, 0)
        })
      })
    })
  })

  describe('#withdraw', () => {
    context('without reserving funds via requestData', () => {
      it('does nothing', async () => {
        let balance = await link.balanceOf(oracleNode)
        assert.equal(0, balance)
        await assertActionThrows(async () => {
          await oc.withdraw(oracleNode, toWei(1), { from: defaultAccount })
        })
        balance = await link.balanceOf(oracleNode)
        assert.equal(0, balance)
      })
    })

    context('reserving funds via requestData', () => {
      let internalId, amount

      beforeEach(async () => {
        amount = 15
        const mock = await deploy('examples/GetterSetter.sol')
        const args = requestDataBytes(to, specId, fHash, 'id', '')
        const tx = await requestDataFrom(oc, link, amount, args)
        assert.equal(3, tx.receipt.logs.length)
        internalId = runRequestId(tx.receipt.logs[2])
      })

      context('but not freeing funds w fulfillData', () => {
        it('does not transfer funds', async () => {
          await assertActionThrows(async () => {
            await oc.withdraw(oracleNode, amount, { from: defaultAccount })
          })
          let balance = await link.balanceOf(oracleNode)
          assert.equal(0, balance)
        })
      })

      context('and freeing funds', () => {
        beforeEach(async () => {
          await oc.fulfillData(internalId, 'Hello World!', { from: oracleNode })
        })

        it('does not allow input greater than the balance', async () => {
          let originalOracleBalance = await link.balanceOf(oc.address)
          let originalStrangerBalance = await link.balanceOf(stranger)
          let withdrawAmount = amount + 1

          assert.isAbove(withdrawAmount, originalOracleBalance.toNumber())
          await assertActionThrows(async () => {
            await oc.withdraw(stranger, withdrawAmount, { from: defaultAccount })
          })

          let newOracleBalance = await link.balanceOf(oc.address)
          let newStrangerBalance = await link.balanceOf(stranger)

          assert.equal(originalOracleBalance.toNumber(), newOracleBalance.toNumber())
          assert.equal(originalStrangerBalance.toNumber(), newStrangerBalance.toNumber())
        })

        it('allows transfer of partial balance by owner to specified address', async () => {
          let partialAmount = 6
          let difference = amount - partialAmount
          await oc.withdraw(stranger, partialAmount, { from: defaultAccount })
          let strangerBalance = await link.balanceOf(stranger)
          let oracleBalance = await link.balanceOf(oc.address)
          assert.equal(partialAmount, strangerBalance)
          assert.equal(difference, oracleBalance)
        })

        it('allows transfer of entire balance by owner to specified address', async () => {
          await oc.withdraw(stranger, amount, { from: defaultAccount })
          let balance = await link.balanceOf(stranger)
          assert.equal(amount, balance)
        })

        it('does not allow a transfer of funds by non-owner', async () => {
          await assertActionThrows(async () => {
            await oc.withdraw(stranger, amount, { from: stranger })
          })
          let balance = await link.balanceOf(stranger)
          assert.equal(0, balance)
        })
      })
    })
  })

  describe('#cancel', () => {
    context('with no pending requests', () => {
      it('fails', async () => {
        await increaseTime5Minutes()
        await assertActionThrows(async () => {
          await oc.cancel(1337, { from: stranger })
        })
      })
    })

    context('with a pending request', () => {
      let internalId, tx, mock, requestAmount, startingBalance
      assert(mock === undefined, 'silence linter')
      let requestId = 'requestId'
      beforeEach(async () => {
        startingBalance = 100
        requestAmount = 20

        mock = await deploy('examples/GetterSetter.sol')
        await link.transfer(consumer, startingBalance)

        let args = requestDataBytes(to, specId, fHash, requestId, '')
        tx = await link.transferAndCall(oc.address, requestAmount, args, { from: consumer })
        assert.equal(3, tx.receipt.logs.length)
        internalId = runRequestId(tx.receipt.logs[2])
      })

      it('has correct initial balances', async () => {
        let oracleBalance = await link.balanceOf(oc.address)
        assert.equal(requestAmount, oracleBalance)

        let consumerAmount = await link.balanceOf(consumer)
        assert.equal(startingBalance - requestAmount, consumerAmount)
      })

      context('from a stranger', () => {
        it('fails', async () => {
          await assertActionThrows(async () => {
            await oc.cancel(requestId, { from: stranger })
          })
        })
      })

      context('from the requester', () => {
        it('refunds the correct amount', async () => {
          await increaseTime5Minutes()
          await oc.cancel(requestId, { from: consumer })
          let balance = await link.balanceOf(consumer)
          assert.equal(startingBalance, balance) // 100
        })

        it('triggers a cancellation event', async () => {
          await increaseTime5Minutes()
          const tx = await oc.cancel(requestId, { from: consumer })

          assert.equal(tx.receipt.logs.length, 2)
          assert.equal(internalId, tx.receipt.logs[1].data)
        })

        context('canceling twice', () => {
          it('fails', async () => {
            await increaseTime5Minutes()
            await oc.cancel(requestId, { from: consumer })
            await assertActionThrows(async () => {
              await oc.cancel(requestId, { from: consumer })
            })
          })
        })
      })
    })
  })
})
