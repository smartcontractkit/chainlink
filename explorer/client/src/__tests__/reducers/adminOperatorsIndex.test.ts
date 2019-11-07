import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import { FetchAdminOperatorsSucceededAction } from '../../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperatorsIndex: { items: ['replace-me'] },
}

describe('reducers/adminOperatorsIndex', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchAdminOperatorsSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminOperatorsIndex).toEqual(INITIAL_STATE.adminOperatorsIndex)
  })

  describe('FETCH_ADMIN_OPERATORS_SUCCEEDED', () => {
    it('can replace items', () => {
      const operators = [{ id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' }]
      const data = {
        meta: {
          currentPageOperators: {
            data: operators,
            meta: { count: 100 },
          },
        },
      }
      const action = {
        type: 'FETCH_ADMIN_OPERATORS_SUCCEEDED',
        data: data,
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.adminOperatorsIndex).toEqual({
        items: ['9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e'],
        count: 100,
      })
    })
  })
})
