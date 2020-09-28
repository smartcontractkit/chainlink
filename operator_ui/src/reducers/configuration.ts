import { Reducer } from 'redux'
import { Actions, ConfigurationAttribute, ResourceActionType } from './actions'

export interface State {
  data: Record<string, ConfigurationAttribute>
}

const INITIAL_STATE: State = {
  data: {},
}

const reducer: Reducer<State, Actions> = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case ResourceActionType.UPSERT_CONFIGURATION: {
      const id = Object.keys(action.data.configPrinters)[0]
      const attributes = action.data.configPrinters[id].attributes

      return { ...state, data: attributes }
    }

    default:
      return state
  }
}

export default reducer
