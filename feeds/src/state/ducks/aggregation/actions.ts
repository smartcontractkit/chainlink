/**
 * aggregation/ORACLES
 */
export interface SetOraclesAction {
  type: 'aggregation/ORACLES'
  payload: any
}

export function setOracles(payload: any) {
  return {
    type: 'aggregation/ORACLES',
    payload,
  }
}

/**
 * aggregation/CURRENT_ANSWER
 */
export interface SetCurrentAnswerAction {
  type: 'aggregation/CURRENT_ANSWER'
  payload: any
}

export function setCurrentAnswer(payload: any) {
  return {
    type: 'aggregation/CURRENT_ANSWER',
    payload,
  }
}

/**
 * aggregation/LATEST_COMPLETED_ANSWER_ID
 */
export interface SetLatestCompletedAnswerIdAction {
  type: 'aggregation/LATEST_COMPLETED_ANSWER_ID'
  payload: any
}

export function setLatestCompletedAnswerId(payload: any) {
  return {
    type: 'aggregation/LATEST_COMPLETED_ANSWER_ID',
    payload,
  }
}

/**
 * aggregation/PENDING_ANSWER_ID
 */
export interface SetPendingAnswerIdAction {
  type: 'aggregation/PENDING_ANSWER_ID'
  payload: any
}

export function setPendingAnswerId(payload: any) {
  return {
    type: 'aggregation/PENDING_ANSWER_ID',
    payload,
  }
}

/**
 * aggregation/NEXT_ANSWER_ID
 */
export interface SetNextAnswerIdAction {
  type: 'aggregation/NEXT_ANSWER_ID'
  payload: any
}

export function setNextAnswerId(payload: any) {
  return {
    type: 'aggregation/NEXT_ANSWER_ID',
    payload,
  }
}

/**
 * aggregation/ORACLE_RESPONSE
 */
export interface SetOracleResponseAction {
  type: 'aggregation/ORACLE_RESPONSE'
  payload: any
}

export function setOracleResponse(payload: any) {
  return {
    type: 'aggregation/ORACLE_RESPONSE',
    payload,
  }
}

/**
 * aggregation/REQUEST_TIME
 */
export interface SetRequestTimeAction {
  type: 'aggregation/REQUEST_TIME'
  payload: any
}

export function setRequestTime(payload: any) {
  return {
    type: 'aggregation/REQUEST_TIME',
    payload,
  }
}

/**
 * aggregation/MINIMUM_RESPONSES
 */
export interface SetMinimumResponsesAction {
  type: 'aggregation/MINIMUM_RESPONSES'
  payload: any
}

export function setMinumumResponses(payload: any) {
  return {
    type: 'aggregation/MINIMUM_RESPONSES',
    payload,
  }
}

/**
 * aggregation/UPDATE_HEIGHT
 */
export interface SetUpdateHeightAction {
  type: 'aggregation/UPDATE_HEIGHT'
  payload: any
}

export function setUpdateHeight(payload: any) {
  return {
    type: 'aggregation/UPDATE_HEIGHT',
    payload,
  }
}

/**
 * aggregation/ANSWER_HISTORY
 */
export interface SetAnswersHistoryAction {
  type: 'aggregation/ANSWER_HISTORY'
  payload: any
}

export function setAnswerHistory(payload: any) {
  return {
    type: 'aggregation/ANSWER_HISTORY',
    payload,
  }
}

/**
 * aggregation/CONTRACT_ADDRESS
 */
export interface SetCurrentAddressAction {
  type: 'aggregation/CONTRACT_ADDRESS'
  payload: any
}

export function setContractAddress(payload: any) {
  return {
    type: 'aggregation/CONTRACT_ADDRESS',
    payload,
  }
}

/**
 * aggregation/OPTIONS
 */
export interface SetOptionsAction {
  type: 'aggregation/OPTIONS'
  payload: any
}

export function setOptions(payload: any) {
  return {
    type: 'aggregation/OPTIONS',
    payload,
  }
}

/**
 * aggregation/CLEAR_STATE
 */
export interface SetClearStateAction {
  type: 'aggregation/CLEAR_STATE'
  payload: any
}

export function clearState() {
  return {
    type: 'aggregation/CLEAR_STATE',
  }
}

/**
 * aggregation/ETHGAS_PRICE
 */
export interface SetEthGasPriceAction {
  type: 'aggregation/ETHGAS_PRICE'
  payload: any
}

export function setEthGasPrice(payload: any) {
  return {
    type: 'aggregation/ETHGAS_PRICE',
    payload,
  }
}
