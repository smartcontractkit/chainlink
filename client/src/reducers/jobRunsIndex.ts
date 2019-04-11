import { JobRunsAction } from './jobRuns'

export interface IState {
  items?: string[]
}

const INITIAL_STATE = { items: undefined }

export default (
  state: IState = INITIAL_STATE,
  action: JobRunsAction
): IState => {
  switch (action.type) {
    case '@@redux/INIT':
    case '@@INIT':
      return INITIAL_STATE
    case 'UPSERT_JOB_RUNS':
      return { items: action.data.result }
    case 'UPSERT_JOB_RUN':
      return INITIAL_STATE
    default:
      return state
  }
}
