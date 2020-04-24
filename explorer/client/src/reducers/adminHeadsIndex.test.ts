import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../reducers'
import { FetchAdminHeadsSucceededAction } from '../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminHeadsIndex: { items: ['replace-me'] },
}

describe('reducers/adminHeadsIndex', () => {
  describe('FETCH_ADMIN_HEADS_SUCCEEDED', () => {
    it('replaces items', () => {
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
    })
  })
})
