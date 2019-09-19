export interface State {
  errors: string[]
}

export type Query = string | undefined

export type Action =
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }
  | { type: 'FETCH_OPERATORS_SUCCEEDED'; data: object[] }

const INITIAL_STATE: State = { errors: [] }

export default (state: State = INITIAL_STATE, action: Action) => {
  switch (action.type) {
    case '@@redux/INIT':
    case '@@INIT':
      return INITIAL_STATE
    case 'FETCH_OPERATORS_SUCCEEDED':
      return { errors: action.data }
    default:
      return state
  }
}
