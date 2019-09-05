export interface IState {
  data: Record<string, Attribute>
}

export interface Action {
  type: ConfigurationActionType.UPSERT
  data: NormalizedResponse
}

enum ConfigurationActionType {
  UPSERT = 'UPSERT_CONFIGURATION'
}

interface NormalizedResponse {
  configWhitelists: Record<string, Attributes>
}
interface Attributes {
  attributes: Record<string, Attribute>
}
type Attribute = string | number | null

const initialState: IState = {
  data: {}
}

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case ConfigurationActionType.UPSERT:
      const id = Object.keys(action.data.configWhitelists)[0]
      const attributes = action.data.configWhitelists[id].attributes

      return { ...state, data: attributes }

    default:
      return state
  }
}
