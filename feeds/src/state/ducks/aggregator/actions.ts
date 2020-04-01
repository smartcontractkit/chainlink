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
 * aggregator/CONFIG
 */
export interface SetConfigAction {
  type: 'aggregator/CONFIG'
  payload: any
}

export function setConfig(payload: any) {
  return {
    type: 'aggregator/CONFIG',
    payload,
  }
}

/**
 * aggregator/CLEAR_STATE
 */
export interface SetClearStateAction {
  type: 'aggregator/CLEAR_STATE'
  payload: any
}

export function clearState() {
  return {
    type: 'aggregator/CLEAR_STATE',
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
