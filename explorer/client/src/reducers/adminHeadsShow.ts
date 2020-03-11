import { Actions } from './actions'
import { Reducer } from 'redux'
import { BigNumber } from 'bignumber.js'

export interface HeadShowData {
  id: number
  parentHash: Buffer
  uncleHash: Buffer
  coinbase: Buffer
  root: Buffer
  txHash: Buffer
  receiptHash: Buffer
  bloom: Buffer
  difficulty: BigNumber
  number: BigNumber
  gasLimit: BigNumber
  gasUsed: BigNumber
  time: BigNumber
  extra: Buffer
  mixDigest: Buffer
  nonce: Buffer
}

export interface State {
  id?: {
    attributes: HeadShowData
  }
}

const INITIAL_STATE: State = {}

export const adminHeadsShow: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_HEAD_SUCCEEDED': {
      console.log('!!!!!!!!! FETCH_ADMIN_HEAD_SUCCEEDED: %o', action)
      return action.data.heads
    }
    case 'FETCH_ADMIN_HEAD_ERROR': {
      console.log('!!!!!!!! FETCH ADMIN HEAD ERROR: %o', action)
      return state
    }
    default: {
      console.log('adminHeadsShow reducer (default): %o', action)
      return state
    }
  }
}

export default adminHeadsShow
