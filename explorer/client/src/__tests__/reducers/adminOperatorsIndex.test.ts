import { partialAsFull } from '@chainlink/ts-helpers'
import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import {
  FetchAdminOperatorsSucceededAction,
  FetchAdminOperatorsErrorAction,
} from '../../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperatorsIndex: { items: ['replace-me'], loaded: false },
}

describe('reducers/adminOperatorsIndex', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchAdminOperatorsSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminOperatorsIndex).toEqual(INITIAL_STATE.adminOperatorsIndex)
  })

  it('FETCH_ADMIN_OPERATORS_SUCCEEDED can replace items', () => {
    const action: FetchAdminOperatorsSucceededAction = {
      type: 'FETCH_ADMIN_OPERATORS_SUCCEEDED',
      data: {
        chainlinkNodes: [],
        meta: {
          currentPageOperators: {
            data: [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }],
            meta: { count: 100 },
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminOperatorsIndex.items).toEqual([
      '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
    ])
    expect(state.adminOperatorsIndex.count).toEqual(100)
    expect(state.adminOperatorsIndex.loaded).toEqual(true)
  })

  it('FETCH_ADMIN_OPERATORS_ERROR sets loaded', () => {
    const action = partialAsFull<FetchAdminOperatorsErrorAction>({
      type: 'FETCH_ADMIN_OPERATORS_ERROR',
      error: new Error(),
    })
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminOperatorsIndex.loaded).toEqual(true)
  })
})
