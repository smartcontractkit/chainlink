import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { AggregatorFactory } from '../../ethers/v0.4/AggregatorFactory'
import { WhitelistedAggregatorProxyFactory } from '../../ethers/v0.6/WhitelistedAggregatorProxyFactory'
import { OracleFactory } from '../../ethers/v0.6/OracleFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const aggregatorFactory = new AggregatorFactory()
const oracleFactory = new OracleFactory()
const whitelistedAggregatorProxyFactory = new WhitelistedAggregatorProxyFactory()

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('AggregatorProxy', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  const response = h.numToBytes32(54321)

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.Instance<AggregatorFactory>
  let oc1: contract.Instance<OracleFactory>
  let proxy: contract.Instance<WhitelistedAggregatorProxyFactory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
    await link.transfer(aggregator.address, deposit)
    proxy = await whitelistedAggregatorProxyFactory
      .connect(defaultAccount)
      .deploy(aggregator.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(whitelistedAggregatorProxyFactory, [
      'aggregator',
      'getAnswer',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestTimestamp',
      'setAggregator',
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // Whitelisted methods:
      'addToWhitelist',
      'disableWhitelist',
      'enableWhitelist',
      'removeFromWhitelist',
      'whitelistEnabled',
      'whitelisted',
    ])
  })

  describe('if the caller is not whitelisted', () => {
    it('latestAnswer reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestAnswer()
      }, 'Not whitelisted')
    })

    it('latestTimestamp reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestTimestamp()
      }, 'Not whitelisted')
    })

    it('getAnswer reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getAnswer(1)
      }, 'Not whitelisted')
    })

    it('getTimestamp reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getTimestamp(1)
      }, 'Not whitelisted')
    })

    it('latestRound reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestRound()
      }, 'Not whitelisted')
    })
  })

  describe('if the caller is whitelisted', () => {
    beforeEach(async () => {
      await proxy.addToWhitelist(defaultAccount.address)

      const requestTx = await aggregator.requestRateUpdate()
      const receipt = await requestTx.wait()
      const request = oracle.decodeRunRequest(receipt.logs?.[3])
      await oc1.fulfillOracleRequest(
        ...oracle.convertFufillParams(request, response),
      )

      matchers.bigNum(
        ethers.utils.bigNumberify(response),
        await aggregator.latestAnswer(),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(response, await proxy.latestAnswer())
      const latestRound = await proxy.latestRound()
      matchers.bigNum(response, await proxy.getAnswer(latestRound))
    })

    it('pulls the timestamp from the aggregator', async () => {
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.latestTimestamp(),
      )
      const latestRound = await proxy.latestRound()
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.getTimestamp(latestRound),
      )
    })
  })
})
