import * as actions from './actions'
import {
  formatAnswer,
  createContract,
  createInfuraProvider,
} from 'contracts/utils'
import feeds from 'feeds.json'
import { MAINNET_ID } from 'utils/'

const answerContract = (contractAddress, provider) => {
  return createContract(contractAddress, provider, [
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
  ])
}

const latestAnswer = async (contractConfig, provider) => {
  const contract = answerContract(contractConfig.contractAddress, provider)
  return contractConfig.contractVersion === 2
    ? await contract.latestAnswer()
    : await contract.currentAnswer()
}

const allAnswers = async provider => {
  const answers = feeds
    .filter(config => config.networkId === MAINNET_ID)
    .map(async config => {
      const payload = await latestAnswer(config, provider)
      return {
        answer: formatAnswer(payload, config.multiply, config.decimalPlaces),
        config,
      }
    })

  return Promise.all(answers)
}

const fetchAnswers = () => {
  return async dispatch => {
    const provider = createInfuraProvider()
    const answerList = await allAnswers(provider)

    dispatch(actions.setAnswers(answerList))
  }
}

const fetchHealthPrice = async config => {
  if (config.health_price) {
    const json = await fetch(config.health_price).then(r => r.json())
    return [config, json[0].current_price]
  }
}

const fetchHealthStatus = groups => {
  return async dispatch => {
    const configs = groups.flatMap(g => g.list.map(l => l.config))
    const priceResponses = await Promise.all(configs.map(fetchHealthPrice))

    priceResponses
      .filter(pr => pr)
      .forEach(pr => {
        dispatch(actions.setHealthPrice(pr))
      })
  }
}

export { fetchAnswers, fetchHealthStatus }
