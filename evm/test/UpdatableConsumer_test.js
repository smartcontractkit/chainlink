import * as h from './support/helpers'
import namehash from 'eth-ens-namehash'
import { assertBigNum } from './support/matchers'
const ENSRegistry = artifacts.require('ENSRegistry.sol')
const Oracle = artifacts.require('Oracle.sol')
const PublicResolver = artifacts.require('PublicResolver.sol')
const UpdatableConsumer = artifacts.require('UpdatableConsumer.sol')

contract('UpdatableConsumer', () => {
  const ensRoot = namehash.hash()
  const tld = 'cltest'
  const tldSubnode = namehash.hash(tld)
  const domain = 'chainlink'
  const domainNode = namehash.hash(`${domain}.${tld}`)
  const tokenSubdomain = 'link'
  const tokenSubnode = namehash.hash(`${tokenSubdomain}.${domain}.${tld}`)
  const oracleSubdomain = 'oracle'
  const oracleSubnode = namehash.hash(`${oracleSubdomain}.${domain}.${tld}`)
  const specId = h.keccak('someSpecID')
  const newOracleAddress = '0xf000000000000000000000000000000000000ba7'

  let ens, ensResolver, link, oc, uc

  beforeEach(async () => {
    link = await h.linkContract()
    oc = await Oracle.new(link.address, { from: h.oracleNode })
    ens = await ENSRegistry.new()
    ensResolver = await PublicResolver.new(ens.address)

    // register tld
    await ens.setSubnodeOwner(ensRoot, h.keccak(tld), h.defaultAccount)
    // register domain
    await ens.setSubnodeOwner(tldSubnode, h.keccak(domain), h.oracleNode)
    await ens.setResolver(domainNode, ensResolver.address, {
      from: h.oracleNode
    })
    // register token subdomain to point to token contract
    await ens.setSubnodeOwner(
      domainNode,
      h.keccak(tokenSubdomain),
      h.oracleNode,
      { from: h.oracleNode }
    )
    await ens.setResolver(tokenSubnode, ensResolver.address, {
      from: h.oracleNode
    })
    await ensResolver.setAddr(tokenSubnode, link.address, {
      from: h.oracleNode
    })
    // register oracle subdomain to point to oracle contract
    await ens.setSubnodeOwner(
      domainNode,
      h.keccak(oracleSubdomain),
      h.oracleNode,
      { from: h.oracleNode }
    )
    await ens.setResolver(oracleSubnode, ensResolver.address, {
      from: h.oracleNode
    })
    await ensResolver.setAddr(oracleSubnode, oc.address, { from: h.oracleNode })

    // deploy updatable consumer contract
    uc = await UpdatableConsumer.new(specId, ens.address, domainNode)
  })

  describe('constructor', () => {
    it('pulls the token contract address from the resolver', async () => {
      assert.equal(link.address, await uc.getChainlinkToken.call())
    })

    it('pulls the oracle contract address from the resolver', async () => {
      assert.equal(oc.address, await uc.getOracle.call())
    })
  })

  describe('#updateOracle', () => {
    describe('when the ENS resolver has been updated', () => {
      beforeEach(async () => {
        await ensResolver.setAddr(oracleSubnode, newOracleAddress, {
          from: h.oracleNode
        })
      })

      it("updates the contract's oracle address", async () => {
        await uc.updateOracle()
        assert.equal(
          newOracleAddress.toLowerCase(),
          (await uc.getOracle.call()).toLowerCase()
        )
      })
    })

    describe('when the ENS resolver has not been updated', () => {
      it('keeps the same oracle address', async () => {
        await uc.updateOracle()

        assert.equal(oc.address, await uc.getOracle.call())
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = '1,000,000.00'
    const currency = 'USD'
    const paymentAmount = h.toWei(1, 'ether')
    let request

    beforeEach(async () => {
      await link.transfer(uc.address, paymentAmount)
      const tx = await uc.requestEthereumPrice(h.toHex(currency))
      request = h.decodeRunRequest(tx.receipt.rawLogs[3])
    })

    it('records the data given to it by the oracle', async () => {
      await h.fulfillOracleRequest(oc, request, response, {
        from: h.oracleNode
      })

      const currentPrice = await uc.currentPrice.call()
      assert.equal(h.toUtf8(currentPrice), response)
    })

    context(
      'when the oracle address is updated before a request is fulfilled',
      () => {
        beforeEach(async () => {
          await ensResolver.setAddr(oracleSubnode, newOracleAddress, {
            from: h.oracleNode
          })
          await uc.updateOracle()
          assert.equal(
            newOracleAddress.toLowerCase(),
            (await uc.getOracle.call()).toLowerCase()
          )
        })

        it('records the data given to it by the old oracle contract', async () => {
          await h.fulfillOracleRequest(oc, request, response, {
            from: h.oracleNode
          })

          const currentPrice = await uc.currentPrice.call()
          assert.equal(h.toUtf8(currentPrice), response)
        })

        it('does not accept responses from the new oracle for the old requests', async () => {
          await h.assertActionThrows(async () => {
            await uc.fulfill(request.id, h.toHex(response), {
              from: h.oracleNode
            })
          })

          const currentPrice = await uc.currentPrice.call()
          assert.equal(h.toUtf8(currentPrice), '')
        })

        it('still allows funds to be withdrawn from the oracle', async () => {
          await h.increaseTime5Minutes()
          assertBigNum(
            0,
            await link.balanceOf.call(uc.address),
            'Initial balance should be 0'
          )

          await uc.cancelRequest(
            request.id,
            request.payment,
            request.callbackFunc,
            request.expiration
          )

          assertBigNum(
            paymentAmount,
            await link.balanceOf.call(uc.address),
            'Oracle should have been repaid on cancellation.'
          )
        })
      }
    )
  })
})
