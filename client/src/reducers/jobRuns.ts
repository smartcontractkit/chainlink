export interface IState {
  items?: IJobRun[]
}

export type JobRunsAction =
  { type: 'UPSERT_JOB_RUNS', items: IJobRun[] } 
  | { type: '@@INIT' }

const initialState = { items: undefined }

export default (state: IState = initialState, action: JobRunsAction) => {
  switch (action.type) {
    case 'UPSERT_JOB_RUNS':
      return Object.assign(
        {},
        state,
        { items: action.items }
      )
    default:
      return state
  }
}
