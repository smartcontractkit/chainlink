import { ListingAnswer } from './operations'
import { Actions } from 'state/actions'

interface State {
  answers?: ListingAnswer[]
}

// TODO: Shouldn't need to export this with the explorer client testing method
export const INITIAL_STATE: State = {
  answers: undefined,
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'listing/SET_ANSWERS':
      return {
        ...state,
        answers: action.payload,
      }
    default:
      return state
  }
}

export default reducer
