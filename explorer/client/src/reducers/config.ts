import { Actions } from './actions'
import { Reducer } from 'redux'

export interface State {
  etherscanHost?: string
}

const initialState: State = { etherscanHost: undefined }

const configReducer: Reducer<State, Actions> = (
  state = initialState,
  action,
) => {
  switch (action.type) {
    case 'FETCH_JOB_RUN_SUCCEEDED':
      return {
        ...state,
        etherscanHost: action.data.meta.jobRun.meta.etherscanHost,
      }
    default:
      return state
  }
}

export default configReducer
