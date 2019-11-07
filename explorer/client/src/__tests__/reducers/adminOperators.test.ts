import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import { FetchAdminOperatorsSucceededAction } from '../../reducers/actions'
import { mockPartial } from '../support/mocks'
import { ChainlinkNode } from 'explorer/models'

const ADMIN_OPERATOR_ID = 5555555
const INITIAL_ADMIN_OPERATOR = mockPartial<ChainlinkNode>({
  id: ADMIN_OPERATOR_ID,
})

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperators: {
    items: { [ADMIN_OPERATOR_ID]: INITIAL_ADMIN_OPERATOR },
  },
}

describe('reducers/adminOperators', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchAdminOperatorsSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminOperators).toEqual(INITIAL_STATE.adminOperators)
  })

  describe('FETCH_ADMIN_OPERATORS_SUCCEEDED', () => {
    it('can replace items', () => {
      const normalizedChainlinkNodes = {
        '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
          id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
        },
      }
      const orderedChainlinkNodes = [
        { id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e' },
      ]
      const data = {
        chainlinkNodes: normalizedChainlinkNodes,
        meta: {
          currentPageOperators: {
            data: orderedChainlinkNodes,
            meta: {},
          },
        },
      }
      const action = {
        type: 'FETCH_ADMIN_OPERATORS_SUCCEEDED',
        data,
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.adminOperators).toEqual({
        items: {
          '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e': {
            id: '9b7d791a-9a1f-4c55-a6be-b4231cf9fd4e',
          },
        },
      })
    })
  })
})
