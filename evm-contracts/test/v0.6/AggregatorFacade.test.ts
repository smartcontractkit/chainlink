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
import { AggregatorFacadeFactory } from '../../ethers/v0.6/AggregatorFacadeFactory'
import { OracleFactory } from '../../ethers/v0.6/OracleFactory'

let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const aggregatorFactory = new AggregatorFactory()
const oracleFactory = new OracleFactory()
const aggregatorFacadeFactory = new AggregatorFacadeFactory()

beforeAll(async () => {
  const users = await setup.users(provider)

  defaultAccount = users.roles.defaultAccount
})

describe('AggregatorFacade', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const previousResponse = h.numToBytes32(54321)
  const response = h.numToBytes32(67890)
  const decimals = 18

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.CallableOverrideInstance<AggregatorFactory>
  let oc1: contract.Instance<OracleFactory>
  let facade: contract.CallableOverrideInstance<AggregatorFacadeFactory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = contract.callableAggregator(
      await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, 0, 1, [oc1.address], [jobId1]),
    )
    facade = contract.callableAggregator(
      await aggregatorFacadeFactory
        .connect(defaultAccount)
        .deploy(aggregator.address, decimals),
    )

    let requestTx = await aggregator.requestRateUpdate()
    let receipt = await requestTx.wait()
    let request = oracle.decodeRunRequest(receipt.logs?.[3])
    await oc1.fulfillOracleRequest(
      ...oracle.convertFufillParams(request, previousResponse),
    )
    requestTx = await aggregator.requestRateUpdate()
    receipt = await requestTx.wait()
    request = oracle.decodeRunRequest(receipt.logs?.[3])
    await oc1.fulfillOracleRequest(
      ...oracle.convertFufillParams(request, response),
    )
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(aggregatorFacadeFactory, [
      'aggregator',
      'decimals',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
    ])
  })

  describe('#getAnswer/latestAnswer', () => {
    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(response, await facade.latestAnswer())
      const latestRound = await facade.latestRound()
      matchers.bigNum(response, await facade.getAnswer(latestRound))
    })
  })

  describe('#getTimestamp/latestTimestamp', () => {
    it('pulls the timestamp from the aggregator', async () => {
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
      matchers.bigNum(height, await facade.latestTimestamp())
      const latestRound = await facade.latestRound()
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await facade.getTimestamp(latestRound),
      )
    })
  })

  describe('#getRoundData', () => {
    it('assembles the requested round data', async () => {
      const previousId = (await facade.latestRound()).sub(1)
      const round = await facade.getRoundData(previousId)
      matchers.bigNum(previousId, round.roundId)
      matchers.bigNum(previousResponse, round.answer)
      matchers.bigNum(await facade.getTimestamp(previousId), round.startedAt)
      matchers.bigNum(await facade.getTimestamp(previousId), round.updatedAt)
      matchers.bigNum(previousId, round.answeredInRound)
    })

    it('returns zero data for non-existing rounds', async () => {
      const roundId = 13371337
      const round = await facade.getRoundData(roundId)
      matchers.bigNum(roundId, round.roundId)
      matchers.bigNum(0, round.answer)
      matchers.bigNum(0, round.startedAt)
      matchers.bigNum(0, round.updatedAt)
      matchers.bigNum(0, round.answeredInRound)
    })
  })

  describe('#latestRoundData', () => {
    it('assembles the requested round data', async () => {
      const latestId = await facade.latestRound()
      const round = await facade.latestRoundData()
      matchers.bigNum(latestId, round.roundId)
      matchers.bigNum(response, round.answer)
      matchers.bigNum(await facade.getTimestamp(latestId), round.startedAt)
      matchers.bigNum(await facade.getTimestamp(latestId), round.updatedAt)
      matchers.bigNum(latestId, round.answeredInRound)
    })
  })
})
