import * as jsonapi from '@chainlink/json-api-client'
import { Dispatch } from 'redux'
import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider } from 'ethers/providers'
import { FeedConfig, Config } from '../../../config'
import {
  fetchFeedsBegin,
  fetchFeedsSuccess,
  fetchFeedsError,
  fetchAnswerSuccess,
  fetchHealthPriceSuccess,
  fetchAnswerTimestampSuccess,
} from './actions'
import { ListingAnswer, HealthPrice, ListingTimestamp } from './types'
import {
  createContract,
  createInfuraProvider,
  formatAnswer,
} from '../../../contracts/utils'

/**
 * feeds
 */
export function fetchFeeds() {
  return async (dispatch: Dispatch) => {
    dispatch(fetchFeedsBegin())
    jsonapi
      .fetchWithTimeout(Config.feedsJson(), {})
      .then((r: Response) => r.json())
      .then((json: FeedConfig[]) => {
        dispatch(fetchFeedsSuccess(json))
      })
      .catch(e => {
        dispatch(fetchFeedsError(e))
      })
  }
}

/**
 * answers
 */
export function fetchLatestData(config: FeedConfig) {
  return async (dispatch: Dispatch) => {
    try {
      const provider = createInfuraProvider()
      const answerPayload = await latestAnswer(config, provider)
      const timestampPayload = await latestTimestamp(config, provider)

      const answer = formatAnswer(
        answerPayload,
        config.multiply,
        config.decimalPlaces,
        config.formatDecimalPlaces,
      )
      const listingAnswer: ListingAnswer = { answer, config }
      const listingAnswerTimestamp: ListingTimestamp = {
        timestamp: Number(timestampPayload),
        config,
      }

      dispatch(fetchAnswerSuccess(listingAnswer))
      dispatch(fetchAnswerTimestampSuccess(listingAnswerTimestamp))
    } catch {
      console.error('Could not fetch answer')
    }
  }
}

async function latestTimestamp(
  contractConfig: FeedConfig,
  provider: JsonRpcProvider,
) {
  const contract = answerContract(contractConfig.contractAddress, provider)
  return await contract.latestTimestamp()
}

const LATEST_ANSWER_CONTRACT_VERSIONS = [2, 3]

async function latestAnswer(
  contractConfig: FeedConfig,
  provider: JsonRpcProvider,
) {
  const contract = answerContract(contractConfig.contractAddress, provider)
  return LATEST_ANSWER_CONTRACT_VERSIONS.includes(
    contractConfig.contractVersion,
  )
    ? await contract.latestAnswer()
    : await contract.currentAnswer()
}

const ANSWER_ABI: FunctionFragment[] = [
  {
    constant: true,
    inputs: [],
    name: 'currentAnswer',
    outputs: [{ name: '', type: 'int256' }],
    payable: false,
    stateMutability: 'view',
    type: 'function',
  },
  {
    constant: true,
    inputs: [],
    name: 'latestAnswer',
    outputs: [{ name: '', type: 'int256' }],
    payable: false,
    stateMutability: 'view',
    type: 'function',
  },
  {
    constant: true,
    inputs: [],
    name: 'latestTimestamp',
    outputs: [{ name: '', type: 'uint256' }],
    payable: false,
    stateMutability: 'view',
    type: 'function',
  },
]

function answerContract(contractAddress: string, provider: JsonRpcProvider) {
  return createContract(contractAddress, provider, ANSWER_ABI)
}

/**
 * health checks
 */

export function fetchHealthStatus(feed: FeedConfig) {
  return async (dispatch: Dispatch) => {
    const priceResponse = await fetchHealthPrice(feed)

    if (priceResponse) {
      dispatch(fetchHealthPriceSuccess(priceResponse))
    }
  }
}

async function fetchHealthPrice(config: any): Promise<HealthPrice | undefined> {
  if (!config.healthPrice) return

  const json = await fetch(config.healthPrice).then(r => r.json())
  return { config, price: json[0].current_price }
}
