export interface State {
  items?: ChainlinkNode[]
}

export interface NormalizedData {
  chainlinkNodes: any[]
}

export type JobRunsAction =
  | { type: 'UPSERT_JOB_RUNS'; data: NormalizedData }
  | { type: 'UPSERT_JOB_RUN'; data: NormalizedData }
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }

const INITIAL_STATE = { items: undefined }

export default (state: State = INITIAL_STATE, action: JobRunsAction) => {
  switch (action.type) {
    case '@@redux/INIT':
    case '@@INIT':
      return INITIAL_STATE
    case 'UPSERT_JOB_RUNS':
      return { items: action.data.chainlinkNodes }
    case 'UPSERT_JOB_RUN':
      return { items: action.data.chainlinkNodes }
    default:
      return state
  }
}
