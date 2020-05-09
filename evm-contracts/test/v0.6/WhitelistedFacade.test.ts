import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { ethers } from 'ethers'
import { AggregatorFactory } from '../../ethers/v0.4/AggregatorFactory'
import { WhitelistedFacadeFactory } from '../../ethers/v0.6/WhitelistedFacadeFactory'
import { OracleFactory } from '../../ethers/v0.6/OracleFactory'

let defaultAccount: ethers.Wallet
let personas: setup.Personas

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const aggregatorFactory = new AggregatorFactory()
const oracleFactory = new OracleFactory()
const whitelistedFacadeFactory = new WhitelistedFacadeFactory()

beforeAll(async () => {
  await setup.users(provider).then(u => (personas = u.personas))
  defaultAccount = personas.Default
})

describe('WhitelistedFacade', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const previousResponse = h.numToBytes32(54321)
  const response = h.numToBytes32(67890)
  const decimals = 18

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.Instance<AggregatorFactory>
  let oc1: contract.Instance<OracleFactory>
  let facade: contract.Instance<WhitelistedFacadeFactory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = await aggregatorFactory
      .connect(personas.Carol)
      .deploy(link.address, 0, 1, [oc1.address], [jobId1])
    facade = await whitelistedFacadeFactory
      .connect(personas.Carol)
      .deploy(aggregator.address, decimals)
    await facade.connect(personas.Carol).addToWhitelist(personas.Carol.address)

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
    matchers.publicAbi(whitelistedFacadeFactory, [
      'aggregator',
      'decimals',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      // Owned methods:
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

  describe('#getAnswer/latestAnswer', () => {
    describe('when the reader is not whitelisted', () => {
      it('does not allow the answer to be read', async () => {
        await matchers.evmRevert(
          facade.connect(personas.Eddy).latestRound(),
          'Not whitelisted',
        )
        const round = await aggregator.connect(personas.Carol).latestRound()
        await matchers.evmRevert(
          facade.connect(personas.Eddy).getAnswer(round),
          'Not whitelisted',
        )
        await matchers.evmRevert(
          facade.connect(personas.Eddy).latestAnswer(),
          'Not whitelisted',
        )
      })
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await facade
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('pulls the answer from the aggregator', async () => {
        const round = await aggregator.connect(personas.Eddy).latestRound()
        const getAnswer = await facade.connect(personas.Eddy).getAnswer(round)
        const latestAnswer = await facade.connect(personas.Eddy).latestAnswer()
        matchers.bigNum(getAnswer, latestAnswer)
      })
    })
  })

  describe('#getTimestamp/latestTimestamp', () => {
    describe('when the reader is not whitelisted', () => {
      it('does not allow the timestamp to be read', async () => {
        const round = await aggregator.connect(personas.Carol).latestRound()
        await matchers.evmRevert(
          facade.connect(personas.Eddy).getTimestamp(round),
          'Not whitelisted',
        )
        await matchers.evmRevert(
          facade.connect(personas.Eddy).latestTimestamp(),
          'Not whitelisted',
        )
      })
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await facade
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('pulls the timestamp from the aggregator', async () => {
        const round = await aggregator.connect(personas.Eddy).latestRound()
        const getTimestamp = await facade
          .connect(personas.Eddy)
          .getTimestamp(round)
        const latestTimestamp = await facade
          .connect(personas.Eddy)
          .latestTimestamp()
        matchers.bigNum(getTimestamp, latestTimestamp)
      })
    })
  })

  describe('#getRoundData/latestRoundData', () => {
    it('does not allow the round data to be read', async () => {
      const round = await aggregator.connect(personas.Carol).latestRound()
      await matchers.evmRevert(
        facade.connect(personas.Eddy).getRoundData(round),
        'Not whitelisted',
      )
      await matchers.evmRevert(
        facade.connect(personas.Eddy).latestRoundData(),
        'Not whitelisted',
      )
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await facade
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('pulls the round data from the aggregator', async () => {
        const round = await aggregator.connect(personas.Eddy).latestRound()
        const getRound = await facade.connect(personas.Eddy).getRoundData(round)
        const latestRound = await facade
          .connect(personas.Eddy)
          .latestRoundData()
        matchers.bigNum(getRound.roundId, latestRound.roundId)
      })
    })
  })
})
