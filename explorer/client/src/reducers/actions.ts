import * as jsonapi from '@chainlink/json-api-client'

/**
 * ADMIN_SIGNIN_SUCCEEDED
 */

export interface FetchAdminSigninSucceededAction {
  type: 'FETCH_ADMIN_SIGNIN_SUCCEEDED'
  data: {
    allowed: boolean
  }
}

/**
 * FETCH_ADMIN_SIGNIN_ERROR
 */

export interface FetchAdminSigninErrorAction {
  type: 'FETCH_ADMIN_SIGNIN_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * FETCH_ADMIN_SIGNOUT_SUCCEEDED
 */

export interface FetchAdminSignoutSucceededAction {
  type: 'FETCH_ADMIN_SIGNOUT_SUCCEEDED'
}

/**
 * FETCH_ADMIN_OPERATORS_BEGIN
 */

export type FetchAdminOperatorsBeginAction = {
  type: 'FETCH_ADMIN_OPERATORS_BEGIN'
}

/**
 * FETCH_ADMIN_OPERATORS_SUCCEEDED
 */

export interface AdminOperatorsNormalizedMeta {
  currentPageOperators: {
    data: any[]
    meta: {
      count: number
    }
  }
}

export interface AdminOperatorsNormalizedData {
  chainlinkNodes: any
  meta: AdminOperatorsNormalizedMeta
}

export type FetchAdminOperatorsSucceededAction = {
  type: 'FETCH_ADMIN_OPERATORS_SUCCEEDED'
  data: AdminOperatorsNormalizedData
}

/**
 * FETCH_ADMIN_OPERATORS_ERROR
 */

export type FetchAdminOperatorsErrorAction = {
  type: 'FETCH_ADMIN_OPERATORS_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * FETCH_ADMIN_OPERATOR_BEGIN
 */

export type FetchAdminOperatorBeginAction = {
  type: 'FETCH_ADMIN_OPERATOR_BEGIN'
}

/**
 * FETCH_ADMIN_OPERATOR_SUCCEEDED
 */

export interface AdminOperatorNormalizedMeta {
  node: {
    data: any[]
  }
}

export interface AdminOperatorNormalizedData {
  chainlinkNodes: any
  meta: AdminOperatorNormalizedMeta
}

export type FetchAdminOperatorSucceededAction = {
  type: 'FETCH_ADMIN_OPERATOR_SUCCEEDED'
  data: AdminOperatorNormalizedData
}

/**
 * FETCH_ADMIN_OPERATOR_ERROR
 */

export type FetchAdminOperatorErrorAction = {
  type: 'FETCH_ADMIN_OPERATOR_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * FETCH_ADMIN_HEADS_BEGIN
 */

export type FetchAdminHeadsBeginAction = {
  type: 'FETCH_ADMIN_HEADS_BEGIN'
}

/**
 * FETCH_ADMIN_HEADS_SUCCEEDED
 */

export type FetchAdminHeadsSucceededAction = {
  type: 'FETCH_ADMIN_HEADS_SUCCEEDED'
  data: AdminHeadsNormalizedData
}

export interface AdminHeadsNormalizedMeta {
  currentPageHeads: {
    data: any[]
    meta: {
      count: number
    }
  }
}

export interface AdminHeadsNormalizedData {
  heads: any
  meta: AdminHeadsNormalizedMeta
}

/**
 * FETCH_ADMIN_HEADS_ERROR
 */

export type FetchAdminHeadsErrorAction = {
  type: 'FETCH_ADMIN_HEADS_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * FETCH_ADMIN_HEAD_BEGIN
 */

export type FetchAdminHeadBeginAction = {
  type: 'FETCH_ADMIN_HEAD_BEGIN'
}

/**
 * FETCH_ADMIN_HEAD_SUCCEEDED
 */

export interface AdminHeadNormalizedMeta {
  node: {
    data: any[]
  }
}

export interface AdminHeadNormalizedData {
  heads: any
  meta: AdminHeadNormalizedMeta
}

export type FetchAdminHeadSucceededAction = {
  type: 'FETCH_ADMIN_HEAD_SUCCEEDED'
  data: AdminHeadNormalizedData
}

/**
 * FETCH_ADMIN_HEAD_ERROR
 */

export type FetchAdminHeadErrorAction = {
  type: 'FETCH_ADMIN_HEAD_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * FETCH_JOB_RUNS_BEGIN
 */

export type FetchJobRunsBeginAction = {
  type: 'FETCH_JOB_RUNS_BEGIN'
}

/**
 * FETCH_JOB_RUNS_SUCCEEDED
 */

export interface JobRunsNormalizedMeta {
  currentPageJobRuns: {
    data: any[]
    meta: {
      count: number
    }
  }
}

export interface JobRunsNormalizedData {
  chainlinkNodes: any[]
  jobRuns: any
  meta: JobRunsNormalizedMeta
}

export type FetchJobRunsSucceededAction = {
  type: 'FETCH_JOB_RUNS_SUCCEEDED'
  data: JobRunsNormalizedData
}

/**
 * FETCH_JOB_RUNS_ERROR
 */

export type FetchJobRunsErrorAction = {
  type: 'FETCH_JOB_RUNS_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * FETCH_JOB_RUN_BEGIN
 */

export type FetchJobRunBeginAction = {
  type: 'FETCH_JOB_RUN_BEGIN'
}

/**
 * FETCH_JOB_RUN_SUCCEEDED
 */

export interface JobRunNormalizedMeta {
  jobRun: any
}

export interface JobRunNormalizedData {
  chainlinkNodes: any[]
  taskRuns: any[]
  jobRuns: any
  meta: JobRunNormalizedMeta
}

export interface FetchJobRunSucceededAction {
  type: 'FETCH_JOB_RUN_SUCCEEDED'
  data: JobRunNormalizedData
}

/**
 * FETCH_JOB_RUN_ERROR
 */

export type FetchJobRunErrorAction = {
  type: 'FETCH_JOB_RUN_ERROR'
  errors: jsonapi.ErrorItem[]
}

/**
 * QUERY_UPDATED
 */

export interface UpdateQueryAction {
  type: 'QUERY_UPDATED'
  data?: string
}

export type Actions =
  | FetchAdminSigninSucceededAction
  | FetchAdminSigninErrorAction
  | FetchAdminSignoutSucceededAction
  | FetchAdminOperatorsBeginAction
  | FetchAdminOperatorsSucceededAction
  | FetchAdminOperatorsErrorAction
  | FetchAdminOperatorBeginAction
  | FetchAdminOperatorSucceededAction
  | FetchAdminOperatorErrorAction
  | FetchAdminHeadsBeginAction
  | FetchAdminHeadsSucceededAction
  | FetchAdminHeadsErrorAction
  | FetchAdminHeadBeginAction
  | FetchAdminHeadSucceededAction
  | FetchAdminHeadErrorAction
  | FetchJobRunsBeginAction
  | FetchJobRunsSucceededAction
  | FetchJobRunsErrorAction
  | FetchJobRunBeginAction
  | FetchJobRunSucceededAction
  | FetchJobRunErrorAction
  | UpdateQueryAction
