import { Actions } from './actions'

export interface State {
  items?: JobRun[]
}

const INITIAL_STATE: State = {}

export default (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'FETCH_JOB_RUNS_SUCCEEDED':
      return { items: { ...action.data.jobRuns } }
    case 'FETCH_JOB_RUN_SUCCEEDED':
      return { items: { ...action.data.jobRuns } }
    default:
      return state
  }
}
