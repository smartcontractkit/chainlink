export interface IState {
  items: object
  currentPage: string[]
  count: number
}

const initialState: IState = {
  items: {},
  currentPage: [],
  count: 0
}

interface NormalizedBridgeResponse {
  bridges: { [id: string]: object }
}

interface NormalizedBridgesResponse {
  bridges: { [id: string]: object }
  meta: {
    currentPageBridges: {
      data: { id: string }[]
      meta: { count: number }
    }
  }
}

export type Action =
  | { type: 'UPSERT_BRIDGES'; data: NormalizedBridgesResponse }
  | { type: 'UPSERT_BRIDGE'; data: NormalizedBridgeResponse }

export default (state: IState = initialState, action: Action) => {
  switch (action.type) {
    case 'UPSERT_BRIDGES': {
      const { bridges, meta } = action.data

      return Object.assign({}, state, {
        items: Object.assign({}, state.items, bridges),
        currentPage: meta.currentPageBridges.data.map(b => b.id),
        count: meta.currentPageBridges.meta.count
      })
    }
    case 'UPSERT_BRIDGE':
      return Object.assign({}, state, {
        items: Object.assign({}, state.items, action.data.bridges)
      })
    default:
      return state
  }
}
