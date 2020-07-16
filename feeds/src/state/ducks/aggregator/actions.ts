import { FeedConfig, OracleNode } from 'config'

/**
 * aggregator/CLEAR_STATE
 */
export interface ClearStateAction {
  type: 'aggregator/CLEAR_STATE'
}

export function clearState() {
  return {
    type: 'aggregator/CLEAR_STATE',
  }
}

/**
 * aggregator/FETCH_FEED_BY_PAIR_BEGIN
 */
export interface FetchFeedByPairBeginAction {
  type: 'aggregator/FETCH_FEED_BY_PAIR_BEGIN'
}

export function fetchFeedByPairBegin() {
  return {
    type: 'aggregator/FETCH_FEED_BY_PAIR_BEGIN',
  }
}

/**
 * aggregator/FETCH_FEED_BY_PAIR_SUCCESS
 */
export interface FetchFeedByPairSuccessAction {
  type: 'aggregator/FETCH_FEED_BY_PAIR_SUCCESS'
  payload: FeedConfig
}

export function fetchFeedByPairSuccess(payload: FeedConfig) {
  return {
    type: 'aggregator/FETCH_FEED_BY_PAIR_SUCCESS',
    payload,
  }
}

/**
 * aggregator/FETCH_FEED_BY_PAIR_ERROR
 */
export interface FetchFeedByPairErrorAction {
  type: 'aggregator/FETCH_FEED_BY_PAIR_ERROR'
  error: string
}

export function fetchFeedByPairError(error: string) {
  return {
    type: 'aggregator/FETCH_FEED_BY_PAIR_ERROR',
    error,
  }
}

/**
 * aggregator/FETCH_FEED_BY_ADDRESS_BEGIN
 */
export interface FetchFeedByAddressBeginAction {
  type: 'aggregator/FETCH_FEED_BY_ADDRESS_BEGIN'
}

export function fetchFeedByAddressBegin() {
  return {
    type: 'aggregator/FETCH_FEED_BY_ADDRESS_BEGIN',
  }
}

/**
 * aggregator/FETCH_FEED_BY_ADDRESS_SUCCESS
 */
export interface FetchFeedByAddressSuccessAction {
  type: 'aggregator/FETCH_FEED_BY_ADDRESS_SUCCESS'
  payload: FeedConfig
}

export function fetchFeedByAddressSuccess(payload: FeedConfig) {
  return {
    type: 'aggregator/FETCH_FEED_BY_ADDRESS_SUCCESS',
    payload,
  }
}

/**
 * aggregator/FETCH_FEED_BY_ADDRESS_ERROR
 */
export interface FetchFeedByAddressErrorAction {
  type: 'aggregator/FETCH_FEED_BY_ADDRESS_ERROR'
  error: string
}

export function fetchFeedByAddressError(error: string) {
  return {
    type: 'aggregator/FETCH_FEED_BY_ADDRESS_ERROR',
    error,
  }
}

/**
 * aggregator/FETCH_ORACLE_NODES_BEGIN
 */
export interface FetchOracleNodesBeginAction {
  type: 'aggregator/FETCH_ORACLE_NODES_BEGIN'
}

export function fetchOracleNodesBegin() {
  return {
    type: 'aggregator/FETCH_ORACLE_NODES_BEGIN',
  }
}

/**
 * aggregator/FETCH_ORACLE_NODES_SUCCESS
 */
export interface FetchOracleNodesSuccessAction {
  type: 'aggregator/FETCH_ORACLE_NODES_SUCCESS'
  payload: OracleNode[]
}

export function fetchOracleNodesSuccess(payload: OracleNode[]) {
  return {
    type: 'aggregator/FETCH_ORACLE_NODES_SUCCESS',
    payload,
  }
}

/**
 * aggregator/FETCH_ORACLE_NODES_ERROR
 */
export interface FetchOracleNodesErrorAction {
  type: 'aggregator/FETCH_ORACLE_NODES_ERROR'
  error: string
}

export function fetchOracleNodesError(error: string) {
  return {
    type: 'aggregator/FETCH_ORACLE_NODES_ERROR',
    error,
  }
}

