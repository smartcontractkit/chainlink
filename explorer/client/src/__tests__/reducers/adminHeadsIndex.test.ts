import { partialAsFull } from '@chainlink/ts-helpers'
import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import {
  FetchAdminHeadsSucceededAction,
  FetchAdminHeadsErrorAction,
} from '../../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminHeadsIndex: { items: ['replace-me'], loaded: false },
}

describe('reducers/adminHeadsIndex', () => {
  it('FETCH_ADMIN_HEADS_SUCCEEDED can replace items', () => {
    const action: FetchAdminHeadsSucceededAction = {
      type: 'FETCH_ADMIN_HEADS_SUCCEEDED',
      data: {
        heads: [],
        meta: {
          currentPageHeads: {
            data: [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }],
            meta: { count: 100 },
          },
        },
      },
    }
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminHeadsIndex.items).toEqual([
      '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
    ])
    expect(state.adminHeadsIndex.count).toEqual(100)
    expect(state.adminHeadsIndex.loaded).toEqual(true)
  })

  it('FETCH_ADMIN_HEADS_ERROR sets loaded', () => {
    const action = partialAsFull<FetchAdminHeadsErrorAction>({
      type: 'FETCH_ADMIN_HEADS_ERROR',
      error: new Error(),
    })
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminHeadsIndex.loaded).toEqual(true)
  })
})
