import * as h from './support/helpers'

contract('Oracle', () => {
  const sourcePath = 'Oracle.sol'
  const fHash = h.functionSelector('requestedBytes32(bytes32,bytes32)')
  const specId = '4c7b7ffb66b344fbaa64995af81e355a'
  const to = '0x80e29acb842498fe6591f020bd82766dce619d43'
  let link, oc

  beforeEach(async () => {
    link = await h.deploy('link_token/contracts/LinkToken.sol')
    oc = await h.deploy(sourcePath, link.address)
    await oc.transferOwnership(h.oracleNode, {from: h.defaultAccount})
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(sourcePath), [
      'cancel',
      'fulfillData',
      'onTokenTransfer',
      'owner',
      'renounceOwnership',
      'requestData',
      'transferOwnership',
      'withdraw'
    ])
  })

  describe('#transferOwnership', () => {
    context('when called by the owner', () => {
      beforeEach(async () => {
        await oc.transferOwnership(h.stranger, {from: h.oracleNode})
      })

      it('can change the owner', async () => {
        let owner = await oc.owner.call()
        assert.isTrue(web3.isAddress(owner))
        assert.equal(h.stranger, owner)
      })
    })

    context('when called by a non-owner', () => {
      it('cannot change the owner', async () => {
        await h.assertActionThrows(async () => {
          await oc.transferOwnership(h.stranger, {from: h.stranger})
        })
      })
    })
  })

  describe('#onTokenTransfer', () => {
    context('when called from any address but the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = h.requestDataBytes(specId, to, fHash, 'id', '')

        await h.assertActionThrows(async () => {
          let tx = await oc.onTokenTransfer(h.oracleNode, 0, callData)
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
      let mock
      const paymentAmount = 1

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
    })

    it('does not allow recursive calls of onTokenTransfer', async () => {
      const requestPayload = h.requestDataBytes(specId, to, fHash, 'id', '')

      const ottSelector = h.functionSelector('onTokenTransfer(address,uint256,bytes)')
      const header = '000000000000000000000000c5fdf4076b8f3a5357c5e395ab970b5b54098fef' + // to
        '0000000000000000000000000000000000000000000000000000000000000539' + // amount
        '0000000000000000000000000000000000000000000000000000000000000060' + // offset
        '0000000000000000000000000000000000000000000000000000000000000136'   // length

      const maliciousPayload = ottSelector + header + requestPayload.slice(2)

      await h.assertActionThrows(async () => {
        const tx = await link.transferAndCall(oc.address, 0, maliciousPayload)
      })
    })
  })

  describe('#requestData', () => {
    context('when called through the LINK token', () => {
      const paid = 100
      let log, tx

      beforeEach(async () => {
        let args = h.requestDataBytes(specId, to, fHash, 'id', '')
        tx = await h.requestDataFrom(oc, link, paid, args)
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2]
      })

      it('logs an event', async () => {
        assert.equal(specId, web3.toUtf8(log.topics[1]))
        assert.equal(h.defaultAccount, web3.toDecimal(log.topics[2]))
        assert.equal(paid, web3.toDecimal(log.topics[3]))
      })

      it('uses the expected event signature', async () => {
        // If updating this test, be sure to update services.RunLogTopic.
        let eventSignature = '0x6d6db1f8fe19d95b1d0fa6a4bce7bb24fbf84597b35a33ff95521fac453c1529'
        assert.equal(eventSignature, log.topics[0])
      })
    })

    context('when not called through the LINK token', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await oc.requestData(0, 0, 1, specId, to, fHash, 'id', '', {from: h.oracleNode})
        })
      })
    })
  })

  describe('#fulfillData', () => {
    let mock, internalId
    let requestId = 'XID'

    context('cooperative consumer', () => {
      beforeEach(async () => {
        mock = await h.deploy('examples/GetterSetter.sol')
        let fHash = h.functionSelector('requestedBytes32(bytes32,bytes32)')
        let args = h.requestDataBytes(specId, mock.address, fHash, requestId, '')
        let req = await h.requestDataFrom(oc, link, 0, args)
        internalId = h.runRequestId(req.receipt.logs[2])
      })

      context('when called by a non-owner', () => {
        it('raises an error', async () => {
          await h.assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'Hello World!', {from: h.stranger})
          })
        })
      })

      context('when called by an owner', () => {
        it('raises an error if the request ID does not exist', async () => {
          await h.assertActionThrows(async () => {
            await oc.fulfillData(0xdeadbeef, 'Hello World!', {from: h.oracleNode})
          })
        })

        it('sets the value on the requested contract', async () => {
          await oc.fulfillData(internalId, 'Hello World!', {from: h.oracleNode})

          let mockRequestId = await mock.requestId.call()
          assert.equal(requestId.toString(), web3.toUtf8(mockRequestId))

          let currentValue = await mock.getBytes32.call()
          assert.equal('Hello World!', web3.toUtf8(currentValue))
        })

        it('does not allow a request to be fulfilled twice', async () => {
          await oc.fulfillData(internalId, 'First message!', {from: h.oracleNode})
          await h.assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'Second message!!', {from: h.oracleNode})
          })
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
          await mock.maliciousRequestCancel()
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
          const req = await mock.requestData('assertFail(bytes32,bytes32)')
          internalId = h.runRequestId(req.receipt.logs[3])
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: h.oracleNode})

          const balance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(h.oracleNode, paymentAmount, {from: h.oracleNode})
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: h.oracleNode})
          await h.assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'hack the planet 102', {from: h.oracleNode})
          })
        })
      })

      context('calls selfdestruct', () => {
        beforeEach(async () => {
          const req = await mock.requestData('doesNothing(bytes32,bytes32)')
          internalId = h.runRequestId(req.receipt.logs[3])
          await mock.remove()
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: h.oracleNode})

          const balance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(h.oracleNode, paymentAmount, {from: h.oracleNode})
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })
      })

      context('request is canceled during fulfillment', () => {
        beforeEach(async () => {
          const req = await mock.requestData('cancelRequestOnFulfill(bytes32,bytes32)')
          internalId = h.runRequestId(req.receipt.logs[3])

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: h.oracleNode})

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))

          const balance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(h.oracleNode, paymentAmount, {from: h.oracleNode})
          const newBalance = await link.balanceOf.call(h.oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: h.oracleNode})
          await h.assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'hack the planet 102', {from: h.oracleNode})
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
    })
  })

  describe('#withdraw', () => {
    context('without reserving funds via requestData', () => {
      it('does nothing', async () => {
        let balance = await link.balanceOf(h.oracleNode)
        assert.equal(0, balance)
        await h.assertActionThrows(async () => {
          await oc.withdraw(h.oracleNode, h.toWei(1), {from: h.oracleNode})
        })
        balance = await link.balanceOf(h.oracleNode)
        assert.equal(0, balance)
      })
    })

    context('reserving funds via requestData', () => {
      let internalId, amount

      beforeEach(async () => {
        amount = 15
        const mock = await h.deploy('examples/GetterSetter.sol')
        const args = h.requestDataBytes(specId, mock.address, fHash, 'id', '')
        const tx = await h.requestDataFrom(oc, link, amount, args)
        assert.equal(3, tx.receipt.logs.length)
        internalId = h.runRequestId(tx.receipt.logs[2])
      })

      context('but not freeing funds w fulfillData', () => {
        it('does not transfer funds', async () => {
          await h.assertActionThrows(async () => {
            await oc.withdraw(h.oracleNode, amount, {from: h.oracleNode})
          })
          let balance = await link.balanceOf(h.oracleNode)
          assert.equal(0, balance)
        })
      })

      context('and freeing funds', () => {
        beforeEach(async () => {
          await oc.fulfillData(internalId, 'Hello World!', {from: h.oracleNode})
        })

        it('does not allow input greater than the balance', async () => {
          let originalOracleBalance = await link.balanceOf(oc.address)
          let originalStrangerBalance = await link.balanceOf(h.stranger)
          let withdrawAmount = amount + 1

          assert.isAbove(withdrawAmount, originalOracleBalance.toNumber())
          await h.assertActionThrows(async () => {
            await oc.withdraw(h.stranger, withdrawAmount, {from: h.oracleNode})
          })

          let newOracleBalance = await link.balanceOf(oc.address)
          let newStrangerBalance = await link.balanceOf(h.stranger)

          assert.equal(originalOracleBalance.toNumber(), newOracleBalance.toNumber())
          assert.equal(originalStrangerBalance.toNumber(), newStrangerBalance.toNumber())
        })

        it('allows transfer of partial balance by owner to specified address', async () => {
          let partialAmount = 6
          let difference = amount - partialAmount
          await oc.withdraw(h.stranger, partialAmount, {from: h.oracleNode})
          let strangerBalance = await link.balanceOf(h.stranger)
          let oracleBalance = await link.balanceOf(oc.address)
          assert.equal(partialAmount, strangerBalance)
          assert.equal(difference, oracleBalance)
        })

        it('allows transfer of entire balance by owner to specified address', async () => {
          await oc.withdraw(h.stranger, amount, {from: h.oracleNode})
          let balance = await link.balanceOf(h.stranger)
          assert.equal(amount, balance)
        })

        it('does not allow a transfer of funds by non-owner', async () => {
          await h.assertActionThrows(async () => {
            await oc.withdraw(h.stranger, amount, {from: h.stranger})
          })
          let balance = await link.balanceOf(h.stranger)
          assert.equal(0, balance)
        })
      })
    })
  })

  describe('#cancel', () => {
    context('with no pending requests', () => {
      it('fails', async () => {
        await h.increaseTime5Minutes();
        await h.assertActionThrows(async () => {
          await oc.cancel(1337, {from: h.stranger})
        })
      })
    })

    context('with a pending request', () => {
      let internalId, tx, mock, requestAmount, startingBalance
      let requestId = 'requestId'
      beforeEach(async () => {
        startingBalance = 100
        requestAmount = 20

        mock = await h.deploy('examples/GetterSetter.sol')
        await link.transfer(h.consumer, startingBalance)

        let args = h.requestDataBytes(specId, h.consumer, fHash, requestId, '')
        tx = await link.transferAndCall(oc.address, requestAmount, args, {from: h.consumer})
        assert.equal(3, tx.receipt.logs.length)
        internalId = h.runRequestId(tx.receipt.logs[2])
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
            await oc.cancel(requestId, {from: h.stranger})
          })
        })
      })

      context('from the requester', () => {
        it('refunds the correct amount', async () => {
          await h.increaseTime5Minutes();
          await oc.cancel(requestId, {from: h.consumer})
          let balance = await link.balanceOf(h.consumer)
          assert.equal(startingBalance, balance) // 100
        })

        it('triggers a cancellation event', async () => {
          await h.increaseTime5Minutes();
          const tx = await oc.cancel(requestId, {from: h.consumer})

          assert.equal(tx.receipt.logs.length, 2)
          assert.equal(internalId, tx.receipt.logs[1].data)
        })

        context('canceling twice', () => {
          it('fails', async () => {
            await h.increaseTime5Minutes();
            await oc.cancel(requestId, {from: h.consumer})
            await h.assertActionThrows(async () => {
              await oc.cancel(requestId, {from: h.consumer})
            })
          })
        })
      })
    })
  })
})
