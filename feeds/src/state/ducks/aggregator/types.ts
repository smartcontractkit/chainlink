import { FeedConfig, OracleNode } from 'config'

export const CLEAR_STATE = 'aggregator/CLEAR_STATE'
export const FETCH_FEED_BY_PAIR_BEGIN = 'aggregator/FETCH_FEED_BY_PAIR_BEGIN'
export const FETCH_FEED_BY_PAIR_SUCCESS =
  'aggregator/FETCH_FEED_BY_PAIR_SUCCESS'
export const FETCH_FEED_BY_PAIR_ERROR = 'aggregator/FETCH_FEED_BY_PAIR_ERROR'
export const FETCH_FEED_BY_ADDRESS_BEGIN =
  'aggregator/FETCH_FEED_BY_ADDRESS_BEGIN'
export const FETCH_FEED_BY_ADDRESS_SUCCESS =
  'aggregator/FETCH_FEED_BY_ADDRESS_SUCCESS'
export const FETCH_FEED_BY_ADDRESS_ERROR =
  'aggregator/FETCH_FEED_BY_ADDRESS_ERROR'
export const FETCH_ORACLE_NODES_BEGIN = 'aggregator/FETCH_ORACLE_NODES_BEGIN'
export const FETCH_ORACLE_NODES_SUCCESS =
  'aggregator/FETCH_ORACLE_NODES_SUCCESS'
export const FETCH_ORACLE_NODES_ERROR = 'aggregator/FETCH_ORACLE_NODES_ERROR'
export const ORACLE_LIST = 'aggregator/ORACLE_LIST'
export const LATEST_ANSWER = 'aggregator/LATEST_ANSWER'
export const LATEST_COMPLETED_ANSWER_ID =
  'aggregator/LATEST_COMPLETED_ANSWER_ID'
export const PENDING_ANSWER_ID = 'aggregator/PENDING_ANSWER_ID'
export const NEXT_ANSWER_ID = 'aggregator/NEXT_ANSWER_ID'
export const ORACLE_ANSWERS = 'aggregator/ORACLE_ANSWERS'
export const LATEST_REQUEST_TIMESTAMP = 'aggregator/LATEST_REQUEST_TIMESTAMP'
export const MINIMUM_ANSWERS = 'aggregator/MINIMUM_ANSWERS'
export const LATEST_ANSWER_TIMESTAMP = 'aggregator/LATEST_ANSWER_TIMESTAMP'
export const ANSWER_HISTORY = 'aggregator/ANSWER_HISTORY'
export const CONTRACT_ADDRESS = 'aggregator/CONTRACT_ADDRESS'
export const ETHGAS_PRICE = 'aggregator/ETHGAS_PRICE'

export interface ClearStateAction {
  type: typeof CLEAR_STATE
}

export interface FetchFeedByPairBeginAction {
  type: typeof FETCH_FEED_BY_PAIR_BEGIN
}

export interface FetchFeedByPairSuccessAction {
  type: typeof FETCH_FEED_BY_PAIR_SUCCESS
  payload: FeedConfig
}

export interface FetchFeedByPairErrorAction {
  type: typeof FETCH_FEED_BY_PAIR_ERROR
  error: string
}

export interface FetchFeedByAddressBeginAction {
  type: typeof FETCH_FEED_BY_ADDRESS_BEGIN
}

export interface FetchFeedByAddressSuccessAction {
  type: typeof FETCH_FEED_BY_ADDRESS_SUCCESS
  payload: FeedConfig
}

export interface FetchFeedByAddressErrorAction {
  type: typeof FETCH_FEED_BY_ADDRESS_ERROR
  error: string
}

export interface FetchOracleNodesBeginAction {
  type: typeof FETCH_ORACLE_NODES_BEGIN
}

export interface FetchOracleNodesSuccessAction {
  type: typeof FETCH_ORACLE_NODES_SUCCESS
  payload: OracleNode[]
}

export interface FetchOracleNodesErrorAction {
  type: typeof FETCH_ORACLE_NODES_ERROR
  error: string
}

export interface SetOracleListAction {
  type: typeof ORACLE_LIST
  payload: any
}

export interface SetLatestAnswerAction {
  type: typeof LATEST_ANSWER
  payload: any
}

export interface SetLatestCompletedAnswerIdAction {
  type: typeof LATEST_COMPLETED_ANSWER_ID
  payload: any
}

export interface SetPendingAnswerIdAction {
  type: typeof PENDING_ANSWER_ID
  payload: any
}

export interface SetNextAnswerIdAction {
  type: typeof NEXT_ANSWER_ID
  payload: any
}

export interface SetOracleAnswersAction {
  type: typeof ORACLE_ANSWERS
  payload: any
}

export interface SetLatestRequestTimestampAction {
  type: typeof LATEST_REQUEST_TIMESTAMP
  payload: any
}

export interface SetMinimumAnswersAction {
  type: typeof MINIMUM_ANSWERS
  payload: any
}

export interface SetLatestAnswerTimestampAction {
  type: typeof LATEST_ANSWER_TIMESTAMP
  payload: any
}

export interface SetAnswersHistoryAction {
  type: typeof ANSWER_HISTORY
  payload: any
}

export interface SetCurrentAddressAction {
  type: typeof CONTRACT_ADDRESS
  payload: any
}

export interface SetEthGasPriceAction {
  type: typeof ETHGAS_PRICE
  payload: any
}

export type AggregatorActionTypes =
  | ClearStateAction
  | FetchFeedByPairBeginAction
  | FetchFeedByPairSuccessAction
  | FetchFeedByPairErrorAction
  | FetchFeedByAddressBeginAction
  | FetchFeedByAddressSuccessAction
  | FetchFeedByAddressErrorAction
  | FetchOracleNodesBeginAction
  | FetchOracleNodesSuccessAction
  | FetchOracleNodesErrorAction
  | SetOracleListAction
  | SetLatestAnswerAction
  | SetLatestCompletedAnswerIdAction
  | SetPendingAnswerIdAction
  | SetNextAnswerIdAction
  | SetOracleAnswersAction
  | SetLatestRequestTimestampAction
  | SetMinimumAnswersAction
  | SetLatestAnswerTimestampAction
  | SetAnswersHistoryAction
  | SetCurrentAddressAction
  | SetEthGasPriceAction
