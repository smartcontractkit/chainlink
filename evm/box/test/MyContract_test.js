'use strict'

// eslint-disable-next-line @typescript-eslint/no-var-requires
const h = require('chainlink-test-helpers')

contract('MyContract', accounts => {
  const LinkToken = artifacts.require('LinkToken.sol')
  const Oracle = artifacts.require('Oracle.sol')
  const MyContract = artifacts.require('MyContract.sol')

  const defaultAccount = accounts[0]
  const oracleNode = accounts[1]
  const stranger = accounts[2]
  const consumer = accounts[3]

  // These parameters are used to validate the data was received
  // on the deployed oracle contract. The Job ID only represents
  // the type of data, but will not work on a public testnet.
  // For the latest JobIDs, visit our docs here:
  // https://docs.chain.link/docs/testnet-oracles
  const jobId = web3.utils.toHex('4c7b7ffb66b344fbaa64995af81e355a')
  const url =
    'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY'
  const path = 'USD'
  const times = 100

  // Represents 1 LINK for testnet requests
  const payment = web3.utils.toWei('1')

  let link, oc, cc

  beforeEach(async () => {
    link = await LinkToken.new()
    oc = await Oracle.new(link.address, { from: defaultAccount })
    cc = await MyContract.new(link.address, { from: consumer })
    await oc.setFulfillmentPermission(oracleNode, true, {
      from: defaultAccount,
    })
  })

  describe('#createRequest', () => {
    context('without LINK', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await cc.createRequestTo(
            oc.address,
            jobId,
            payment,
            url,
            path,
            times,
            { from: consumer },
          )
        })
      })
    })

    context('with LINK', () => {
      let request

      beforeEach(async () => {
        await link.transfer(cc.address, web3.utils.toWei('1', 'ether'))
      })

      context('sending a request to a specific oracle contract address', () => {
        it('triggers a log event in the new Oracle contract', async () => {
          const tx = await cc.createRequestTo(
            oc.address,
            jobId,
            payment,
            url,
            path,
            times,
            { from: consumer },
          )
          request = h.decodeRunRequest(tx.receipt.rawLogs[3])
          assert.equal(oc.address, tx.receipt.rawLogs[3].address)
          assert.equal(
            request.topic,
            web3.utils.keccak256(
              'OracleRequest(bytes32,address,bytes32,uint256,address,bytes4,uint256,uint256,bytes)',
            ),
          )
        })
      })
    })
  })

  describe('#fulfill', () => {
    const expected = 50000
    const response = web3.utils.toHex(expected)
    let request

    beforeEach(async () => {
      await link.transfer(cc.address, web3.utils.toWei('1', 'ether'))
      const tx = await cc.createRequestTo(
        oc.address,
        jobId,
        payment,
        url,
        path,
        times,
        { from: consumer },
      )
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
      await h.fulfillOracleRequest(oc, request, response, { from: oracleNode })
    })

    it('records the data given to it by the oracle', async () => {
      const currentPrice = await cc.data.call()
      assert.equal(
        web3.utils.toHex(currentPrice),
        web3.utils.padRight(expected, 64),
      )
    })

    context('when my contract does not recognize the request ID', () => {
      const otherId = web3.utils.toHex('otherId')

      beforeEach(async () => {
        request.id = otherId
      })

      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: oracleNode,
          })
        })
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await cc.fulfill(request.id, response, { from: stranger })
        })
      })
    })
  })

  describe('#cancelRequest', () => {
    let request

    beforeEach(async () => {
      await link.transfer(cc.address, web3.utils.toWei('1', 'ether'))
      const tx = await cc.createRequestTo(
        oc.address,
        jobId,
        payment,
        url,
        path,
        times,
        { from: consumer },
      )
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
    })

    context('before the expiration time', () => {
      it('cannot cancel a request', async () => {
        await h.assertActionThrows(async () => {
          await cc.cancelRequest(
            request.id,
            request.payment,
            request.callbackFunc,
            request.expiration,
            { from: consumer },
          )
        })
      })
    })

    context('after the expiration time', () => {
      beforeEach(async () => {
        await h.increaseTime5Minutes()
      })

      context('when called by a non-owner', () => {
        it('cannot cancel a request', async () => {
          await h.assertActionThrows(async () => {
            await cc.cancelRequest(
              request.id,
              request.payment,
              request.callbackFunc,
              request.expiration,
              { from: stranger },
            )
          })
        })
      })

      context('when called by an owner', () => {
        it('can cancel a request', async () => {
          await cc.cancelRequest(
            request.id,
            request.payment,
            request.callbackFunc,
            request.expiration,
            { from: consumer },
          )
        })
      })
    })
  })

  describe('#withdrawLink', () => {
    beforeEach(async () => {
      await link.transfer(cc.address, web3.utils.toWei('1', 'ether'))
    })

    context('when called by a non-owner', () => {
      it('cannot withdraw', async () => {
        await h.assertActionThrows(async () => {
          await cc.withdrawLink({ from: stranger })
        })
      })
    })

    context('when called by the owner', () => {
      it('transfers LINK to the owner', async () => {
        const beforeBalance = await link.balanceOf(consumer)
        assert.equal(beforeBalance, '0')
        await cc.withdrawLink({ from: consumer })
        const afterBalance = await link.balanceOf(consumer)
        assert.equal(afterBalance, web3.utils.toWei('1', 'ether'))
      })
    })
  })
})
