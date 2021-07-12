import { assert } from 'chai'
import { ethers } from 'ethers'
import { FluxAggregator__factory } from '@chainlink/contracts/ethers/v0.6/factories/FluxAggregator__factory'
import { contract, helpers as h, matchers } from '@chainlink/test-helpers'
import ChainlinkClient from '../test-helpers/chainlinkClient'
import fluxMonitorJobTemplate from '../fixtures/flux-monitor-job'
import * as t from '../test-helpers/common'

jest.unmock('execa').unmock('dockerode')

const {
  NODE_1_CONTAINER,
  NODE_2_CONTAINER,
  CLIENT_NODE_URL,
  CLIENT_NODE_2_URL,
  EXTERNAL_ADAPTER_URL,
  EXTERNAL_ADAPTER_2_URL,
  MINIMUM_CONTRACT_PAYMENT_LINK_JUELS,
} = t.getEnvVars([
  'NODE_1_CONTAINER',
  'NODE_2_CONTAINER',
  'CLIENT_NODE_URL',
  'CLIENT_NODE_2_URL',
  'EXTERNAL_ADAPTER_URL',
  'EXTERNAL_ADAPTER_2_URL',
  'MINIMUM_CONTRACT_PAYMENT_LINK_JUELS',
])

const provider = t.createProvider()
const carol = ethers.Wallet.createRandom().connect(provider)
const linkTokenFactory = new contract.LinkTokenFactory(carol)
const fluxAggregatorFactory = new FluxAggregator__factory(carol)
const deposit = h.toWei('1000')
const emptyAddress = '0x0000000000000000000000000000000000000000'

const answerUpdated = fluxAggregatorFactory.interface.events.AnswerUpdated.name
const oracleAdded =
  fluxAggregatorFactory.interface.events.OraclePermissionsUpdated.name
const submissionReceived =
  fluxAggregatorFactory.interface.events.SubmissionReceived.name
const roundDetailsUpdated =
  fluxAggregatorFactory.interface.events.RoundDetailsUpdated.name
const faEventsToListenTo = [
  answerUpdated,
  oracleAdded,
  submissionReceived,
  roundDetailsUpdated,
]

const clClient1 = new ChainlinkClient(
  'node 1',
  CLIENT_NODE_URL,
  NODE_1_CONTAINER,
)
const clClient2 = new ChainlinkClient(
  'node 2',
  CLIENT_NODE_2_URL,
  NODE_2_CONTAINER,
)

// TODO import JobSpecRequest from operator_ui/@types/core/store/models.d.ts
// https://www.pivotaltracker.com/story/show/171715396
let fluxMonitorJob: any
let linkToken: contract.Instance<contract.LinkTokenFactory>
let fluxAggregator: contract.Instance<FluxAggregator__factory>

let node1Address: string
let node2Address: string
let txInterval: number | undefined

async function assertAggregatorValues(
  latestAnswer: number,
  latestRound: number,
  reportingRound: number,
  roundState1: number,
  roundState2: number,
  msg: string,
): Promise<void> {
  const [la, lr, rr, ls1, ls2] = await Promise.all([
    fluxAggregator.latestRoundData().then((res) => res.answer),
    fluxAggregator.latestRoundData().then((res) => res.roundId),
    // get earliest eligible round ID by checking a non-existent address
    fluxAggregator
      .oracleRoundState(emptyAddress, 0)
      .then((res) => res._roundId),
    fluxAggregator
      .oracleRoundState(node1Address, 0)
      .then((res) => res._roundId),
    fluxAggregator
      .oracleRoundState(node2Address, 0)
      .then((res) => res._roundId),
  ])
  matchers.bigNum(latestAnswer, la, `${msg} : latest answer`)
  matchers.bigNum(latestRound, lr, `${msg} : latest round`)
  matchers.bigNum(reportingRound, rr, `${msg} : reporting round`)
  matchers.bigNum(roundState1, ls1, `${msg} : node 1 round state round ID`)
  matchers.bigNum(roundState2, ls2, `${msg} : node 2 round state round ID`)
}

async function assertLatestAnswerEq(n: number) {
  await t.assertAsync(async () => {
    const round = await fluxAggregator.latestRoundData()
    return round.answer.eq(n)
  }, `latestAnswer should eventually equal ${n}`)
}

