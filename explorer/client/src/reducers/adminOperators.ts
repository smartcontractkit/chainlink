export interface State {
  errors: string[]
}

export type Action = { type: 'FETCH_OPERATORS_SUCCEEDED'; data: object[] }

const INITIAL_STATE: State = { errors: [] }

export default (state: State = INITIAL_STATE, action: Action) => {
  switch (action.type) {
    case 'FETCH_OPERATORS_SUCCEEDED':
      return { errors: action.data }
    default:
      return state
  }
}
