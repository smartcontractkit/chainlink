import { Actions } from 'state/actions'

export interface State {
  config: null | any
  contractAddress: null | any
  oracleList: Array<string> | any
  latestAnswer: null | any
  latestCompletedAnswerId: null | any
  pendingAnswerId: null | any
  nextAnswerId: null | any
  oracleAnswers: null | any
  latestRequestTimestamp: null | any
  minimumAnswers: null | any
  latestAnswerTimestamp: null | any
  answerHistory: null | any
  ethGasPrice: null | any
}

export const INITIAL_STATE: State = {
  config: null,
  contractAddress: null,
  oracleList: null,
  latestAnswer: null,
  latestCompletedAnswerId: null,
  pendingAnswerId: null,
  nextAnswerId: null,
  oracleAnswers: [],
  latestRequestTimestamp: null,
  minimumAnswers: null,
  latestAnswerTimestamp: null,
  answerHistory: null,
  ethGasPrice: null,
}

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'aggregator/CLEAR_STATE':
      return INITIAL_STATE

    case 'aggregator/CONFIG':
      return {
        ...state,
        config: action.payload,
      }

    case 'aggregator/CONTRACT_ADDRESS':
      return {
        ...state,
        contractAddress: action.payload,
      }

    case 'aggregator/ORACLE_LIST':
      return {
        ...state,
        oracleList: action.payload,
      }

    case 'aggregator/LATEST_ANSWER':
      return {
        ...state,
        latestAnswer: action.payload,
      }

    case 'aggregator/LATEST_COMPLETED_ANSWER_ID':
      return {
        ...state,
        latestCompletedAnswerId: action.payload,
      }

    case 'aggregator/PENDING_ANSWER_ID':
      return {
        ...state,
        pendingAnswerId: action.payload,
      }

    case 'aggregator/NEXT_ANSWER_ID':
      return {
        ...state,
        nextAnswerId: action.payload,
      }

    case 'aggregator/ORACLE_ANSWERS':
      return {
        ...state,
        oracleAnswers: action.payload,
      }

    case 'aggregator/LATEST_REQUEST_TIMESTAMP':
      return {
        ...state,
        latestRequestTimestamp: action.payload,
      }

    case 'aggregator/MINIMUM_ANSWERS':
      return {
        ...state,
        minimumAnswers: action.payload,
      }

    case 'aggregator/LATEST_ANSWER_TIMESTAMP':
      return {
        ...state,
        latestAnswerTimestamp: action.payload,
      }

    case 'aggregator/ANSWER_HISTORY':
      return {
        ...state,
        answerHistory: action.payload,
      }

    case 'aggregator/ETHGAS_PRICE':
      return {
        ...state,
        ethGasPrice: action.payload,
      }

    default:
      return state
  }
}

export default reducer
