import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../reducers'
import { OperatorShowData } from '../reducers/adminOperatorsShow'
import { FetchAdminOperatorSucceededAction } from '../reducers/actions'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminOperatorsShow: {},
}

const MOCK_PAYLOAD: OperatorShowData = {
  id: '1',
  name: 'name',
  url: 'url',
  createdAt: 'mm-yy',
  uptime: 1,
  jobCounts: {
    completed: 2,
    errored: 3,
    inProgress: 4,
    total: 9,
  },
}

const MOCK_ACTION: FetchAdminOperatorSucceededAction = {
  type: 'FETCH_ADMIN_OPERATOR_SUCCEEDED',
  data: {
    chainlinkNodes: {
      '1': MOCK_PAYLOAD,
    },
    meta: {
      node: {
        data: [],
      },
    },
  },
}

describe('reducers/adminOperatorsShow', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchAdminOperatorSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminOperatorsShow).toEqual(INITIAL_STATE.adminOperatorsShow)
  })

  describe('FETCH_ADMIN_OPERATOR_SUCCEEDED', () => {
    it('can receive new operator', () => {
      const state = reducer(INITIAL_STATE, MOCK_ACTION)
      expect(state.adminOperatorsShow).toEqual({ 1: MOCK_PAYLOAD })
    })
  })
})
