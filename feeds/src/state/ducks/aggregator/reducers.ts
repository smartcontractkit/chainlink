import 'core-js/stable/object/from-entries'
import { Reducer } from 'redux'
import { OracleNode } from 'config'
import {
  CLEAR_STATE,
  FETCH_FEED_BY_PAIR_BEGIN,
  FETCH_FEED_BY_PAIR_SUCCESS,
  FETCH_FEED_BY_PAIR_ERROR,
  FETCH_FEED_BY_ADDRESS_BEGIN,
  FETCH_FEED_BY_ADDRESS_SUCCESS,
  FETCH_FEED_BY_ADDRESS_ERROR,
  FETCH_ORACLE_NODES_BEGIN,
  FETCH_ORACLE_NODES_SUCCESS,
  FETCH_ORACLE_NODES_ERROR,
  ORACLE_LIST,
  LATEST_ANSWER,
  LATEST_COMPLETED_ANSWER_ID,
  PENDING_ANSWER_ID,
  NEXT_ANSWER_ID,
  ORACLE_ANSWERS,
  LATEST_REQUEST_TIMESTAMP,
  MINIMUM_ANSWERS,
  LATEST_ANSWER_TIMESTAMP,
  ANSWER_HISTORY,
  CONTRACT_ADDRESS,
  ETHGAS_PRICE,
  AggregatorActionTypes,
} from './types'

export interface State {
  loadingFeed: boolean
  errorFeed?: string
  config: null | any
  loadingOracleNodes: boolean
  errorOracleNodes?: string
  oracleNodes: Record<OracleNode['address'], OracleNode>
  contractAddress: null | any
  oracleList: Array<OracleNode['address']> | any
  latestAnswer: null | any
  latestCompletedAnswerId: null | any
  pendingAnswerId: null | any
  nextAnswerId: null | any
  oracleAnswers: any[]
  latestRequestTimestamp: null | any
  minimumAnswers: null | any
  latestAnswerTimestamp: null | any
  answerHistory: null | any
  ethGasPrice: null | any
}

export const INITIAL_STATE: State = {
  loadingFeed: false,
  config: null,
  loadingOracleNodes: false,
  oracleNodes: {},
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

const reducer: Reducer<State, AggregatorActionTypes> = (
  state = INITIAL_STATE,
  action,
) => {
  switch (action.type) {
    case CLEAR_STATE: {
      return {
        ...INITIAL_STATE,
      }
    }

    case FETCH_FEED_BY_PAIR_BEGIN: {
      return {
        ...INITIAL_STATE,
        loadingFeed: true,
        errorFeed: undefined,
      }
    }

    case FETCH_FEED_BY_PAIR_SUCCESS: {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: undefined,
        config: action.payload,
      }
    }

    case FETCH_FEED_BY_PAIR_ERROR: {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: action.error,
      }
    }

    case FETCH_FEED_BY_ADDRESS_BEGIN: {
      return {
        ...INITIAL_STATE,
        loadingFeed: true,
        errorFeed: undefined,
      }
    }

    case FETCH_FEED_BY_ADDRESS_SUCCESS: {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: undefined,
        config: action.payload,
      }
    }

    case FETCH_FEED_BY_ADDRESS_ERROR: {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: action.error,
      }
    }

    case FETCH_ORACLE_NODES_BEGIN: {
      return {
        ...state,
        loadingOracleNodes: true,
        errorOracleNodes: undefined,
      }
    }

    case FETCH_ORACLE_NODES_SUCCESS: {
      const oracleNodes = Object.fromEntries(
        // action.payload.map(n => [n.address, n]),
        action.payload.map(n => [n.oracleAddress, n]),
      )

      return {
        ...state,
        loadingOracleNodes: false,
        errorOracleNodes: undefined,
        oracleNodes,
      }
    }

    case FETCH_ORACLE_NODES_ERROR: {
      return {
        ...state,
        loadingOracleNodes: false,
        errorOracleNodes: action.error,
      }
    }

    case CONTRACT_ADDRESS:
      return {
        ...state,
        contractAddress: action.payload,
      }

    case ORACLE_LIST:
      return {
        ...state,
        oracleList: action.payload,
      }

    case LATEST_ANSWER:
      return {
        ...state,
        latestAnswer: action.payload,
      }

    case LATEST_COMPLETED_ANSWER_ID:
      return {
        ...state,
        latestCompletedAnswerId: action.payload,
      }

    case PENDING_ANSWER_ID:
      return {
        ...state,
        pendingAnswerId: action.payload,
      }

    case NEXT_ANSWER_ID:
      return {
        ...state,
        nextAnswerId: action.payload,
      }

    case ORACLE_ANSWERS:
      return {
        ...state,
        oracleAnswers: action.payload,
      }

    case LATEST_REQUEST_TIMESTAMP:
      return {
        ...state,
        latestRequestTimestamp: action.payload,
      }

    case MINIMUM_ANSWERS:
      return {
        ...state,
        minimumAnswers: action.payload,
      }

    case LATEST_ANSWER_TIMESTAMP:
      return {
        ...state,
        latestAnswerTimestamp: action.payload,
      }

    case ANSWER_HISTORY:
      return {
        ...state,
        answerHistory: action.payload,
      }

    case ETHGAS_PRICE:
      return {
        ...state,
        ethGasPrice: action.payload,
      }

    default:
      return state
  }
}

export default reducer
