export interface State {
  items?: ChainlinkNode[]
}

export interface NormalizedData {
  chainlinkNodes: any[]
}

export type JobRunsAction =
  | { type: 'UPSERT_JOB_RUNS'; data: NormalizedData }
  | { type: 'UPSERT_JOB_RUN'; data: NormalizedData }

const INITIAL_STATE: State = { items: undefined }

export default (state: State = INITIAL_STATE, action: JobRunsAction) => {
  switch (action.type) {
    case 'UPSERT_JOB_RUNS':
      return { items: action.data.chainlinkNodes }
    case 'UPSERT_JOB_RUN':
      return { items: action.data.chainlinkNodes }
    default:
      return state
  }
}
