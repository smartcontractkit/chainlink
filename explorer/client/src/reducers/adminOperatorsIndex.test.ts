import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../reducers'
import { FetchAdminOperatorsSucceededAction } from '../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperatorsIndex: { items: ['replace-me'] },
}

describe('reducers/adminOperatorsIndex', () => {
  describe('FETCH_ADMIN_OPERATORS_SUCCEEDED', () => {
    it('replaces items', () => {
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
    })
  })
})
