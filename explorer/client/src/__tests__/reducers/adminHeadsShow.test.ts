import reducer, {
  INITIAL_STATE as initialRootState,
  AppState,
} from '../../reducers'
import { HeadShowData } from '../../reducers/adminHeadsShow'
import { FetchAdminHeadSucceededAction } from '../../reducers/actions'
import { BigNumber } from 'bignumber.js'

const INITIAL_STATE: AppState = {
  ...initialRootState,
  adminHeadsShow: {},
}

const MOCK_PAYLOAD: HeadShowData = {
  id: 1,
  parentHash: Buffer.from([]),
  uncleHash: Buffer.from([]),
  coinbase: Buffer.from([]),
  root: Buffer.from([]),
  txHash: Buffer.from([]),
  receiptHash: Buffer.from([]),
  bloom: Buffer.from([]),
  difficulty: new BigNumber(0),
  number: new BigNumber(0),
  gasLimit: new BigNumber(0),
  gasUsed: new BigNumber(0),
  time: new BigNumber(0),
  extra: Buffer.from([]),
  mixDigest: Buffer.from([]),
  nonce: Buffer.from([]),
}

const MOCK_ACTION: FetchAdminHeadSucceededAction = {
  type: 'FETCH_ADMIN_HEAD_SUCCEEDED',
  data: {
    heads: {
      '1': MOCK_PAYLOAD,
    },
    meta: {
      node: {
        data: [],
      },
    },
  },
}

describe('reducers/adminHeadsShow', () => {
  it('returns the current state for other actions', () => {
    const action = {} as FetchAdminHeadSucceededAction
    const state = reducer(INITIAL_STATE, action)

    expect(state.adminHeadsShow).toEqual(INITIAL_STATE.adminHeadsShow)
  })

  describe('FETCH_ADMIN_HEAD_SUCCEEDED', () => {
    it('can receive new operator', () => {
      const state = reducer(INITIAL_STATE, MOCK_ACTION)
      expect(state.adminHeadsShow).toEqual({ 1: MOCK_PAYLOAD })
    })
  })
})
