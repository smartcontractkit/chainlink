type Attribute = string | number | null

export interface IState {
  data: { [key: string]: Attribute }
}

const initialState: IState = {
  data: {}
}

interface NormalizedResponse {
  configWhitelists: {
    [id: string]: {
      attributes: { [key: string]: Attribute }
    }
  }
}

export type Action = {
  type: 'UPSERT_CONFIGURATION'
  data: NormalizedResponse
}

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_CONFIGURATION':
      const id = Object.keys(action.data.configWhitelists)[0]
      const attributes = action.data.configWhitelists[id].attributes

      return {
        ...state,
        data: attributes
      }
    default:
      return state
  }
}
