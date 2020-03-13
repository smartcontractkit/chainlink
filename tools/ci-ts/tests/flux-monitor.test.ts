import { FluxAggregatorFactory } from '@chainlink/contracts/ethers/v0.6/FluxAggregatorFactory'
import { contract, helpers as h, matchers } from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import 'isomorphic-unfetch'
import { JobSpec } from '../../../operator_ui/@types/operator_ui'
import ChainlinkClient from '../test-helpers/chainlink-cli'
import fluxMonitorJob from '../fixtures/flux-monitor-job'
import {
  assertAsync,
  createProvider,
  fundAddress,
  txWait,
  wait,
} from '../test-helpers/common'

const NODE_1_URL = 'http://node:6688'
const NODE_2_URL = 'http://node-2:6688'
const NODE_1_CONTAINER = 'chainlink-node'
const NODE_2_CONTAINER = 'chainlink-node-2'
const EA_1_URL = 'http://external-adapter:6644'
const EA_2_URL = 'http://external-adapter-2:6644'

const provider = createProvider()
const carol = ethers.Wallet.createRandom().connect(provider)
const linkTokenFactory = new contract.LinkTokenFactory(carol)
const fluxAggregatorFactory = new FluxAggregatorFactory(carol)
const deposit = h.toWei('1000')
const clClient1 = new ChainlinkClient(NODE_1_URL, NODE_1_CONTAINER)
const clClient2 = new ChainlinkClient(NODE_2_URL, NODE_2_CONTAINER)

let linkToken: contract.Instance<contract.LinkTokenFactory>
let fluxAggregator: contract.Instance<FluxAggregatorFactory>
let node1Address: string
let node2Address: string

async function changePriceFeed(adapter: string, value: number) {
  const url = new URL('result', adapter).href
  const response = await fetch(url, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ result: value }),
  })
  assert(response.ok)
}

async function assertJobRun(
  clClient: ChainlinkClient,
  jobId: string,
  count: number,
  errorMessage: string,
) {
  await assertAsync(() => {
    const jobRuns = clClient.getJobRuns()
    const jobRun = jobRuns[jobRuns.length - 1]
    return (
      clClient.getJobRuns().length === count &&
      jobRun.status === 'completed' &&
      jobRun.jobId === jobId
    )
  }, errorMessage)
}

beforeAll(async () => {
  clClient1.login()
  clClient2.login()
  node1Address = clClient1.getAdminInfo()[0].address
  node2Address = clClient2.getAdminInfo()[0].address
  await fundAddress(carol.address)
  await fundAddress(node1Address)
  await fundAddress(node2Address)
  linkToken = await linkTokenFactory.deploy()
  await linkToken.deployed()
})

beforeEach(async () => {
  fluxAggregator = await fluxAggregatorFactory.deploy(
    linkToken.address,
    1,
    30,
    1,
    ethers.utils.formatBytes32String('ETH/USD'),
  )
  await fluxAggregator.deployed()
  await changePriceFeed(EA_1_URL, 100)
  await changePriceFeed(EA_2_URL, 120)
})

describe('FluxMonitor / FluxAggregator integration with one node', () => {
  let job: JobSpec

  afterEach(async () => {
    clClient1.archiveJob(job.id)
  })

  it('updates the price', async () => {
    await fluxAggregator
      .addOracle(node1Address, node1Address, 1, 1, 0)
      .then(txWait)
    await linkToken.transfer(fluxAggregator.address, deposit).then(txWait)
    await fluxAggregator.updateAvailableFunds().then(txWait)

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
    fluxMonitorJob.initiators[0].params.feeds = [EA_1_URL]
    job = clClient1.createJob(JSON.stringify(fluxMonitorJob))
    assert.equal(clClient1.getJobs().length, initialJobCount + 1)

    // Job should trigger initial FM run
    await assertJobRun(
      clClient1,
      job.id,
      initialRunCount + 1,
      'initial job never run',
    )
    matchers.bigNum(10000, await fluxAggregator.latestAnswer())

    // Nominally change price feed
    await changePriceFeed(EA_1_URL, 101)
    await wait(10000)
    assert.equal(
      clClient1.getJobRuns().length,
      initialRunCount + 1,
      'Flux Monitor should not run job after nominal price deviation',
    )

    // Significantly change price feed
    await changePriceFeed(EA_1_URL, 110)
    await assertJobRun(
      clClient1,
      job.id,
      initialRunCount + 2,
      'second job never run',
    )
    matchers.bigNum(11000, await fluxAggregator.latestAnswer())
  })
})

