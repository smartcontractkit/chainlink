import * as jsonapi from '@chainlink/json-api-client'
import { Dispatch } from 'redux'
import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider } from 'ethers/providers'
import { FeedConfig } from 'config'
import * as actions from './actions'
import {
  formatAnswer,
  createContract,
  createInfuraProvider,
} from '../../../contracts/utils'

/**
 * feeds
 */
export function fetchFeeds() {
  return async (dispatch: Dispatch) => {
    dispatch(actions.fetchFeedsBegin())
    jsonapi
      .fetchWithTimeout('/feeds.json', {})
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
    const provider = createInfuraProvider()
    const payload = await latestAnswer(config, provider)
    const answer = formatAnswer(payload, config.multiply, config.decimalPlaces)
    const listingAnswer: actions.ListingAnswer = { answer, config }

    dispatch(actions.fetchAnswerSuccess(listingAnswer))
  }
}

const LATEST_ANSWER_CONTRACT_VERSION = 2

async function latestAnswer(
  contractConfig: FeedConfig,
  provider: JsonRpcProvider,
) {
  const contract = answerContract(contractConfig.contractAddress, provider)
  return contractConfig.contractVersion === LATEST_ANSWER_CONTRACT_VERSION
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
interface HealthPrice {
  config: any
  price: number
}

export function fetchHealthStatus(feed: FeedConfig) {
  return async (dispatch: Dispatch) => {
    const priceResponse = await fetchHealthPrice(feed)

    if (priceResponse) {
      dispatch(actions.fetchHealthPriceSuccess(priceResponse))
    }
  }
}

async function fetchHealthPrice(config: any): Promise<HealthPrice | undefined> {
  if (!config.healthPrice) return

  const json = await fetch(config.healthPrice).then(r => r.json())
  return { config, price: json[0].current_price }
}
