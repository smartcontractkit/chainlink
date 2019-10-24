import { Actions } from './actions'

export interface State {
  etherscanHost?: string
}

const initialState: State = { etherscanHost: undefined }

export default (state: State = initialState, action: Actions) => {
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
