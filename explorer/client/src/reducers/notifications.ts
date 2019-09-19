export interface State {
  errors: string[]
}

export type Query = string | undefined

export type Action =
  | { type: '@@redux/INIT' }
  | { type: '@@INIT' }
  | { type: 'NOTIFY_ERROR'; text: string }

const INITIAL_STATE: State = { errors: [] }

export default (state: State = INITIAL_STATE, action: Action) => {
  switch (action.type) {
    case '@@redux/INIT':
    case '@@INIT':
      return INITIAL_STATE
    case 'NOTIFY_ERROR':
      return { errors: [action.text] }
    default:
      return state
  }
}
