export interface State {
  items?: JobRun[]
}

export interface Meta {
  count: number
}

export interface NormalizedEndpoint {
  data: any[]
  meta: Meta
}

export interface NormalizedMeta {
  jobRuns: NormalizedEndpoint
}

export interface NormalizedData {
  jobRuns: any
  meta: NormalizedMeta
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
      return { items: action.data.jobRuns }
    case 'UPSERT_JOB_RUN':
      return { items: action.data.jobRuns }
    default:
      return state
  }
}