/**
 * aggregator/ORACLE_LIST
 */
export interface SetOracleListAction {
  type: 'aggregator/ORACLE_LIST'
  payload: any
}

export function setOracleList(payload: Array<string>) {
  return {
    type: 'aggregator/ORACLE_LIST',
    payload,
  }
}

/**
 * aggregator/LATEST_ANSWER
 */
export interface SetLatestAnswerAction {
  type: 'aggregator/LATEST_ANSWER'
  payload: any
}

export function setLatestAnswer(payload: any) {
  return {
    type: 'aggregator/LATEST_ANSWER',
    payload,
  }
}

/**
 * aggregator/LATEST_COMPLETED_ANSWER_ID
 */
export interface SetLatestCompletedAnswerIdAction {
  type: 'aggregator/LATEST_COMPLETED_ANSWER_ID'
  payload: any
}

export function setLatestCompletedAnswerId(payload: any) {
  return {
    type: 'aggregator/LATEST_COMPLETED_ANSWER_ID',
    payload,
  }
}

/**
 * aggregator/PENDING_ANSWER_ID
 */
export interface SetPendingAnswerIdAction {
  type: 'aggregator/PENDING_ANSWER_ID'
  payload: any
}

export function setPendingAnswerId(payload: any) {
  return {
    type: 'aggregator/PENDING_ANSWER_ID',
    payload,
  }
}

/**
 * aggregator/NEXT_ANSWER_ID
 */
export interface SetNextAnswerIdAction {
  type: 'aggregator/NEXT_ANSWER_ID'
  payload: any
}

export function setNextAnswerId(payload: any) {
  return {
    type: 'aggregator/NEXT_ANSWER_ID',
    payload,
  }
}

/**
 * aggregator/ORACLE_ANSWERS
 */
export interface SetOracleAnswersAction {
  type: 'aggregator/ORACLE_ANSWERS'
  payload: any
}

export function setOracleAnswers(payload: any) {
  return {
    type: 'aggregator/ORACLE_ANSWERS',
    payload,
  }
}

/**
 * aggregator/LATEST_REQUEST_TIMESTAMP
 */
export interface SetLatestRequestTimestampAction {
  type: 'aggregator/LATEST_REQUEST_TIMESTAMP'
  payload: any
}

export function setLatestRequestTimestamp(payload: any) {
  return {
    type: 'aggregator/LATEST_REQUEST_TIMESTAMP',
    payload,
  }
}

/**
 * aggregator/MINIMUM_ANSWERS
 */
export interface SetMinimumAnswersAction {
  type: 'aggregator/MINIMUM_ANSWERS'
  payload: any
}

export function setMinumumAnswers(payload: any) {
  return {
    type: 'aggregator/MINIMUM_ANSWERS',
    payload,
  }
}

/**
 * aggregator/LATEST_ANSWER_TIMESTAMP
 */
export interface SetLatestAnswerTimestampAction {
  type: 'aggregator/LATEST_ANSWER_TIMESTAMP'
  payload: any
}

export function setLatestAnswerTimestamp(payload: any) {
  return {
    type: 'aggregator/LATEST_ANSWER_TIMESTAMP',
    payload,
  }
}

/**
 * aggregator/ANSWER_HISTORY
 */
export interface SetAnswersHistoryAction {
  type: 'aggregator/ANSWER_HISTORY'
  payload: any
}

export function setAnswerHistory(payload: any) {
  return {
    type: 'aggregator/ANSWER_HISTORY',
    payload,
  }
}

/**
 * aggregator/CONTRACT_ADDRESS
 */
export interface SetCurrentAddressAction {
  type: 'aggregator/CONTRACT_ADDRESS'
  payload: any
}

export function setContractAddress(payload: any) {
  return {
    type: 'aggregator/CONTRACT_ADDRESS',
    payload,
  }
}

/**
 * aggregator/ETHGAS_PRICE
 */
export interface SetEthGasPriceAction {
  type: 'aggregator/ETHGAS_PRICE'
  payload: any
}

export function setEthGasPrice(payload: any) {
  return {
    type: 'aggregator/ETHGAS_PRICE',
    payload,
  }
}
