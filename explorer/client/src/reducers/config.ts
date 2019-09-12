export interface State {
  etherscanHost?: string
}

export type Action =
  | { type: 'UPSERT_JOB_RUN'; data: any }
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }

const initialState = { etherscanHost: undefined }

export default (state: State = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_JOB_RUN':
      return Object.assign({}, state, {
        etherscanHost: action.data.meta.jobRun.meta.etherscanHost,
      })
    default:
      return state
  }
}
