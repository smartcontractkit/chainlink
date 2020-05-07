import * as jsonapi from '@chainlink/json-api-client'
import { Dispatch } from 'redux'
import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider } from 'ethers/providers'
import { FeedConfig, Config } from '../../../config'
import * as actions from './actions'
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
    dispatch(actions.fetchFeedsBegin())
    jsonapi
      .fetchWithTimeout(Config.feedsJson(), {})
      .then((r: Response) => r.json())
      .then((json: FeedConfig[]) => {
        dispatch(actions.fetchFeedsSuccess(json))
      })
      .catch(e => {
        dispatch(actions.fetchFeedsError(e))
      })
  }
}

/**
 * answers
 */
export function fetchAnswer(config: FeedConfig) {
  return async (dispatch: Dispatch) => {
    try {
      const provider = createInfuraProvider()
      const payload = await latestAnswer(config, provider)
      const answer = formatAnswer(
        payload,
        config.multiply,
        config.decimalPlaces,
      )
      const listingAnswer: actions.ListingAnswer = { answer, config }

      dispatch(actions.fetchAnswerSuccess(listingAnswer))
    } catch {
      console.error('Could not fetch answer')
    }
  }
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
      dispatch(actions.fetchHealthPriceSuccess(priceResponse))
    }
  }
}

async function fetchHealthPrice(
  config: any,
): Promise<actions.HealthPrice | undefined> {
  if (!config.healthPrice) return

  const json = await fetch(config.healthPrice).then(r => r.json())
  return { config, price: json[0].current_price }
}
