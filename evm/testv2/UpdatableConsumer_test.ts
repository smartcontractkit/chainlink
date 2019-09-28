import * as h from '../src/helpersV2'
import ganache from 'ganache-core'
import { assertBigNum } from '../src/matchersV2'
import { AbstractContract } from '../src/contract'
import { linkToken } from '../src/linkToken'
import { ENSRegistry } from 'contracts/ENSRegistry'
import { Oracle } from 'contracts/Oracle'
import { PublicResolver } from 'contracts/PublicResolver'
import { UpdatableConsumer } from 'contracts/UpdatableConsumer'
import { LinkTokenInterface } from 'contracts/LinkTokenInterface'
import { FF } from 'src/types'
import { ethers } from 'ethers'
import { assert } from 'chai'

const ganacheProvider: any = ganache.provider()

const LinkContract = AbstractContract.fromBuildArtifact(linkToken)
const ENSRegistryContract = AbstractContract.fromArtifactName('ENSRegistry')
const OracleContract = AbstractContract.fromArtifactName('Oracle')
const PublicResolverContract = AbstractContract.fromArtifactName(
  'PublicResolver',
)
const UpdatableConsumerContract = AbstractContract.fromArtifactName(
  'UpdatableConsumer',
)

let roles: h.Roles

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(ganacheProvider)

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

  let ens: FF<ENSRegistry>
  let ensResolver: FF<PublicResolver>
  let link: FF<LinkTokenInterface>
  let oc: FF<Oracle>
  let uc: FF<UpdatableConsumer>

  beforeEach(async () => {
    link = await LinkContract.deploy(roles.defaultAccount)
    oc = await OracleContract.deploy(roles.oracleNode, [link.address])

    ens = await ENSRegistryContract.deploy(roles.defaultAccount)
    ensResolver = await PublicResolverContract.deploy(roles.defaultAccount, [
      ens.address,
    ])
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
    uc = await UpdatableConsumerContract.deploy(roles.defaultAccount, [
      specId,
      ens.address,
      domainNode,
    ])
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

    context(
      'when the oracle address is updated before a request is fulfilled',
      () => {
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
          await h.increaseTime5Minutes(ganacheProvider)
          assertBigNum(
            0,
            (await link.balanceOf(uc.address)) as any,
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
            (await link.balanceOf(uc.address)) as any,
            'Oracle should have been repaid on cancellation.',
          )
        })
      },
    )
  })
})
