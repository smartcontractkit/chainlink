import { assert } from 'chai'
import { ethers } from 'ethers'
import { FluxAggregatorFactory } from '@chainlink/contracts/ethers/v0.6/FluxAggregatorFactory'
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
  MINIMUM_CONTRACT_PAYMENT,
} = t.getEnvVars([
  'NODE_1_CONTAINER',
  'NODE_2_CONTAINER',
  'CLIENT_NODE_URL',
  'CLIENT_NODE_2_URL',
  'EXTERNAL_ADAPTER_URL',
  'EXTERNAL_ADAPTER_2_URL',
  'MINIMUM_CONTRACT_PAYMENT',
])

const provider = t.createProvider()
const carol = ethers.Wallet.createRandom().connect(provider)
const linkTokenFactory = new contract.LinkTokenFactory(carol)
const fluxAggregatorFactory = new FluxAggregatorFactory(carol)
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
let fluxAggregator: contract.Instance<FluxAggregatorFactory>

let node1Address: string
let node2Address: string

async function assertAggregatorValues(
  latestAnswer: number,
  latestRound: number,
  reportingRound: number,
  latestSubmission1: number,
  latestSubmission2: number,
  msg: string,
): Promise<void> {
  const [la, lr, rr, ls1, ls2] = await Promise.all([
    fluxAggregator.latestAnswer(),
    fluxAggregator.latestRound(),
    fluxAggregator.reportingRound(),
    fluxAggregator.latestSubmission(node1Address).then(res => res[1]),
    fluxAggregator.latestSubmission(node2Address).then(res => res[1]),
  ])

  matchers.bigNum(latestAnswer, la, `${msg} : latest answer`)
  matchers.bigNum(latestRound, lr, `${msg} : latest round`)
  matchers.bigNum(reportingRound, rr, `${msg} : reporting round`)
  matchers.bigNum(latestSubmission1, ls1, `${msg} : node 1 latest submission`)
  matchers.bigNum(latestSubmission2, ls2, `${msg} : node 2 latest submission`)
}

async function assertLatestAnswerEq(n: number) {
  await t.assertAsync(
    async () => (await fluxAggregator.latestAnswer()).eq(n),
    `latestAnswer should eventually equal ${n}`,
  )
}

beforeAll(async () => {
  t.printHeading('Flux Monitor Test')

  clClient1.login()
  clClient2.login()
  node1Address = clClient1.getAdminInfo()[0].address
  node2Address = clClient2.getAdminInfo()[0].address

  await t.fundAddress(carol.address)
  await t.fundAddress(node1Address)
  await t.fundAddress(node2Address)
  linkToken = await linkTokenFactory.deploy()
  await linkToken.deployed()

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
    MINIMUM_CONTRACT_PAYMENT,
    10,
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
  await Promise.all([clClient1.unpause(), clClient2.unpause()])
  clClient1.getJobs().forEach(job => clClient1.archiveJob(job.id))
  clClient2.getJobs().forEach(job => clClient2.archiveJob(job.id))
  await Promise.all([
    t.changePriceFeed(EXTERNAL_ADAPTER_URL, 100), // original price
    t.changePriceFeed(EXTERNAL_ADAPTER_2_URL, 100),
  ])
  fluxAggregator.removeAllListeners('*')
})

describe('FluxMonitor / FluxAggregator integration with one node', () => {
  it('updates the price', async () => {
    await linkToken.transfer(fluxAggregator.address, deposit).then(t.txWait)
    await fluxAggregator.updateAvailableFunds().then(t.txWait)

    await fluxAggregator
      .addOracles([node1Address], [node1Address], 1, 1, 0)
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

    await fluxAggregator
      .addOracles(
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
    clClient1.createJob(JSON.stringify(fluxMonitorJob))
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_2_URL]
    clClient2.createJob(JSON.stringify(fluxMonitorJob))

    // initial job run
    await t.assertJobRun(clClient1, node1InitialRunCount + 1, 'initial update')
    await t.assertJobRun(clClient2, node2InitialRunCount + 1, 'initial update')
    await assertAggregatorValues(10000, 1, 1, 1, 1, 'initial round')
    await clClient1.pause()
    await clClient2.pause()

    // node 1 should still begin round even with unresponsive node 2
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 110)
    await t.changePriceFeed(EXTERNAL_ADAPTER_2_URL, 120)
    await clClient1.unpause()
    await t.assertJobRun(clClient1, node1InitialRunCount + 2, 'second update')
    await assertAggregatorValues(10000, 1, 2, 2, 1, 'node 1 only')
    await clClient1.pause()

    // node 2 should finish round
    await clClient2.unpause()
    await t.assertJobRun(clClient2, node2InitialRunCount + 2, 'second update')
    await assertAggregatorValues(11500, 2, 2, 2, 2, 'second round')
    await clClient2.pause()

    // reduce minAnswers to 1
    await (
      await fluxAggregator.updateFutureRounds(
        MINIMUM_CONTRACT_PAYMENT,
        1,
        2,
        0,
        5,
      )
    ).wait()
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 130)
    await clClient1.unpause()
    await t.assertJobRun(clClient1, node1InitialRunCount + 3, 'third update')
    await assertAggregatorValues(13000, 3, 3, 3, 2, 'third round')
    await clClient1.pause()

    // node should continue to start new rounds alone
    await t.changePriceFeed(EXTERNAL_ADAPTER_URL, 140)
    await clClient1.unpause()
    await t.assertJobRun(clClient1, node1InitialRunCount + 4, 'fourth update')
    await assertAggregatorValues(14000, 4, 4, 4, 2, 'fourth round')

    await clClient2.unpause()
  })

  it('respects the idle timer duration', async () => {
    await (
      await fluxAggregator.updateFutureRounds(
        MINIMUM_CONTRACT_PAYMENT,
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
    clClient1.createJob(JSON.stringify(fluxMonitorJob))
    fluxMonitorJob.initiators[0].params.feeds = [EXTERNAL_ADAPTER_2_URL]
    clClient2.createJob(JSON.stringify(fluxMonitorJob))

    // initial job run
    await t.assertJobRun(clClient1, node1InitialRunCount + 1, 'initial update')
    await t.assertJobRun(clClient2, node2InitialRunCount + 1, 'initial update')
    await assertAggregatorValues(10000, 1, 1, 1, 1, 'initial round')

    // second job run
    await t.assertJobRun(clClient1, node1InitialRunCount + 2, 'second update')
    await t.assertJobRun(clClient2, node2InitialRunCount + 2, 'second update')
    await assertAggregatorValues(10000, 2, 2, 2, 2, 'second round')

    // third job run without node 2
    await clClient2.pause()
    await t.assertJobRun(clClient1, node1InitialRunCount + 3, 'third update')
    await assertAggregatorValues(10000, 2, 3, 3, 2, 'third round')
  })
})