beforeAll(async () => {
  t.printHeading('Flux Monitor Test')

  clClient1.login()
  clClient2.login()
  node1Address = clClient1.newEthKey().address
  node2Address = clClient2.newEthKey().address
  console.log('new eth keys', node1Address, node2Address)

  await t.fundAddress(carol.address)
  await t.fundAddress(node1Address)
  await t.fundAddress(node2Address)
  linkToken = await linkTokenFactory.deploy()
  await linkToken.deployed()

  // Set up a recurring tx send in the background to ensure that the
  // ETH node is continually mining blocks despite being in dev mode.
  const miner = ethers.Wallet.createRandom().connect(provider)
  await t.fundAddress(miner.address)
  txInterval = t.setRecurringTx(miner)

  console.log(`Chainlink Node 1 address: ${node1Address}`)
  console.log(`Chainlink Node 2 address: ${node2Address}`)
  console.log(`Contract creator's address: ${carol.address}`)
  console.log(`Deployed LinkToken contract: ${linkToken.address}`)
})

beforeEach(async () => {
  t.printHeading('Running Test')
  fluxMonitorJob = JSON.parse(JSON.stringify(fluxMonitorJobTemplate)) // perform a deep clone
  const minSubmissionValue = 1
  const maxSubmissionValue = 1000000000
  const deployingContract = await fluxAggregatorFactory.deploy(
    linkToken.address,
    MINIMUM_CONTRACT_PAYMENT_LINK_JUELS,
    300,
    emptyAddress,
    minSubmissionValue,
    maxSubmissionValue,
    1,
    ethers.utils.formatBytes32String('ETH/USD'),
  )
  await deployingContract.deployed()
  fluxAggregator = deployingContract
  t.logEvents(fluxAggregator as any, 'FluxAggregator', faEventsToListenTo)
  console.log(`Deployed FluxAggregator contract: ${fluxAggregator.address}`)
})

afterEach(async () => {
  clClient1.getJobs().forEach((job) => clClient1.archiveJob(job.id))
  clClient2.getJobs().forEach((job) => clClient2.archiveJob(job.id))
  await Promise.all([
    t.changePriceFeed(EXTERNAL_ADAPTER_URL, 100), // original price
    t.changePriceFeed(EXTERNAL_ADAPTER_2_URL, 100),
  ])
  fluxAggregator.removeAllListeners('*')
})

afterAll(() => {
  clearInterval(txInterval)
})

describe('FluxMonitor / FluxAggregator integration with one node', () => {
  it('updates the price', async () => {
    await linkToken.transfer(fluxAggregator.address, deposit).then(t.txWait)
    await fluxAggregator.updateAvailableFunds().then(t.txWait)

    await fluxAggregator
      .changeOracles([], [node1Address], [node1Address], 1, 1, 0)
      .then(t.txWait)

    expect(await fluxAggregator.getOracles()).toEqual([node1Address])
    matchers.bigNum(
      await linkToken.balanceOf(fluxAggregator.address),
      deposit,
      'Unable to fund FluxAggregator',
    )

    const initialJobCount = clClient1.getJobs().length
    const initialRunCount = clClient1.getJobRuns().length

    // create FM job
    fluxMonitorJob.initiators[0].params.address = fluxAggregator.address
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_URL]
    fluxMonitorJob.tasks[2].params.fromAddress = node1Address
    clClient1.createJob(JSON.stringify(fluxMonitorJob))
    assert.equal(clClient1.getJobs().length, initialJobCount + 1)

    // Job should trigger initial FM run
    await t.assertJobRun(clClient1, initialRunCount + 1, 'initial update')
    await assertLatestAnswerEq(10000)

    // Nominally change price feed
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 101)
    await t.wait(10000)
    assert.equal(
      clClient1.getJobRuns().length,
      initialRunCount + 1,
      'Flux Monitor should not run job after nominal price deviation',
    )

    // Significantly change price feed
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 110)
    await t.assertJobRun(clClient1, initialRunCount + 2, 'second update')
    await assertLatestAnswerEq(11000)
  })
})

