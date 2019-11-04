import * as h from '../src/helpersV2'
import { assertBigNum } from '../src/matchersV2'
import { ENSRegistryFactory } from '../src/generated/ENSRegistryFactory'
import { PublicResolverFactory } from '../src/generated/PublicResolverFactory'
import { UpdatableConsumerFactory } from '../src/generated/UpdatableConsumerFactory'
import { ethers } from 'ethers'
import { assert } from 'chai'
import { LinkTokenFactory } from '../src/generated/LinkTokenFactory'
import { OracleFactory } from '../src/generated/OracleFactory'
import { Instance } from '../src/contract'
import ganache from 'ganache-core'

const linkTokenFactory = new LinkTokenFactory()
const ensRegistryFactory = new ENSRegistryFactory()
const oracleFactory = new OracleFactory()
const publicResolverFacotory = new PublicResolverFactory()
const updatableConsumerFactory = new UpdatableConsumerFactory()

const provider = new ethers.providers.Web3Provider(ganache.provider() as any)

let roles: h.Roles

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)

  roles = rolesAndPersonas.roles
})

describe('UpdatableConsumer', () => {
  // https://github.com/ethers-io/ethers-ens/blob/master/src.ts/index.ts#L631
  const ensRoot = ethers.utils.namehash('')
  const tld = 'cltest'
  const tldSubnode = ethers.utils.namehash(tld)
  const domain = 'chainlink'
  const domainNode = ethers.utils.namehash(`${domain}.${tld}`)
  const tokenSubdomain = 'link'
  const tokenSubnode = ethers.utils.namehash(
    `${tokenSubdomain}.${domain}.${tld}`,
  )
  const oracleSubdomain = 'oracle'
  const oracleSubnode = ethers.utils.namehash(
    `${oracleSubdomain}.${domain}.${tld}`,
  )
  const specId = ethers.utils.formatBytes32String('someSpecID')
  const newOracleAddress = '0xf000000000000000000000000000000000000ba7'

  let ens: Instance<ENSRegistryFactory>
  let ensResolver: Instance<PublicResolverFactory>
  let link: Instance<LinkTokenFactory>
  let oc: Instance<OracleFactory>
  let uc: Instance<UpdatableConsumerFactory>
  const deployment = h.useSnapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    oc = await oracleFactory.connect(roles.oracleNode).deploy(link.address)
    ens = await ensRegistryFactory.connect(roles.defaultAccount).deploy()

    ensResolver = await publicResolverFacotory
      .connect(roles.defaultAccount)
      .deploy(ens.address)
    const ensOracleNode = ens.connect(roles.oracleNode)
    const ensResolverOracleNode = ensResolver.connect(roles.oracleNode)

    // register tld
    await ens.setSubnodeOwner(
      ensRoot,
      h.keccak(ethers.utils.toUtf8Bytes(tld)),
      roles.defaultAccount.address,
    )

    // register domain
    await ens.setSubnodeOwner(
      tldSubnode,
      h.keccak(ethers.utils.toUtf8Bytes(domain)),
      roles.oracleNode.address,
    )

    await ensOracleNode.setResolver(domainNode, ensResolver.address)

    // register token subdomain to point to token contract
    await ensOracleNode.setSubnodeOwner(
      domainNode,
      h.keccak(ethers.utils.toUtf8Bytes(tokenSubdomain)),
      roles.oracleNode.address,
    )
    await ensOracleNode.setResolver(tokenSubnode, ensResolver.address)
    await ensResolverOracleNode.setAddr(tokenSubnode, link.address)

    // register oracle subdomain to point to oracle contract
    await ensOracleNode.setSubnodeOwner(
      domainNode,
      h.keccak(ethers.utils.toUtf8Bytes(oracleSubdomain)),
      roles.oracleNode.address,
    )
    await ensOracleNode.setResolver(oracleSubnode, ensResolver.address)
    await ensResolverOracleNode.setAddr(oracleSubnode, oc.address)

    // deploy updatable consumer contract
    uc = await updatableConsumerFactory
      .connect(roles.defaultAccount)
      .deploy(specId, ens.address, domainNode)
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('constructor', () => {
    it('pulls the token contract address from the resolver', async () => {
      assert.equal(link.address, await uc.getChainlinkToken())
    })

    it('pulls the oracle contract address from the resolver', async () => {
      assert.equal(oc.address, await uc.getOracle())
    })
  })

  describe('#updateOracle', () => {
    describe('when the ENS resolver has been updated', () => {
      beforeEach(async () => {
        await ensResolver
          .connect(roles.oracleNode)
          .setAddr(oracleSubnode, newOracleAddress)
      })

      it("updates the contract's oracle address", async () => {
        await uc.updateOracle()
        assert.equal(
          newOracleAddress.toLowerCase(),
          (await uc.getOracle()).toLowerCase(),
        )
      })
    })

    describe('when the ENS resolver has not been updated', () => {
      it('keeps the same oracle address', async () => {
        await uc.updateOracle()

        assert.equal(oc.address, await uc.getOracle())
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = ethers.utils.formatBytes32String('1,000,000.00')
    const currency = 'USD'
    const paymentAmount = h.toWei('1')
    let request: h.RunRequest

    beforeEach(async () => {
      await link.transfer(uc.address, paymentAmount)
      const tx = await uc.requestEthereumPrice(
        h.toHex(ethers.utils.toUtf8Bytes(currency)),
      )
      const receipt = await tx.wait()
      request = h.decodeRunRequest(receipt.logs![3])
    })

    it('records the data given to it by the oracle', async () => {
      await h.fulfillOracleRequest(oc, request, response)

      const currentPrice = await uc.currentPrice()
      assert.equal(currentPrice, response)
    })

    describe('when the oracle address is updated before a request is fulfilled', () => {
      beforeEach(async () => {
        await ensResolver
          .connect(roles.oracleNode)
          .setAddr(oracleSubnode, newOracleAddress)
        await uc.updateOracle()
        assert.equal(
          newOracleAddress.toLowerCase(),
          (await uc.getOracle()).toLowerCase(),
        )
      })

      it('records the data given to it by the old oracle contract', async () => {
        await h.fulfillOracleRequest(oc, request, response)

        const currentPrice = await uc.currentPrice()
        assert.equal(currentPrice, response)
      })

      it('does not accept responses from the new oracle for the old requests', async () => {
        await h.assertActionThrows(async () => {
          await uc
            .connect(roles.oracleNode)
            .fulfill(request.id, h.toHex(response))
        })

        const currentPrice = await uc.currentPrice()
        assert.equal(ethers.utils.parseBytes32String(currentPrice), '')
      })

      it('still allows funds to be withdrawn from the oracle', async () => {
        await h.increaseTime5Minutes(provider)
        assertBigNum(
          0,
          await link.balanceOf(uc.address),
          'Initial balance should be 0',
        )

        await uc.cancelRequest(
          request.id,
          request.payment,
          request.callbackFunc,
          request.expiration,
        )

        assertBigNum(
          paymentAmount,
          await link.balanceOf(uc.address),
          'Oracle should have been repaid on cancellation.',
        )
      })
    })
  })
})
