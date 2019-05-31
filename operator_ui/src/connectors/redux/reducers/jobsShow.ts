export interface IState {
  runsCount?: number
}

const initialState: IState = {
  runsCount: undefined
}

interface IData {
  specs: any[]
  runs?: any[]
}

export type Action =
  | { type: 'UPSERT_JOB'; data: IData }
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_JOB': {
      const keys = Object.keys(action.data.runs || {})
      return Object.assign({}, state, { runsCount: keys.length })
    }
    default:
      return state
  }
}
