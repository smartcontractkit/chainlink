import { partialAsFull } from '@chainlink/ts-helpers'
import { ChainlinkNode } from 'explorer/models'
import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../reducers'
import {
  AdminOperatorsNormalizedData,
  FetchAdminOperatorsBeginAction,
  FetchAdminOperatorsSucceededAction,
  FetchAdminOperatorsErrorAction,
} from '../reducers/actions'

const ADMIN_OPERATOR_ID = 5555555
const INITIAL_ADMIN_OPERATOR = partialAsFull<ChainlinkNode>({
  id: ADMIN_OPERATOR_ID,
})

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperators: {
    loading: false,
    error: false,
    items: { [ADMIN_OPERATOR_ID]: INITIAL_ADMIN_OPERATOR },
  },
}

describe('reducers/adminOperators', () => {
  const fetchAdminOperatorsBeginAction: FetchAdminOperatorsBeginAction = {
    type: 'FETCH_ADMIN_OPERATORS_BEGIN',
  }

  describe('FETCH_ADMIN_OPERATORS_BEGIN', () => {
    it('sets loading to true', () => {
      const state = reducer(INITIAL_STATE, fetchAdminOperatorsBeginAction)

      expect(state.adminOperators.loading).toEqual(true)
    })
  })

  describe('FETCH_ADMIN_OPERATORS_SUCCEEDED', () => {
    const data = partialAsFull<AdminOperatorsNormalizedData>({
      chainlinkNodes: {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      },
      meta: {
        currentPageOperators: {
          data: [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }],
          meta: {
            count: 1,
          },
        },
      },
    })
    const action: FetchAdminOperatorsSucceededAction = {
      type: 'FETCH_ADMIN_OPERATORS_SUCCEEDED',
      data,
    }

    it('replaces items', () => {
      const state = reducer(INITIAL_STATE, action)

      expect(state.adminOperators.items).toEqual({
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      })
    })

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchAdminOperatorsBeginAction)
      expect(state.adminOperators.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.adminOperators.loading).toEqual(false)
    })
  })

  describe('FETCH_ADMIN_OPERATORS_ERROR', () => {
    const action: FetchAdminOperatorsErrorAction = {
      type: 'FETCH_ADMIN_OPERATORS_ERROR',
      errors: [
        {
          status: 500,
          detail: 'An error',
        },
      ],
    }

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchAdminOperatorsBeginAction)
      expect(state.adminOperators.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.adminOperators.loading).toEqual(false)
    })

    it('sets error to true & resets it to false when a new fetch starts', () => {
      let state = reducer(INITIAL_STATE, action)
      expect(state.adminOperators.error).toEqual(true)

      state = reducer(state, fetchAdminOperatorsBeginAction)
      expect(state.adminOperators.error).toEqual(false)
    })
  })
})
