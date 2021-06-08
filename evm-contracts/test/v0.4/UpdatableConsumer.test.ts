import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ENSRegistry__factory } from '../../ethers/v0.4/factories/ENSRegistry__factory'
import { Oracle__factory } from '../../ethers/v0.4/factories/Oracle__factory'
import { PublicResolver__factory } from '../../ethers/v0.4/factories/PublicResolver__factory'
import { UpdatableConsumer__factory } from '../../ethers/v0.4/factories/UpdatableConsumer__factory'

const linkTokenFactory = new contract.LinkToken__factory()
const ensRegistryFactory = new ENSRegistry__factory()
const oracleFactory = new Oracle__factory()
const publicResolverFacotory = new PublicResolver__factory()
const updatableConsumerFactory = new UpdatableConsumer__factory()

const provider = setup.provider()

let roles: setup.Roles

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
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

  let ens: contract.Instance<ENSRegistry__factory>
  let ensResolver: contract.Instance<PublicResolver__factory>
  let link: contract.Instance<contract.LinkToken__factory>
  let oc: contract.Instance<Oracle__factory>
  let uc: contract.Instance<UpdatableConsumer__factory>
  const deployment = setup.snapshot(provider, async () => {
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
    let request: oracle.RunRequest

    beforeEach(async () => {
      await link.transfer(uc.address, paymentAmount)
      const tx = await uc.requestEthereumPrice(
        h.toHex(ethers.utils.toUtf8Bytes(currency)),
      )
      const receipt = await tx.wait()
      request = oracle.decodeRunRequest(receipt.logs?.[3])
    })

    it('records the data given to it by the oracle', async () => {
      await oc.fulfillOracleRequest(
        ...oracle.convertFufillParams(request, response),
      )

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
        await oc.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, response),
        )

        const currentPrice = await uc.currentPrice()
        assert.equal(currentPrice, response)
      })

      it('does not accept responses from the new oracle for the old requests', async () => {
        await matchers.evmRevert(async () => {
          await uc
            .connect(roles.oracleNode)
            .fulfill(request.requestId, h.toHex(response))
        })

        const currentPrice = await uc.currentPrice()
        assert.equal(ethers.utils.parseBytes32String(currentPrice), '')
      })

      it('still allows funds to be withdrawn from the oracle', async () => {
        await h.increaseTime5Minutes(provider)
        matchers.bigNum(
          0,
          await link.balanceOf(uc.address),
          'Initial balance should be 0',
        )

        await uc.cancelRequest(
          request.requestId,
          request.payment,
          request.callbackFunc,
          request.expiration,
        )

        matchers.bigNum(
          paymentAmount,
          await link.balanceOf(uc.address),
          'Oracle should have been repaid on cancellation.',
        )
      })
    })
  })
})
