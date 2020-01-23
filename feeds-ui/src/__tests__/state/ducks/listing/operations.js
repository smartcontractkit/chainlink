import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import * as types from 'state/ducks/listing/types'
import { initialState } from 'state/ducks/listing/reducers'
import * as operations from 'state/ducks/listing/operations'
import feeds from 'feeds.json'
import { MAINNET_ID } from 'utils'

const mainnetContracts = feeds.filter(config => config.networkId === MAINNET_ID)

import * as utils from 'contracts/utils'

const middlewares = [thunk]
const mockStore = configureStore(middlewares)
const store = mockStore(initialState)

jest.mock('contracts/utils')

const dispatchWrapper = f => (...args) => {
  return f(...args)(store.dispatch, store.getState)
}

const formatAnswerSpy = jest
  .spyOn(utils, 'formatAnswer')
  .mockImplementation(answer => answer)

const createContractSpy = jest
  .spyOn(utils, 'createContract')
  .mockImplementation(() => ({
    latestAnswer: () => 'latestAnswer',
    currentAnswer: () => 'currentAnswer',
  }))

describe('state/ducks/listing', () => {
  describe('fetchAnswers', () => {
    beforeEach(() => {
      store.clearActions()
      jest.clearAllMocks()
    })

    it('should fetch answer list', async () => {
      await dispatchWrapper(operations.fetchAnswers)()
      const actions = store.getActions()[0]
      expect(actions.type).toEqual(types.SET_ANSWERS)
      expect(actions.payload).toHaveLength(mainnetContracts.length)

      const contractVersionOne = actions.payload.filter(
        answer => answer.config.contractVersion === 1,
      )[0]

      const contractVersionTwo = actions.payload.filter(
        answer => answer.config.contractVersion === 2,
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