describe('FluxMonitor / FluxAggregator integration with two nodes', () => {
  beforeEach(async () => {
    await linkToken.transfer(fluxAggregator.address, deposit).then(t.txWait)
    await fluxAggregator.updateAvailableFunds().then(t.txWait)

    console.log(await fluxAggregator.getOracles())
    await fluxAggregator
      .changeOracles(
        [],
        [node1Address, node2Address],
        [node1Address, node2Address],
        2,
        2,
        0,
      )
      .then(t.txWait)

    expect(await fluxAggregator.getOracles()).toEqual([
      node1Address,
      node2Address,
    ])
    matchers.bigNum(
      await linkToken.balanceOf(fluxAggregator.address),
      deposit,
      'Unable to fund FluxAggregator',
    )
  })

  it('updates the price', async () => {
    const node1InitialRunCount = clClient1.getJobRuns().length
    const node2InitialRunCount = clClient2.getJobRuns().length

    fluxMonitorJob.initiators[0].params.address = fluxAggregator.address
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_URL]
    fluxMonitorJob.tasks[2].params.fromAddress = node1Address
    clClient1.createJob(JSON.stringify(fluxMonitorJob))
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_2_URL]
    fluxMonitorJob.tasks[2].params.fromAddress = node2Address
    console.log(`using keys`, node1Address, node2Address)
    clClient2.createJob(JSON.stringify(fluxMonitorJob))

    // initial job run
    await t.assertJobRun(clClient1, node1InitialRunCount + 1, 'update 1 node 1')
    await t.assertJobRun(clClient2, node2InitialRunCount + 1, 'update 1 node 2')
    await assertAggregatorValues(10000, 1, 2, 2, 2, 'initial round')

    // node 1 should still begin round even with unresponsive node 2
    await clClient2.pause()
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 110)
    await t.changePriceFeed(EXTERNAL_ADAPTER_2_URL, 120)
    await t.assertJobRun(clClient1, node1InitialRunCount + 2, 'update 2 node 1')
    await assertAggregatorValues(10000, 1, 2, 2, 2, 'node 1 only')

    // node 2 should finish round
    await clClient2.unpause()
    await t.assertJobRun(clClient2, node2InitialRunCount + 2, 'update 2 node 2')
    await assertAggregatorValues(11500, 2, 3, 3, 3, 'second round')
    await clClient2.pause()

    // reduce minAnswers to 1
    await (
      await fluxAggregator.updateFutureRounds(
        MINIMUM_CONTRACT_PAYMENT_LINK_JUELS,
        1,
        2,
        0,
        300,
      )
    ).wait()
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 130)
    await t.assertJobRun(clClient1, node1InitialRunCount + 3, 'update 3')
    await assertAggregatorValues(13000, 3, 3, 4, 3, 'third round')

    // node should continue to start new rounds alone
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 140)
    await t.assertJobRun(clClient1, node1InitialRunCount + 4, 'update 4')
    await assertAggregatorValues(14000, 4, 4, 5, 4, 'fourth round')
    await clClient2.unpause()
  })

  it('respects the idle timer duration', async () => {
    await (
      await fluxAggregator.updateFutureRounds(
        MINIMUM_CONTRACT_PAYMENT_LINK_JUELS,
        2,
        2,
        0,
        10,
      )
    ).wait()

    const node1InitialRunCount = clClient1.getJobRuns().length
    const node2InitialRunCount = clClient2.getJobRuns().length

    fluxMonitorJob.initiators[0].params.idleTimer.disabled = false
    fluxMonitorJob.initiators[0].params.idleTimer.duration = '15s'
    fluxMonitorJob.initiators[0].params.pollTimer.disabled = true
    fluxMonitorJob.initiators[0].params.pollTimer.period = '0'
    fluxMonitorJob.initiators[0].params.address = fluxAggregator.address
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_URL]
    fluxMonitorJob.tasks[2].params.fromAddress = node1Address
    clClient1.createJob(JSON.stringify(fluxMonitorJob))
    fluxMonitorJob.tasks[2].params.fromAddress = node2Address
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_2_URL]
    clClient2.createJob(JSON.stringify(fluxMonitorJob))

    // initial job run
    await t.assertJobRun(clClient1, node1InitialRunCount + 1, 'initial update')
    await t.assertJobRun(clClient2, node2InitialRunCount + 1, 'initial update')
    await assertAggregatorValues(10000, 1, 2, 2, 2, 'initial round')

    // second job run
    await t.assertJobRun(clClient1, node1InitialRunCount + 2, 'second update')
    await t.assertJobRun(clClient2, node2InitialRunCount + 2, 'second update')
    await clClient2.pause()
    await assertAggregatorValues(10000, 2, 3, 3, 3, 'second round')

    // third job run without node 2
    await t.assertJobRun(clClient1, node1InitialRunCount + 3, 'third update')
    await assertAggregatorValues(10000, 2, 3, 3, 3, 'third round')

    await clClient2.unpause()
  })
})
