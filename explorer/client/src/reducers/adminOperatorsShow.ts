import { Actions } from './actions'
import { Reducer } from 'redux'

export interface OperatorShowData {
  id: string
  name: string
  url?: string
  createdAt: string
  uptime: number
  coreVersion?: string
  coreSha?: string
  jobCounts: {
    completed: number
    errored: number
    inProgress: number
    total: number
  }
}

export interface State {
  id?: {
    attributes: OperatorShowData
  }
}

const INITIAL_STATE: State = {}

export const adminOperatorsShow: Reducer<State, Actions> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case 'FETCH_ADMIN_OPERATOR_SUCCEEDED': {
      return action.data.chainlinkNodes
    }
    default: {
      return state
    }
  }
}

export default adminOperatorsShow
