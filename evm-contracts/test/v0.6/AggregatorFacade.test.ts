import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { Aggregator__factory } from '../../ethers/v0.4/factories/Aggregator__factory'
import { AggregatorFacade__factory } from '../../ethers/v0.6/factories/AggregatorFacade__factory'
import { Oracle__factory } from '../../ethers/v0.6/factories/Oracle__factory'

let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkToken__factory()
const aggregatorFactory = new Aggregator__factory()
const oracleFactory = new Oracle__factory()
const aggregatorFacadeFactory = new AggregatorFacade__factory()

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
  const description = 'LINK / USD: Historic Aggregator Facade'

  let link: contract.Instance<contract.LinkToken__factory>
  let aggregator: contract.Instance<Aggregator__factory>
  let oc1: contract.Instance<Oracle__factory>
  let facade: contract.Instance<AggregatorFacade__factory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(link.address, 0, 1, [oc1.address], [jobId1])
    facade = await aggregatorFacadeFactory
      .connect(defaultAccount)
      .deploy(aggregator.address, decimals, description)

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
      'description',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'version',
    ])
  })

  describe('#constructor', () => {
    it('uses the decimals set in the constructor', async () => {
      matchers.bigNum(decimals, await facade.decimals())
    })

    it('uses the description set in the constructor', async () => {
      assert.equal(description, await facade.description())
    })

    it('sets the version to 2', async () => {
      matchers.bigNum(2, await facade.version())
    })
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
      await matchers.evmRevert(facade.getRoundData(roundId), 'No data present')
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

    describe('when there is no latest round', () => {
      beforeEach(async () => {
        aggregator = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, 0, 1, [oc1.address], [jobId1])
        facade = await aggregatorFacadeFactory
          .connect(defaultAccount)
          .deploy(aggregator.address, decimals, description)
      })

      it('assembles the requested round data', async () => {
        await matchers.evmRevert(facade.latestRoundData(), 'No data present')
      })
    })
  })
})
