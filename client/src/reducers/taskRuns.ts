export interface IState {
  items?: ITaskRun[]
}

export interface INormalizedData {
  entities: any
  result: any
}

export type TaskRunsAction =
  | { type: 'UPSERT_JOB_RUN'; data: INormalizedData }
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }

const INITIAL_STATE = { items: undefined }

export default (state: IState = INITIAL_STATE, action: TaskRunsAction) => {
  switch (action.type) {
    case 'UPSERT_JOB_RUN':
      return { items: action.data.entities.taskRuns }
    default:
      return state
  }
}
