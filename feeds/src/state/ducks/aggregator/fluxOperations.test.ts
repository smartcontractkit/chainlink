import configureStore from 'redux-mock-store'
import thunk from 'redux-thunk'
import FluxOperations from './fluxOperations'
import { INITIAL_STATE } from './reducers'

const middlewares = [thunk]
const mockStore = configureStore(middlewares)
const store = mockStore(INITIAL_STATE)

let mockState = {}
let mockSubmissionEventLog = {}

jest.spyOn(store, 'getState').mockImplementation(() => {
  return mockState
})

const dispatchWrapper = (f: any) => (...args: any[]) => {
  return f(...args)(store.dispatch, store.getState)
}

const contractInstance: any = {
  listenSubmissionReceivedEvent: jest.fn(callback =>
    callback(mockSubmissionEventLog),
  ),
  listenNewRoundEvent: jest.fn(),
  latestAnswer: jest.fn(),
  latestTimestamp: jest.fn(),
}

describe('state/ducks/aggregator/fluxOperations', () => {
  describe('fetchAnswers', () => {
    beforeEach(() => {
      store.clearActions()
      jest.clearAllMocks()
      FluxOperations.contractInstance = contractInstance
    })

    it('should update answer from log event', async () => {
      mockState = {
        aggregator: {
          minimumAnswers: 2,
          oracleAnswers: [
            {
              sender: '0x1',
              answerId: 1,
            },
          ],
        },
      }

      mockSubmissionEventLog = { answerId: 2, sender: '0x1' }

      await dispatchWrapper(FluxOperations.initListeners)()
      const actions = store.getActions()[0]
      expect(actions.type).toEqual('aggregator/ORACLE_ANSWERS')
      expect(actions.payload).toEqual([mockSubmissionEventLog])
    })

    it('should update answer from log event', async () => {
      mockState = {
        aggregator: {
          minimumAnswers: 2,
          oracleAnswers: [
            {
              sender: '0x1',
              answerId: 1,
            },
            {
              sender: '0x2',
              answerId: 2,
            },
          ],
        },
      }

      mockSubmissionEventLog = { answerId: 2, sender: '0x1' }

      await dispatchWrapper(FluxOperations.initListeners)()
      const actions = store.getActions()[0]
      expect(actions.type).toEqual('aggregator/ORACLE_ANSWERS')
      expect(actions.payload).toEqual([
        {
          sender: '0x1',
          answerId: 2,
        },
        {
          sender: '0x2',
          answerId: 2,
        },
      ])
    })

    it('should fetch latest answer and latest timestamp ', async () => {
      mockState = {
        aggregator: {
          minimumAnswers: 1,
          oracleAnswers: [
            {
              sender: '0x1',
              answerId: 1,
            },
          ],
        },
      }

      mockSubmissionEventLog = { answerId: 2, sender: '0x1' }

      await dispatchWrapper(FluxOperations.initListeners)()
      const actions = store.getActions()
      expect(actions[0].type).toEqual('aggregator/ORACLE_ANSWERS')
      expect(actions[1].type).toEqual('aggregator/LATEST_ANSWER')
      expect(actions[2].type).toEqual('aggregator/LATEST_ANSWER_TIMESTAMP')
    })

    it('should not fetch latest answer and latest timestamp ', async () => {
      mockState = {
        aggregator: {
          minimumAnswers: 2,
          oracleAnswers: [
            {
              sender: '0x1',
              answerId: 1,
            },
          ],
        },
      }

      mockSubmissionEventLog = { answerId: 2, sender: '0x1' }

      await dispatchWrapper(FluxOperations.initListeners)()
      const actions = store.getActions()
      expect(actions[0].type).toEqual('aggregator/ORACLE_ANSWERS')
      expect(actions.length).toEqual(1)
    })
  })
})
