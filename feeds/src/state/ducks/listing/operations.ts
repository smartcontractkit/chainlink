import { FeedConfig, getFeedsConfig } from 'config'
import { JsonRpcProvider } from 'ethers/providers'
import { FunctionFragment } from 'ethers/utils'
import { Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk'
import { AppState } from 'state'
import { Actions } from 'state/actions'
import {
  createContract,
  createInfuraProvider,
  formatAnswer,
} from '../../../contracts/utils'
import * as actions from './actions'
import { ListingAnswer } from './reducers'
import { ListingGroup } from './selectors'

interface HealthPrice {
  config: any
  price: number
}

export function fetchHealthStatus(groups: ListingGroup[]) {
  return async (dispatch: Dispatch) => {
    const configs = groups.flatMap(g => g.feeds)
    const priceResponses = await Promise.all(configs.map(fetchHealthPrice))

    priceResponses
      .filter(pr => pr)
      .forEach(pr => {
        dispatch(actions.setHealthPrice(pr))
      })
  }
}

async function fetchHealthPrice(config: any): Promise<HealthPrice | undefined> {
  if (!config.healthPrice) return

  const json = await fetch(config.healthPrice).then(r => r.json())
  return { config, price: json[0].current_price }
}

export function fetchAnswers(): ThunkAction<void, AppState, void, Actions> {
  return async dispatch => {
    const provider = createInfuraProvider()
    const answerList = await allAnswers(provider)
    dispatch(actions.setAnswers(answerList))
  }
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

const LATEST_ANSWER_CONTRACT_VERSION = 2

async function latestAnswer(
  contractConfig: FeedConfig,
  provider: JsonRpcProvider,
) {
  const contract = answerContract(contractConfig.contractAddress, provider)
  return contractConfig.contractVersion === LATEST_ANSWER_CONTRACT_VERSION ||
    contractConfig.contractVersion === 3
    ? await contract.latestAnswer()
    : await contract.currentAnswer()
}

async function allAnswers(provider: JsonRpcProvider): Promise<ListingAnswer[]> {
  const answers = getFeedsConfig()
    .filter(config => config.listing)
    .map(async config => {
      try {
        const payload = await latestAnswer(config, provider)
        const answer = formatAnswer(
          payload,
          config.multiply,
          config.decimalPlaces,
        )
        return { answer, config }
      } catch {
        console.log('Could not fetch the answer')
      }
      return { config }
    })

  return Promise.all(answers)
}
