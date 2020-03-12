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
  getArgs,
  wait,
} from '../test-helpers/common'

const [node1URL, node2URL] = ['http://node:6688', 'http://node-2:6688']
const { EXTERNAL_ADAPTER_URL } = getArgs(['EXTERNAL_ADAPTER_URL'])

const provider = createProvider()
const carol = ethers.Wallet.createRandom().connect(provider)
const linkTokenFactory = new contract.LinkTokenFactory(carol)
const fluxAggregatorFactory = new FluxAggregatorFactory(carol)
const adapterURL = new URL('result', EXTERNAL_ADAPTER_URL).href
const deposit = h.toWei('1000')
const clClient1 = new ChainlinkClient(node1URL)
const clClient2 = new ChainlinkClient(node2URL)

console.log(node2URL)

let linkToken: contract.Instance<contract.LinkTokenFactory>
let node1Address: string
let node2Address: string

async function changePriceFeed(value: number) {
  const response = await fetch(adapterURL, {
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
  console.log('node1Address', node1Address)
  console.log('node2Address', node2Address)
  await fundAddress(carol.address)
  await fundAddress(node1Address)
  await fundAddress(node2Address)
  linkToken = await linkTokenFactory.deploy()
  await linkToken.deployed()
})

// describe('test', () => {
//   it('works', () => {
//     assert(true)
//   })
// })

describe('FluxMonitor / FluxAggregator integration with one node', () => {
  let fluxAggregator: contract.Instance<FluxAggregatorFactory>
  let job: JobSpec

  afterEach(async () => {
    clClient1.archiveJob(job.id)
    await changePriceFeed(100) // original price
  })

  it('updates the price with a single node', async () => {
    const clClient = new ChainlinkClient(node1URL)

    fluxAggregator = await fluxAggregatorFactory.deploy(
      linkToken.address,
      1,
      600,
      1,
      ethers.utils.formatBytes32String('ETH/USD'),
    )
    await fluxAggregator.deployed()

    const tx1 = await fluxAggregator.addOracle(
      node1Address,
      node1Address,
      1,
      1,
      0,
    )
    await tx1.wait()
    const tx2 = await linkToken.transfer(fluxAggregator.address, deposit)
    await tx2.wait()
    const tx3 = await fluxAggregator.updateAvailableFunds()
    await tx3.wait()

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
    job = clClient1.createJob(JSON.stringify(fluxMonitorJob))
    assert.equal(clClient1.getJobs().length, initialJobCount + 1)

    // Job should trigger initial FM run
    await assertJobRun(
      clClient,
      job.id,
      initialRunCount + 1,
      'initial job never run',
    )
    matchers.bigNum(10000, await fluxAggregator.latestAnswer())

    // Nominally change price feed
    await changePriceFeed(101)
    await wait(10000)
    assert.equal(
      clClient1.getJobRuns().length,
      initialRunCount + 1,
      'Flux Monitor should not run job after nominal price deviation',
    )

    // Significantly change price feed
    await changePriceFeed(110)
    await assertJobRun(
      clClient,
      job.id,
      initialRunCount + 2,
      'second job never run',
    )
    matchers.bigNum(11000, await fluxAggregator.latestAnswer())
  })
})

describe('FluxMonitor / FluxAggregator integration with two nodes', () => {
  it.skip('works', () => {
    assert(true)
  })
})
