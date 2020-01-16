import { Reducer } from 'redux'
import { Actions, ConfigurationAttribute } from './actions'

export interface State {
  data: Record<string, ConfigurationAttribute>
}

const INITIAL_STATE: State = {
  data: {},
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case 'UPSERT_CONFIGURATION': {
      const id = Object.keys(action.data.configWhitelists)[0]
      const attributes = action.data.configWhitelists[id].attributes

      return { ...state, data: attributes }
    }

    default:
      return state
  }
}

export default reducer
