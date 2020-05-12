import { partialAsFull } from '@chainlink/ts-helpers'
import { Head } from 'explorer/models'
import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../reducers'
import {
  AdminHeadsNormalizedData,
  FetchAdminHeadsBeginAction,
  FetchAdminHeadsSucceededAction,
  FetchAdminHeadsErrorAction,
} from '../reducers/actions'

const INITIAL_ADMIN_HEAD = partialAsFull<Head>({ id: 1 })

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminHeads: {
    loading: false,
    error: false,
    items: { [1]: INITIAL_ADMIN_HEAD },
  },
}

describe('reducers/adminHeads', () => {
  const fetchAdminHeadsBeginAction: FetchAdminHeadsBeginAction = {
    type: 'FETCH_ADMIN_HEADS_BEGIN',
  }

  describe('FETCH_ADMIN_HEADS_BEGIN', () => {
    it('sets loading to true', () => {
      const state = reducer(INITIAL_STATE, fetchAdminHeadsBeginAction)

      expect(state.adminHeads.loading).toEqual(true)
    })
  })

  describe('FETCH_ADMIN_HEADS_SUCCEEDED', () => {
    const data = partialAsFull<AdminHeadsNormalizedData>({
      heads: {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      },
      meta: {
        currentPageHeads: {
          data: [{ id: 1 }],
          meta: { count: 100 },
        },
      },
    })
    const action: FetchAdminHeadsSucceededAction = {
      type: 'FETCH_ADMIN_HEADS_SUCCEEDED',
      data,
    }

    it('replaces items', () => {
      const state = reducer(INITIAL_STATE, action)

      expect(state.adminHeads.items).toEqual({
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      })
    })

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchAdminHeadsBeginAction)
      expect(state.adminHeads.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.adminHeads.loading).toEqual(false)
    })
  })

  describe('FETCH_ADMIN_HEADS_ERROR', () => {
    const action: FetchAdminHeadsErrorAction = {
      type: 'FETCH_ADMIN_HEADS_ERROR',
      errors: [
        {
          status: 500,
          detail: 'An error',
        },
      ],
    }

    it('sets loading to false', () => {
      let state = reducer(INITIAL_STATE, fetchAdminHeadsBeginAction)
      expect(state.adminHeads.loading).toEqual(true)

      state = reducer(state, action)
      expect(state.adminHeads.loading).toEqual(false)
    })

    it('sets error to true & resets it to false when a new fetch starts', () => {
      let state = reducer(INITIAL_STATE, action)
      expect(state.adminHeads.error).toEqual(true)

      state = reducer(state, fetchAdminHeadsBeginAction)
      expect(state.adminHeads.error).toEqual(false)
    })
  })
})
