import { FeedConfig, OracleNode } from 'config'
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
  STORE_AGGREGATOR_CONFIG,
  AggregatorActionTypes,
} from './types'

export function clearState(): AggregatorActionTypes {
  return {
    type: CLEAR_STATE,
  }
}

export function fetchFeedByPairBegin(): AggregatorActionTypes {
  return {
    type: FETCH_FEED_BY_PAIR_BEGIN,
  }
}

export function fetchFeedByPairSuccess(
  payload: FeedConfig,
): AggregatorActionTypes {
  return {
    type: FETCH_FEED_BY_PAIR_SUCCESS,
    payload,
  }
}

export function fetchFeedByPairError(error: string): AggregatorActionTypes {
  return {
    type: FETCH_FEED_BY_PAIR_ERROR,
    error,
  }
}

export function fetchFeedByAddressBegin(): AggregatorActionTypes {
  return {
    type: FETCH_FEED_BY_ADDRESS_BEGIN,
  }
}

export function fetchFeedByAddressSuccess(
  payload: FeedConfig,
): AggregatorActionTypes {
  return {
    type: FETCH_FEED_BY_ADDRESS_SUCCESS,
    payload,
  }
}

export function fetchFeedByAddressError(error: string): AggregatorActionTypes {
  return {
    type: FETCH_FEED_BY_ADDRESS_ERROR,
    error,
  }
}

export function fetchOracleNodesBegin(): AggregatorActionTypes {
  return {
    type: FETCH_ORACLE_NODES_BEGIN,
  }
}

export function fetchOracleNodesSuccess(
  payload: OracleNode[],
): AggregatorActionTypes {
  return {
    type: FETCH_ORACLE_NODES_SUCCESS,
    payload,
  }
}

export function fetchOracleNodesError(error: string): AggregatorActionTypes {
  return {
    type: FETCH_ORACLE_NODES_ERROR,
    error,
  }
}

export function setOracleList(payload: Array<string>): AggregatorActionTypes {
  return {
    type: ORACLE_LIST,
    payload,
  }
}

export function setLatestAnswer(payload: any): AggregatorActionTypes {
  return {
    type: LATEST_ANSWER,
    payload,
  }
}

export function setLatestCompletedAnswerId(
  payload: any,
): AggregatorActionTypes {
  return {
    type: LATEST_COMPLETED_ANSWER_ID,
    payload,
  }
}

export function setPendingAnswerId(payload: any): AggregatorActionTypes {
  return {
    type: PENDING_ANSWER_ID,
    payload,
  }
}

export function setNextAnswerId(payload: any): AggregatorActionTypes {
  return {
    type: NEXT_ANSWER_ID,
    payload,
  }
}

export function setOracleAnswers(payload: any): AggregatorActionTypes {
  return {
    type: ORACLE_ANSWERS,
    payload,
  }
}

export function setLatestRequestTimestamp(payload: any): AggregatorActionTypes {
  return {
    type: LATEST_REQUEST_TIMESTAMP,
    payload,
  }
}

export function setMinumumAnswers(payload: any): AggregatorActionTypes {
  return {
    type: MINIMUM_ANSWERS,
    payload,
  }
}

export function setLatestAnswerTimestamp(payload: any): AggregatorActionTypes {
  return {
    type: LATEST_ANSWER_TIMESTAMP,
    payload,
  }
}

export function setAnswerHistory(payload: any): AggregatorActionTypes {
  return {
    type: ANSWER_HISTORY,
    payload,
  }
}

export function setContractAddress(payload: any): AggregatorActionTypes {
  return {
    type: CONTRACT_ADDRESS,
    payload,
  }
}

export function setEthGasPrice(payload: any): AggregatorActionTypes {
  return {
    type: ETHGAS_PRICE,
    payload,
  }
}

export function storeAggregatorConfig(config: FeedConfig): AggregatorActionTypes {
  return {
    type: STORE_AGGREGATOR_CONFIG,
    payload: { config },
  }
}
