import cbor from 'cbor'
import * as h from '../src/helpersV2'
import { AbstractContract } from '../src/contract'
import { assertBigNum } from '../src/matchersV2'
import ganache from 'ganache-core'
import { ethers } from 'ethers'
import { assert } from 'chai'
import { linkToken } from '../src/linkToken'
import { LinkTokenInterface } from 'contracts/LinkTokenInterface'
import { Oracle } from 'contracts/Oracle'
import { BasicConsumer } from 'contracts/BasicConsumer'
import { FF } from 'src/types'
const BasicConsumerContract = AbstractContract.fromArtifactName('BasicConsumer')
const OracleContract = AbstractContract.fromArtifactName('Oracle')
const LinkContract = AbstractContract.fromBuildArtifact(linkToken)

// create a web3js instance connected to in-memory ganache blockchain
const ganacheProvider: any = ganache.provider()

// create ethers provider from that web3js instance
const provider = new ethers.providers.Web3Provider(ganacheProvider)

let roles: h.Roles

before(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(ganacheProvider)

  roles = rolesAndPersonas.roles
})

describe('BasicConsumer', () => {
  const specId = '0x4c7b7ffb66b344fbaa64995af81e355a'.padEnd(66, '0')

  const currency = 'USD'
  let link: FF<LinkTokenInterface>
  let oc: FF<Oracle>
  let cc: FF<BasicConsumer>

  beforeEach(async () => {
    link = await LinkContract.deploy(roles.defaultAccount)
    oc = await OracleContract.deploy(roles.oracleNode, [link.address])
    cc = await BasicConsumerContract.deploy(roles.defaultAccount, [
      link.address,
      oc.address,
      specId,
    ])
  })

  it('has a predictable gas price', async () => {
    const rec = await provider.getTransactionReceipt(cc.deployTransaction.hash!)
    assert.isBelow(rec.gasUsed!.toNumber(), 1700000)
  })

  describe('#requestEthereumPrice', () => {
    context('without LINK', () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await cc.requestEthereumPrice(currency)
        })
      })
    })

    context('with LINK', () => {
      beforeEach(async () => {
        await link.transfer(cc.address, ethers.utils.parseEther('1'))
      })

      it('triggers a log event in the Oracle contract', async () => {
        const tx = await cc.requestEthereumPrice(currency)
        const receipt = await tx.wait()
        const log = receipt.logs![3]
        assert.equal(log.address, oc.address)

        const request = h.decodeRunRequest(log)
        const expected = {
          path: ['USD'],
          get:
            'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY',
        }

        assert.equal(h.toHex(specId), request.jobId)
        assertBigNum(h.toWei('1'), request.payment)
        assert.equal(cc.address, request.requester)
        assert.equal(1, request.dataVersion)
        assert.deepEqual(expected, cbor.decodeFirstSync(request.data))
      })

      it('has a reasonable gas cost', async () => {
        const tx = await cc.requestEthereumPrice(currency)
        const receipt = await tx.wait()
        assert.isBelow(receipt.gasUsed!.toNumber(), 120000)
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = ethers.utils.formatBytes32String('1,000,000.00')
    let request: h.RunRequest

    beforeEach(async () => {
      await link.transfer(cc.address, h.toWei('1'))
      const tx = await cc.requestEthereumPrice(currency)
      const receipt = await tx.wait()
      request = h.decodeRunRequest(receipt.logs![3])
    })

    it('records the data given to it by the oracle', async () => {
      await h.fulfillOracleRequest(
        oc.connect(roles.oracleNode),
        request,
        response,
        { gasLimit: 1000000 }, // FIXME: incorrect gas estimation
      )

      const currentPrice = await cc.currentPrice()
      assert.equal(currentPrice, response)
    })

    it('logs the data given to it by the oracle', async () => {
      const tx = await h.fulfillOracleRequest(
        oc.connect(roles.oracleNode),
        request,
        response,
        { gasLimit: 1000000 }, // FIXME: incorrect gas estimation
      )
      const receipt = await tx.wait()
      assert.equal(2, receipt.logs!.length)
      const log = receipt.logs![1]

      assert.equal(log.topics[2], response)
    })

    context('when the consumer does not recognize the request ID', () => {
      let otherRequest: h.RunRequest

      beforeEach(async () => {
        const funcSig = ethers.utils
          .id('fulfill(bytes32,bytes32)')
          .slice(2)
          .slice(0, 8)
        const args = h.requestDataBytes(
          specId,
          cc.address,
          `0x${funcSig}`,
          43,
          '0x0',
        )
        const tx = await h.requestDataFrom(oc, link, 0, args)
        const receipt = await tx.wait()
        otherRequest = h.decodeRunRequest(receipt.logs![2])
      })

      it('does not accept the data provided', async () => {
        await h.fulfillOracleRequest(
          oc.connect(roles.oracleNode),
          otherRequest,
          response,
          { gasLimit: 1000000 }, // FIXME: incorrect gas estimation
        )

        const received = await cc.currentPrice()

        assert.equal(received, ethers.constants.HashZero)
      })
    })

    context('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await h.assertActionThrows(async () => {
          await cc.connect(roles.oracleNode).fulfill(request.id, response)
        })

        const received = await cc.currentPrice()
        assert.equal(received, ethers.constants.HashZero)
      })
    })
  })

  describe('#cancelRequest', () => {
    const depositAmount = h.toWei('1')
    let request: h.RunRequest

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      const tx = await cc.requestEthereumPrice(currency)
      const receipt = await tx.wait()
      request = h.decodeRunRequest(receipt.logs![3])
    })

    context('before 5 minutes', () => {
      it('cant cancel the request', async () => {
        await h.assertActionThrows(async () => {
          await cc
            .connect(roles.consumer)
            .cancelRequest(
              request.id,
              request.payment,
              request.callbackFunc,
              request.expiration,
            )
        })
      })
    })

    context('after 5 minutes', () => {
      it('can cancel the request', async () => {
        await h.increaseTime5Minutes(ganacheProvider)

        await cc
          .connect(roles.consumer)
          .cancelRequest(
            request.id,
            request.payment,
            request.callbackFunc,
            request.expiration,
          )
      })
    })
  })

  describe('#withdrawLink', () => {
    const depositAmount = h.toWei('1')

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      const balance = (await link.balanceOf(cc.address)) as any
      assertBigNum(balance, depositAmount)
    })

    it('transfers LINK out of the contract', async () => {
      await cc.connect(roles.consumer).withdrawLink()
      const ccBalance = (await link.balanceOf(cc.address)) as any
      const consumerBalance = ethers.utils.bigNumberify((await link.balanceOf(
        roles.consumer.address,
      )) as any)
      assertBigNum(ccBalance, 0)
      assertBigNum(consumerBalance, depositAmount)
    })
  })
})
