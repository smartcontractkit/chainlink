import {
  contract,
  debug,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import cbor from 'cbor'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { BasicConsumer__factory } from '../../ethers/v0.5/factories/BasicConsumer__factory'
import { Oracle__factory } from '../../ethers/v0.5/factories/Oracle__factory'

const d = debug.makeDebug('BasicConsumer')
const basicConsumerFactory = new BasicConsumer__factory()
const oracleFactory = new Oracle__factory()
const linkTokenFactory = new contract.LinkToken__factory()

// create ethers provider from that web3js instance
const provider = setup.provider()

let roles: setup.Roles

beforeAll(async () => {
  const users = await setup.users(provider)

  roles = users.roles
})

describe('BasicConsumer', () => {
  const specId = '0x4c7b7ffb66b344fbaa64995af81e355a'.padEnd(66, '0')
  const currency = 'USD'
  const payment = h.toWei('1')
  let link: contract.Instance<contract.LinkToken__factory>
  let oc: contract.Instance<Oracle__factory>
  let cc: contract.Instance<BasicConsumer__factory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    oc = await oracleFactory.connect(roles.oracleNode).deploy(link.address)
    cc = await basicConsumerFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, oc.address, specId)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a predictable gas price', async () => {
    const rec = await provider.getTransactionReceipt(
      cc.deployTransaction.hash ?? '',
    )
    assert.isBelow(rec.gasUsed?.toNumber() ?? -1, 1750000)
  })

  describe('#requestEthereumPrice', () => {
    describe('without LINK', () => {
      it('reverts', () =>
        matchers.evmRevert(cc.requestEthereumPrice(currency, payment)))
    })

    describe('with LINK', () => {
      beforeEach(async () => {
        await link.transfer(cc.address, h.toWei('1'))
      })

      it('triggers a log event in the Oracle contract', async () => {
        const tx = await cc.requestEthereumPrice(currency, payment)
        const receipt = await tx.wait()

        const log = receipt?.logs?.[3]
        assert.equal(log?.address.toLowerCase(), oc.address.toLowerCase())

        const request = oracle.decodeRunRequest(log)
        const expected = {
          path: ['USD'],
          get:
            'https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY',
        }

        assert.equal(h.toHex(specId), request.specId)
        matchers.bigNum(h.toWei('1'), request.payment)
        assert.equal(cc.address.toLowerCase(), request.requester.toLowerCase())
        assert.equal(1, request.dataVersion)
        assert.deepEqual(expected, cbor.decodeFirstSync(request.data))
      })

      it('has a reasonable gas cost', async () => {
        const tx = await cc.requestEthereumPrice(currency, payment)
        const receipt = await tx.wait()

        assert.isBelow(receipt?.gasUsed?.toNumber() ?? -1, 130000)
      })
    })
  })

  describe('#fulfillOracleRequest', () => {
    const response = ethers.utils.formatBytes32String('1,000,000.00')
    let request: oracle.RunRequest

    beforeEach(async () => {
      await link.transfer(cc.address, h.toWei('1'))
      const tx = await cc.requestEthereumPrice(currency, payment)
      const receipt = await tx.wait()

      request = oracle.decodeRunRequest(receipt?.logs?.[3])
    })

    it('records the data given to it by the oracle', async () => {
      await oc
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...oracle.convertFufillParams(request, response))

      const currentPrice = await cc.currentPrice()
      assert.equal(currentPrice, response)
    })

    it('logs the data given to it by the oracle', async () => {
      const tx = await oc
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...oracle.convertFufillParams(request, response))
      const receipt = await tx.wait()

      assert.equal(2, receipt?.logs?.length)
      const log = receipt?.logs?.[1]

      assert.equal(log?.topics[2], response)
    })

    describe('when the consumer does not recognize the request ID', () => {
      let otherRequest: oracle.RunRequest

      beforeEach(async () => {
        // Create a request directly via the oracle, rather than through the
        // chainlink client (consumer). The client should not respond to
        // fulfillment of this request, even though the oracle will faithfully
        // forward the fulfillment to it.
        const args = oracle.encodeOracleRequest(
          h.toHex(specId),
          cc.address,
          basicConsumerFactory.interface.functions.fulfill.sighash,
          43,
          '0x0',
        )
        const tx = await link.transferAndCall(oc.address, 0, args)
        const receipt = await tx.wait()

        otherRequest = oracle.decodeRunRequest(receipt?.logs?.[2])
      })

      it('does not accept the data provided', async () => {
        d('otherRequest %s', otherRequest)
        await oc
          .connect(roles.oracleNode)
          .fulfillOracleRequest(
            ...oracle.convertFufillParams(otherRequest, response),
          )

        const received = await cc.currentPrice()

        assert.equal(ethers.utils.parseBytes32String(received), '')
      })
    })

    describe('when called by anyone other than the oracle contract', () => {
      it('does not accept the data provided', async () => {
        await matchers.evmRevert(
          cc.connect(roles.oracleNode).fulfill(request.requestId, response),
        )

        const received = await cc.currentPrice()
        assert.equal(ethers.utils.parseBytes32String(received), '')
      })
    })
  })

  describe('#cancelRequest', () => {
    const depositAmount = h.toWei('1')
    let request: oracle.RunRequest

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount)
      const tx = await cc.requestEthereumPrice(currency, payment)
      const receipt = await tx.wait()

      request = oracle.decodeRunRequest(receipt.logs?.[3])
    })

    describe('before 5 minutes', () => {
      it('cant cancel the request', () =>
        matchers.evmRevert(
          cc
            .connect(roles.consumer)
            .cancelRequest(
              oc.address,
              request.requestId,
              request.payment,
              request.callbackFunc,
              request.expiration,
            ),
        ))
    })

    describe('after 5 minutes', () => {
      it('can cancel the request', async () => {
        await h.increaseTime5Minutes(provider)

        await cc
          .connect(roles.consumer)
          .cancelRequest(
            oc.address,
            request.requestId,
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
      const balance = await link.balanceOf(cc.address)
      matchers.bigNum(balance, depositAmount)
    })

    it('transfers LINK out of the contract', async () => {
      await cc.connect(roles.consumer).withdrawLink()
      const ccBalance = await link.balanceOf(cc.address)
      const consumerBalance = ethers.utils.bigNumberify(
        await link.balanceOf(roles.consumer.address),
      )
      matchers.bigNum(ccBalance, 0)
      matchers.bigNum(consumerBalance, depositAmount)
    })
  })
})
