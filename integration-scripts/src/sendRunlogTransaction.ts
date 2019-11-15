import { RunLogFactory } from './generated'
import { generated as chainlink } from 'chainlink'
import { RunLog } from './generated/RunLog'
import { ethers } from 'ethers'
import url from 'url'
import {
  DEVNET_ADDRESS,
  registerPromiseHandler,
  getArgs,
  credentials,
  createProvider,
} from './common'
const request = require('request-promise').defaults({ jar: true })

async function main() {
  registerPromiseHandler()
  const args = getArgs([
    'CHAINLINK_URL',
    'ECHO_SERVER_URL',
    'RUN_LOG_ADDRESS',
    'LINK_TOKEN_ADDRESS',
  ])

  await sendRunlogTransaction({
    chainlinkUrl: args.CHAINLINK_URL,
    echoServerUrl: args.ECHO_SERVER_URL,
    linkTokenAddress: args.LINK_TOKEN_ADDRESS,
    runLogAddress: args.RUN_LOG_ADDRESS,
  })
}
main()

interface Args {
  runLogAddress: string
  linkTokenAddress: string
  chainlinkUrl: string
  echoServerUrl: string
}
async function sendRunlogTransaction({
  runLogAddress,
  linkTokenAddress,
  chainlinkUrl,
  echoServerUrl,
}: Args) {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  const runLogFactory = new RunLogFactory(signer)
  const linkTokenFactory = new chainlink.LinkTokenFactory(signer)
  const runLog = runLogFactory.attach(runLogAddress)
  const linkToken = linkTokenFactory.attach(linkTokenAddress)

  // transfer link to runlog address
  const linkAmount = ethers.utils.parseEther('1000')
  try {
    await linkToken.transfer(runLog.address, linkAmount, {
      gasLimit: 100000,
    })
  } catch (error) {
    console.error('Error transferring link to RunLog')
    throw Error(error)
  }

  console.log(`Transferred ${linkAmount} to RunLog at: ${runLog.address}`)

  await signIn(chainlinkUrl)
  const job = await createJob(chainlinkUrl, runLog.address, echoServerUrl)
  await makeRunlogRequest(runLog, job)
}

/**
 * Sign into a chainlink node by creating a session
 * @param chainlinkUrl The chainlink node to send the signin request to
 */
async function signIn(chainlinkUrl: string) {
  const sessionsUrl = url.resolve(chainlinkUrl, '/sessions')
  await request.post(sessionsUrl, { json: credentials })
}

/**
 * Calculate the current date + offset in seconds
 *
 * @param seconds The number of seconds to add as an offset
 */
function futureOffsetSeconds(seconds: number): number {
  const nowInSeconds = Math.ceil(new Date().getTime() / 1000)

  return nowInSeconds + seconds
}

/**
 * Create a chainlink job
 *
 * @param requesterAddress The requester of the job to be created
 */
async function createJob(
  chainlinkUrl: string,
  requesterAddress: string,
  echoServerUrl: string,
) {
  const job = {
    initiators: [
      { type: 'runlog', params: { requesters: [requesterAddress] } },
    ],
    tasks: [
      // 10 seconds to ensure the time has not elapsed by the time the run is triggered
      { type: 'Sleep', params: { until: futureOffsetSeconds(10) } },
      { type: 'HttpPost', params: { url: echoServerUrl } },
      { type: 'EthTx' },
    ],
  }

  const specsUrl = url.resolve(chainlinkUrl, '/v2/specs')
  const Job = await request.post(specsUrl, { json: job }).catch((e: any) => {
    throw Error(`Error creating Job ${e}`)
  })
  return Job
}

async function makeRunlogRequest(runLog: RunLog, job: any) {
  try {
    await runLog.request(ethers.utils.toUtf8Bytes(job.data.id), {
      gasLimit: 2000000,
    })
  } catch (error) {
    console.error('Error making runlog request')
    throw error
  }
  console.log(`Made RunLog request`)
}
