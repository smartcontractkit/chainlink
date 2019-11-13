export interface State {
  errors: string[]
}

export type Action = { type: 'NOTIFY_ERROR'; text: string }

const INITIAL_STATE: State = { errors: [] }

export default (state: State = INITIAL_STATE, action: Action) => {
  switch (action.type) {
    case 'NOTIFY_ERROR':
      return { errors: [action.text] }
    default:
      return state
  }
}
