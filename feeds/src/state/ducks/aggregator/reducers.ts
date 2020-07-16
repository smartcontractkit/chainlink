import 'core-js/stable/object/from-entries'
import { Actions } from 'state/actions'
import { OracleNode } from 'config'

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

const reducer = (state: State = INITIAL_STATE, action: Actions) => {
  switch (action.type) {
    case 'aggregator/CLEAR_STATE': {
      return {
        ...INITIAL_STATE,
      }
    }

    case 'aggregator/FETCH_FEED_BY_PAIR_BEGIN': {
      return {
        ...INITIAL_STATE,
        loadingFeed: true,
        errorFeed: undefined,
      }
    }

    case 'aggregator/FETCH_FEED_BY_PAIR_SUCCESS': {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: undefined,
        config: action.payload,
      }
    }

    case 'aggregator/FETCH_FEED_BY_PAIR_ERROR': {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: action.error,
      }
    }

    case 'aggregator/FETCH_FEED_BY_ADDRESS_BEGIN': {
      return {
        ...INITIAL_STATE,
        loadingFeed: true,
        errorFeed: undefined,
      }
    }

    case 'aggregator/FETCH_FEED_BY_ADDRESS_SUCCESS': {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: undefined,
        config: action.payload,
      }
    }

    case 'aggregator/FETCH_FEED_BY_ADDRESS_ERROR': {
      return {
        ...state,
        loadingFeed: false,
        errorFeed: action.error,
      }
    }

    case 'aggregator/FETCH_ORACLE_NODES_BEGIN': {
      return {
        ...state,
        loadingOracleNodes: true,
        errorOracleNodes: undefined,
      }
    }

    case 'aggregator/FETCH_ORACLE_NODES_SUCCESS': {
      const oracleNodes = Object.fromEntries(
        action.payload.map(n => [n.address, n]),
      )

      return {
        ...state,
        loadingOracleNodes: false,
        errorOracleNodes: undefined,
        oracleNodes,
      }
    }

    case 'aggregator/FETCH_ORACLE_NODES_ERROR': {
      return {
        ...state,
        loadingOracleNodes: false,
        errorOracleNodes: action.error,
      }
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
