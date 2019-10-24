import { Actions } from './actions'

export interface State {
  items?: ChainlinkNode[]
}

const INITIAL_STATE: State = { items: undefined }

export default (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'FETCH_JOB_RUNS_SUCCEEDED':
      return { items: { ...action.data.chainlinkNodes } }
    case 'FETCH_JOB_RUN_SUCCEEDED':
      return { items: { ...action.data.chainlinkNodes } }
    default:
      return state
  }
}
