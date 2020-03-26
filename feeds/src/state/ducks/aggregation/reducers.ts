import { Actions } from 'state/actions'

export interface State {
  options: null | any
  contractAddress: null | any
  oracles: null | any
  currentAnswer: null | any
  latestCompletedAnswerId: null | any
  pendingAnswerId: null | any
  nextAnswerId: null | any
  oracleResponse: null | any
  requestTime: null | any
  minimumResponses: null | any
  updateHeight: null | any
  answerHistory: null | any
  ethGasPrice: null | any
}

export const INITIAL_STATE: State = {
  options: null,
  contractAddress: null,
  oracles: null,
  currentAnswer: null,
  latestCompletedAnswerId: null,
  pendingAnswerId: null,
  nextAnswerId: null,
  oracleResponse: null,
  requestTime: null,
  minimumResponses: null,
  updateHeight: null,
  answerHistory: null,
  ethGasPrice: null,
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'aggregation/CLEAR_STATE':
      return INITIAL_STATE

    case 'aggregation/OPTIONS':
      return {
        ...state,
        options: action.payload,
      }

    case 'aggregation/CONTRACT_ADDRESS':
      return {
        ...state,
        contractAddress: action.payload,
      }

    case 'aggregation/ORACLES':
      return {
        ...state,
        oracles: action.payload,
      }

    case 'aggregation/CURRENT_ANSWER':
      return {
        ...state,
        currentAnswer: action.payload,
      }

    case 'aggregation/LATEST_COMPLETED_ANSWER_ID':
      return {
        ...state,
        latestCompletedAnswerId: action.payload,
      }

    case 'aggregation/PENDING_ANSWER_ID':
      return {
        ...state,
        pendingAnswerId: action.payload,
      }

    case 'aggregation/NEXT_ANSWER_ID':
      return {
        ...state,
        nextAnswerId: action.payload,
      }

    case 'aggregation/ORACLE_RESPONSE':
      return {
        ...state,
        oracleResponse: action.payload,
      }

    case 'aggregation/REQUEST_TIME':
      return {
        ...state,
        requestTime: action.payload,
      }

    case 'aggregation/MINIMUM_RESPONSES':
      return {
        ...state,
        minimumResponses: action.payload,
      }

    case 'aggregation/UPDATE_HEIGHT':
      return {
        ...state,
        updateHeight: action.payload,
      }

    case 'aggregation/ANSWER_HISTORY':
      return {
        ...state,
        answerHistory: action.payload,
      }

    case 'aggregation/ETHGAS_PRICE':
      return {
        ...state,
        ethGasPrice: action.payload,
      }

    default:
      return state
  }
}

export default reducer
