import { partialAsFull } from '@chainlink/ts-test-helpers'
import { ChainlinkNode } from 'explorer/models'
import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import { FetchAdminOperatorsSucceededAction } from '../../reducers/actions'

const ADMIN_OPERATOR_ID = 5555555
const INITIAL_ADMIN_OPERATOR = partialAsFull<ChainlinkNode>({
  id: ADMIN_OPERATOR_ID,
})

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperators: {
    items: { [ADMIN_OPERATOR_ID]: INITIAL_ADMIN_OPERATOR },
  },
}

describe('reducers/adminOperators', () => {
  describe('FETCH_ADMIN_OPERATORS_SUCCEEDED', () => {
    it('can replace items', () => {
      const action: FetchAdminOperatorsSucceededAction = {
        type: 'FETCH_ADMIN_OPERATORS_SUCCEEDED',
        data: {
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
        },
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
