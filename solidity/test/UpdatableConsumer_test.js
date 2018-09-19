import {
  assertActionThrows,
  defaultAccount,
  deploy,
  getLatestEvent,
  oracleNode,
} from './support/helpers'
import namehash from 'eth-ens-namehash'

contract('UpdatableConsumer', () => {
  const sourcePath = 'examples/UpdatableConsumer.sol'

  const ensRoot = namehash.hash()
  const tld = 'test'
  const tldSubnode = namehash.hash(tld)
  const tokenDomain = 'link'
  const tokenHash = namehash.hash(`${tokenDomain}.${tld}`)
  const oracleDomain = 'oracle'
  const oracleHash = namehash.hash(`${oracleDomain}.${tld}`)
  const specId = web3.sha3('someSpecID')
  const newOracleAddress = '0xf000000000000000000000000000000000000ba7'
  const currency = 'USD'

  let ens, ensResolver, link, oc, uc

  beforeEach(async () => {
    link = await deploy('LinkToken.sol')
    oc = await deploy('Oracle.sol', link.address)
    await oc.transferOwnership(oracleNode, {from: defaultAccount})
    ens = await deploy('ENSRegistry.sol')
    ensResolver = await deploy('PublicResolver.sol', ens.address)

    // register domain
    await ens.setSubnodeOwner(ensRoot, web3.sha3(tld), oracleNode)
    await ens.setResolver(tldSubnode, ensResolver.address, {from: oracleNode})
    await ensResolver.setAddr(tldSubnode, oc.address, {from: oracleNode})

     // register token subdomain
    await ens.setSubnodeOwner(tldSubnode, web3.sha3(tokenDomain), oracleNode, {from: oracleNode})
    await ensResolver.setAddr(tokenHash, link.address, {from: oracleNode})

    // register oracle subdomain
    await ens.setSubnodeOwner(tldSubnode, web3.sha3(oracleDomain), oracleNode, {from: oracleNode})
    await ensResolver.setAddr(oracleHash, oc.address, {from: oracleNode})

    uc = await deploy(sourcePath, specId, ens.address, tldSubnode)
  })

  describe('constructor', () => {
    it('pulls the token contract address from the resolver', async () => {
      assert.equal(link.address, await uc.publicLinkToken.call())
    })

    it('pulls the oracle contract address from the resolver', async () => {
      assert.equal(oc.address, await uc.publicOracle.call())
    })
  })

  describe('#updateOracle', () => {
    describe('when the ENS resolver has been updated', () => {
      beforeEach(async () => {
        await ensResolver.setAddr(oracleHash, newOracleAddress, {from: oracleNode})
      })

      it("updates the contract's oracle address", async () => {
        await uc.updateOracle()

        assert.equal(newOracleAddress, await uc.publicOracle.call())
      })
    })

    describe("when the ENS resolver has not been updated", () => {
      it("keeps the same oracle address", async () => {
        await uc.updateOracle()

        assert.equal(oc.address, await uc.publicOracle.call())
      })
    })
  })

  describe('#fulfillData', () => {
    const response = '1,000,000.00'
    let internalId, requestId

    beforeEach(async () => {
      await link.transfer(uc.address, web3.toWei('1', 'ether'))
      await uc.requestEthereumPrice(currency)
      const event = await getLatestEvent(oc)
      internalId = event.args.internalId

      const event2 = await getLatestEvent(uc)
      requestId = event2.args.id
    })

    it('records the data given to it by the oracle', async () => {
      await oc.fulfillData(internalId, response, {from: oracleNode})

      const currentPrice = await uc.currentPrice.call()
      assert.equal(web3.toUtf8(currentPrice), response)
    })

    context('when the oracle address is updated before a request is fulfilled', () => {
      beforeEach(async () => {
        await ensResolver.setAddr(oracleHash, newOracleAddress, {from: oracleNode})
        await uc.updateOracle()
        assert.equal(newOracleAddress, await uc.publicOracle.call())
      })

      it('records the data given to it by the old oracle contract', async () => {
        await oc.fulfillData(internalId, response, {from: oracleNode})

        const currentPrice = await uc.currentPrice.call()
        assert.equal(web3.toUtf8(currentPrice), response)
      })

      it('does not accept responses from the new oracle for the old requests', async () => {
        await assertActionThrows(async () => {
          await uc.fulfill(requestId, response, {from: oracleNode})
        })

        const currentPrice = await uc.currentPrice.call()
        assert.equal(web3.toUtf8(currentPrice), '')
      })

      it('does not accept responses from the new oracle for the old requests', async () => {
        await assertActionThrows(async () => {
          await uc.fulfill(requestId, response, {from: oracleNode})
        })

        const currentPrice = await uc.currentPrice.call()
        assert.equal(web3.toUtf8(currentPrice), '')
      })
    })
  })
})
