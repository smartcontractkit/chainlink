import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider } from 'ethers/providers'
import { FeedConfig } from 'feeds'
import * as actions from './actions'
import {
  formatAnswer,
  createContract,
  createInfuraProvider,
} from '../../../contracts/utils'
import { Networks } from '../../../utils'
import feeds from '../../../feeds.json'

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
  return contractConfig.contractVersion === LATEST_ANSWER_CONTRACT_VERSION
    ? await contract.latestAnswer()
    : await contract.currentAnswer()
}

async function allAnswers(provider: JsonRpcProvider) {
  const answers = feeds
    .filter(config => config.networkId === Networks.MAINNET)
    .map(async config => {
      const payload = await latestAnswer(config, provider)
      const answer = formatAnswer(
        payload,
        config.multiply,
        config.decimalPlaces,
      )

      return { answer, config }
    })

  return Promise.all(answers)
}

export interface ListingAnswer {
  answer: string
  config: FeedConfig
}

function fetchAnswers() {
  return async (dispatch: any) => {
    const provider = createInfuraProvider()
    const answerList = await allAnswers(provider)
    dispatch(actions.setAnswers(answerList))
  }
}

export { fetchAnswers }