describe('FluxMonitor / FluxAggregator integration with two nodes', () => {
  // let job: JobSpec

  it('updates the price', async () => {
    await fluxAggregator
      .addOracle(node1Address, node1Address, 1, 1, 0)
      .then(txWait)
    await fluxAggregator
      .addOracle(node2Address, node2Address, 2, 2, 0)
      .then(txWait)
    await linkToken.transfer(fluxAggregator.address, deposit).then(txWait)
    await fluxAggregator.updateAvailableFunds().then(txWait)

    expect(await fluxAggregator.getOracles()).toEqual([
      node1Address,
      node2Address,
    ])
    matchers.bigNum(
      await linkToken.balanceOf(fluxAggregator.address),
      deposit,
      'Unable to fund FluxAggregator',
    )

    const node1InitialJobCount = clClient1.getJobs().length
    const node1InitialRunCount = clClient1.getJobRuns().length
    const node2InitialJobCount = clClient2.getJobs().length
    const node2InitialRunCount = clClient2.getJobRuns().length

    // TODO reset flux monitor job b/t tests (re-read from file?)
    fluxMonitorJob.initiators[0].params.address = fluxAggregator.address
    fluxMonitorJob.initiators[0].params.feeds = [EA_1_URL]
    const job1 = clClient1.createJob(JSON.stringify(fluxMonitorJob))
    fluxMonitorJob.initiators[0].params.feeds = [EA_2_URL]
    const job2 = clClient2.createJob(JSON.stringify(fluxMonitorJob))

    assert.equal(clClient1.getJobs().length, node1InitialJobCount + 1)
    assert.equal(clClient2.getJobs().length, node2InitialJobCount + 1)

    await assertJobRun(
      clClient1,
      job1.id,
      node1InitialRunCount + 1,
      'initial update never run by node 1',
    )
    await assertJobRun(
      clClient2,
      job2.id,
      node2InitialRunCount + 1,
      'initial update never run by node 2',
    )

    matchers.bigNum(11000, await fluxAggregator.latestAnswer())
    matchers.bigNum(1, await fluxAggregator.latestRound())
    matchers.bigNum(1, await fluxAggregator.reportingRound())
    matchers.bigNum(1, (await fluxAggregator.latestSubmission(node1Address))[1])
    matchers.bigNum(1, (await fluxAggregator.latestSubmission(node2Address))[1])
    matchers.bigNum(1, await fluxAggregator.withdrawable(node1Address))
    matchers.bigNum(1, await fluxAggregator.withdrawable(node2Address))

    clClient2.pause()
    await changePriceFeed(EA_1_URL, 110)
    await changePriceFeed(EA_1_URL, 120)

    await assertJobRun(
      clClient1,
      job1.id,
      node1InitialRunCount + 2,
      "node 1's second update not run",
    )
    // await wait(5000) // allow for node 2 to potentially respond
    matchers.bigNum(11000, await fluxAggregator.latestAnswer())
    matchers.bigNum(1, await fluxAggregator.latestRound())
    matchers.bigNum(2, await fluxAggregator.reportingRound())
    matchers.bigNum(2, (await fluxAggregator.latestSubmission(node1Address))[1])
    matchers.bigNum(1, (await fluxAggregator.latestSubmission(node2Address))[1])
    matchers.bigNum(
      2,
      await fluxAggregator.withdrawable(node1Address),
      "node 1 wasn't paid",
    )
    matchers.bigNum(
      1,
      await fluxAggregator.withdrawable(node2Address),
      "node 2 was paid and shouldn't have been",
    )

    clClient2.unpause()

    await assertJobRun(
      clClient2,
      job2.id,
      node2InitialRunCount + 2,
      "node 2's second update not run",
    )
    matchers.bigNum(11500, await fluxAggregator.latestAnswer())
    matchers.bigNum(2, await fluxAggregator.latestRound())
    matchers.bigNum(2, await fluxAggregator.reportingRound())
    matchers.bigNum(2, (await fluxAggregator.latestSubmission(node1Address))[1])
    matchers.bigNum(2, (await fluxAggregator.latestSubmission(node2Address))[1])
    matchers.bigNum(
      2,
      await fluxAggregator.withdrawable(node1Address),
      "node 1 wasn't paid",
    )
    matchers.bigNum(
      2,
      await fluxAggregator.withdrawable(node2Address),
      "node 2 wasn't paid",
    )
  })
})
