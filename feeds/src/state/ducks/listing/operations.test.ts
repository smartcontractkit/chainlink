import { partialAsFull } from '@chainlink/ts-helpers'
import { getFeedsConfig } from 'config'
import { Contract } from 'ethers'
import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import * as utils from '../../../contracts/utils'
import { Networks } from '../../../utils'
import * as operations from './operations'
import { INITIAL_STATE } from './reducers'

jest.mock('../../../contracts/utils')

const formatAnswerSpy = jest
  .spyOn(utils, 'formatAnswer')
  .mockImplementation(answer => answer)

const createContractSpy = jest
  .spyOn(utils, 'createContract')
  .mockImplementation(() => {
    const contract = partialAsFull<Contract>({
      latestAnswer: () => 'latestAnswer',
      currentAnswer: () => 'currentAnswer',
    })
    return contract
  })

const mainnetContracts = getFeedsConfig().filter(
  config => config.networkId === Networks.MAINNET,
)

const middlewares = [thunk]
const mockStore = configureStore(middlewares)
const store = mockStore(INITIAL_STATE)

const dispatchWrapper = (f: any) => (...args: any[]) => {
  return f(...args)(store.dispatch, store.getState)
}

describe('state/ducks/listing', () => {
  describe('fetchAnswers', () => {
    beforeEach(() => {
      store.clearActions()
      jest.clearAllMocks()
    })

    it('should fetch answer list', async () => {
      await dispatchWrapper(operations.fetchAnswers)()
      const actions = store.getActions()[0]
      expect(actions.type).toEqual('listing/SET_ANSWERS')
      expect(actions.payload).toHaveLength(mainnetContracts.length)

      const contractVersionOne = actions.payload.filter(
        (answer: any) => answer.config.contractVersion === 1,
      )[0]

      const contractVersionTwo = actions.payload.filter(
        (answer: any) => answer.config.contractVersion === 2,
      )[0]

      expect(contractVersionOne.answer).toEqual('currentAnswer')
      expect(contractVersionTwo.answer).toEqual('latestAnswer')
    })

    it('should build a list of objects', async () => {
      await dispatchWrapper(operations.fetchAnswers)()
      const actions = store.getActions()[0]
      expect(actions.payload[0]).toHaveProperty('answer')
      expect(actions.payload[0]).toHaveProperty('config')
    })

    it('should format answers', async () => {
      await dispatchWrapper(operations.fetchAnswers)()
      expect(formatAnswerSpy).toHaveBeenCalledTimes(mainnetContracts.length)
    })

    it('should create a contracts for each config', async () => {
      await dispatchWrapper(operations.fetchAnswers)()
      expect(createContractSpy).toHaveBeenCalledTimes(mainnetContracts.length)
    })
  })
})
