'use strict'

require('./support/helpers.js')

contract('Oracle', () => {
  const sourcePath = 'Oracle.sol'
  const fHash = functionSelector('requestedBytes32(bytes32,bytes32)')
  const specId = '4c7b7ffb66b344fbaa64995af81e355a'
  const to = '0x80e29acb842498fe6591f020bd82766dce619d43'
  let link, oc

  beforeEach(async () => {
    link = await deploy('linkToken/contracts/LinkToken.sol')
    oc = await deploy(sourcePath, link.address)
    await oc.transferOwnership(oracleNode, {from: defaultAccount})
  })

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [
      'cancel',
      'fulfillData',
      'onTokenTransfer',
      'owner',
      'requestData',
      'transferOwnership',
      'withdraw'
    ])
  })

  describe('#transferOwnership', () => {
    context('when called by the owner', () => {
      beforeEach(async () => {
        await oc.transferOwnership(stranger, {from: oracleNode})
      })

      it('can change the owner', async () => {
        let owner = await oc.owner.call()
        assert.isTrue(web3.isAddress(owner))
        assert.equal(stranger, owner)
      })
    })

    context('when called by a non-owner', () => {
      it('cannot change the owner', async () => {
        await assertActionThrows(async () => {
          await oc.transferOwnership(stranger, {from: stranger})
        })
      })
    })
  })

  describe('#onTokenTransfer', () => {
    let mock

    context('when called from the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = requestDataBytes(specId, to, fHash, 'id', '')

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

    context('when called from any address but the LINK token', () => {
      it('triggers the intended method', async () => {
        let callData = requestDataBytes(specId, to, fHash, 'id', '')

        await assertActionThrows(async () => {
          let tx = await oc.onTokenTransfer(oracleNode, 0, callData)
        })
      })
    })

    context('malicious requester', () => {
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
  })

  describe('#requestData', () => {
    context('when called through the LINK token', () => {
      let log, tx
      beforeEach(async () => {
        let args = requestDataBytes(specId, to, fHash, 'id', '')
        tx = await requestDataFrom(oc, link, 0, args)
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2]
      })

      it('logs an event', async () => {
        assert.equal(specId, web3.toUtf8(log.topics[2]))
      })

      it('uses the expected event signature', async () => {
        // If updating this test, be sure to update services.RunLogTopic.
        let eventSignature = '0x3fab86a1207bdcfe3976d0d9df25f263d45ae8d381a60960559771a2b223974d'
        assert.equal(eventSignature, log.topics[0])
      })
    })

    context('when not called through the LINK token', () => {
      it('reverts', async () => {
        await assertActionThrows(async () => {
          await oc.requestData(1, specId, to, fHash, 'id', '', {from: oracleNode})
        })
      })
    })
  })

  describe('#fulfillData', () => {
    let mock, internalId
    let requestId = 'XID'

    context('successful consumer', () => {
      beforeEach(async () => {
        mock = await deploy('examples/GetterSetter.sol')
        let fHash = functionSelector('requestedBytes32(bytes32,bytes32)')
        let args = requestDataBytes(specId, mock.address, fHash, requestId, '')
        let req = await requestDataFrom(oc, link, 0, args)
        internalId = req.receipt.logs[2].topics[1]
      })

      context('when called by a non-owner', () => {
        it('raises an error', async () => {
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'Hello World!', {from: stranger})
          })
        })
      })

      context('when called by an owner', () => {
        it('raises an error if the request ID does not exist', async () => {
          await assertActionThrows(async () => {
            await oc.fulfillData(0xdeadbeef, 'Hello World!', {from: oracleNode})
          })
        })

        it('sets the value on the requested contract', async () => {
          await oc.fulfillData(internalId, 'Hello World!', {from: oracleNode})

          let mockRequestId = await mock.requestId.call()
          assert.equal(requestId.toString(), web3.toUtf8(mockRequestId))

          let currentValue = await mock.getBytes32.call()
          assert.equal('Hello World!', web3.toUtf8(currentValue))
        })

        it('does not allow a request to be fulfilled twice', async () => {
          await oc.fulfillData(internalId, 'First message!', {from: oracleNode})
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'Second message!!', {from: oracleNode})
          })
        })
      })
    })

    context('malicious consumer', () => {
      const paymentAmount = toWei(1)

      context('fails during fulfillment', () => {
        beforeEach(async () => {
          mock = await deploy('examples/MaliciousConsumer.sol', link.address, oc.address)
          await link.transfer(mock.address, paymentAmount)

          const req = await mock.requestData('assertFail(bytes32,bytes32)')
          internalId = req.receipt.logs[2].topics[1]
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: oracleNode})

          const balance = await link.balanceOf.call(oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(oracleNode, paymentAmount, {from: oracleNode})
          const newBalance = await link.balanceOf.call(oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: oracleNode})
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'hack the planet 102', {from: oracleNode})
          })
        })
      })

      context('calls selfdestruct', () => {
        beforeEach(async () => {
          mock = await deploy('examples/MaliciousConsumer.sol', link.address, oc.address)
          await link.transfer(mock.address, paymentAmount)

          const req = await mock.requestData('doesNothing(bytes32,bytes32)')
          internalId = req.receipt.logs[2].topics[1]
          await mock.remove()
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: oracleNode})

          const balance = await link.balanceOf.call(oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(oracleNode, paymentAmount, {from: oracleNode})
          const newBalance = await link.balanceOf.call(oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })
      })

      context('request is canceled during fulfillment', () => {
        beforeEach(async () => {
          mock = await deploy('examples/MaliciousConsumer.sol', link.address, oc.address)
          await link.transfer(mock.address, paymentAmount)

          const req = await mock.requestData('cancelRequestOnFulfill(bytes32,bytes32)')
          internalId = req.receipt.logs[2].topics[1]

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))
        })

        it('allows the oracle node to receive their payment', async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: oracleNode})

          const mockBalance = await link.balanceOf.call(mock.address)
          assert.isTrue(mockBalance.equals(0))

          const balance = await link.balanceOf.call(oracleNode)
          assert.isTrue(balance.equals(0))

          await oc.withdraw(oracleNode, paymentAmount, {from: oracleNode})
          const newBalance = await link.balanceOf.call(oracleNode)
          assert.isTrue(paymentAmount.equals(newBalance))
        })

        it("can't fulfill the data again", async () => {
          await oc.fulfillData(internalId, 'hack the planet 101', {from: oracleNode})
          await assertActionThrows(async () => {
            await oc.fulfillData(internalId, 'hack the planet 102', {from: oracleNode})
          })
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
          await oc.withdraw(oracleNode, toWei(1), {from: oracleNode})
        })
        balance = await link.balanceOf(oracleNode)
        assert.equal(0, balance)
      })
    })

    context('reserving funds via requestData', () => {
      let log, tx, mock, internalId, amount
      beforeEach(async () => {
        amount = 15
        mock = await deploy('examples/GetterSetter.sol')
        let args = requestDataBytes(specId, mock.address, fHash, 'id', '')
        tx = await requestDataFrom(oc, link, amount, args)
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2]
        internalId = log.topics[1]
      })

      context('but not freeing funds w fulfillData', () => {
        it('does not transfer funds', async () => {
          await assertActionThrows(async () => {
            await oc.withdraw(oracleNode, amount, {from: oracleNode})
          })
          let balance = await link.balanceOf(oracleNode)
          assert.equal(0, balance)
        })
      })

      context('and freeing funds', () => {
        beforeEach(async () => {
          await oc.fulfillData(internalId, 'Hello World!', {from: oracleNode})
        })

        it('does not allow input greater than the balance', async () => {
          let originalOracleBalance = await link.balanceOf(oc.address)
          let originalStrangerBalance = await link.balanceOf(stranger)
          let withdrawAmount = amount + 1

          assert.isAbove(withdrawAmount, originalOracleBalance.toNumber())
          await assertActionThrows(async () => {
            await oc.withdraw(stranger, withdrawAmount, {from: oracleNode})
          })

          let newOracleBalance = await link.balanceOf(oc.address)
          let newStrangerBalance = await link.balanceOf(stranger)

          assert.equal(originalOracleBalance.toNumber(), newOracleBalance.toNumber())
          assert.equal(originalStrangerBalance.toNumber(), newStrangerBalance.toNumber())
        })

        it('allows transfer of partial balance by owner to specified address', async () => {
          let partialAmount = 6
          let difference = amount - partialAmount
          await oc.withdraw(stranger, partialAmount, {from: oracleNode})
          let strangerBalance = await link.balanceOf(stranger)
          let oracleBalance = await link.balanceOf(oc.address)
          assert.equal(partialAmount, strangerBalance)
          assert.equal(difference, oracleBalance)
        })

        it('allows transfer of entire balance by owner to specified address', async () => {
          await oc.withdraw(stranger, amount, {from: oracleNode})
          let balance = await link.balanceOf(stranger)
          assert.equal(amount, balance)
        })

        it('does not allow a transfer of funds by non-owner', async () => {
          await assertActionThrows(async () => {
            await oc.withdraw(stranger, amount, {from: stranger})
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
        await assertActionThrows(async () => {
          await oc.cancel(1337, {from: stranger})
        })
      })
    })

    context('with a pending request', () => {
      let log, tx, mock, requestAmount, startingBalance
      let requestId = 'requestId'
      beforeEach(async () => {
        startingBalance = 100
        requestAmount = 20

        mock = await deploy('examples/GetterSetter.sol')
        await link.transfer(consumer, startingBalance)

        let args = requestDataBytes(specId, consumer, fHash, requestId, '')
        tx = await link.transferAndCall(oc.address, requestAmount, args, {from: consumer})
        assert.equal(3, tx.receipt.logs.length)
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
            await oc.cancel(requestId, {from: stranger})
          })
        })
      })

      context('from the requester', () => {
        it('refunds the correct amount', async () => {
          await oc.cancel(requestId, {from: consumer})
          let balance = await link.balanceOf(consumer)
          assert.equal(startingBalance, balance) // 100
        })

        context('canceling twice', () => {
          it('fails', async () => {
            await oc.cancel(requestId, {from: consumer})
            await assertActionThrows(async () => {
              await oc.cancel(requestId, {from: consumer})
            })
          })
        })
      })
    })
  })
})
