import React from 'react'
import { Action } from 'redux'
import * as jsonapi from 'utils/json-api-client'

export interface InitialStateAction extends Action<'INITIAL_STATE'> {}

export enum RouterActionType {
  REDIRECT = 'REDIRECT',
  MATCH_ROUTE = 'MATCH_ROUTE',
}

/**
 * REDIRECT
 */

export interface RedirectAction extends Action<RouterActionType.REDIRECT> {
  to: string
}

/**
 * MATCH_ROUTE
 */

export interface MatchRouteAction extends Action<RouterActionType.MATCH_ROUTE> {
  pathname: string
}

export enum NotifyActionType {
  NOTIFY_SUCCESS = 'NOTIFY_SUCCESS',
  NOTIFY_SUCCESS_MSG = 'NOTIFY_SUCCESS_MSG',
  NOTIFY_ERROR = 'NOTIFY_ERROR',
  NOTIFY_ERROR_MSG = 'NOTIFY_ERROR_MSG',
}

/**
 * NOTIFY_SUCCESS
 */

export interface NotifySuccessAction
  extends Action<NotifyActionType.NOTIFY_SUCCESS> {
  component: React.FC<any>
  props: any
}

/**
 * NOTIFY_SUCCESS_MSG
 */

export interface NotifySuccessMsgAction
  extends Action<NotifyActionType.NOTIFY_SUCCESS_MSG> {
  msg: string
}

/**
 * NOTIFY_ERROR
 */

export interface NotifyErrorAction
  extends Action<NotifyActionType.NOTIFY_ERROR> {
  component: React.FC<any>
  error: {
    errors: jsonapi.ErrorItem[]
  }
}

/**
 * NOTIFY_ERROR_MSG
 */

export interface NotifyErrorMsgAction
  extends Action<NotifyActionType.NOTIFY_ERROR_MSG> {
  msg: string
}

export enum AuthActionType {
  REQUEST_SIGNIN = 'REQUEST_SIGNIN',
  RECEIVE_SIGNIN_SUCCESS = 'RECEIVE_SIGNIN_SUCCESS',
  RECEIVE_SIGNIN_FAIL = 'RECEIVE_SIGNIN_FAIL',
  RECEIVE_SIGNIN_ERROR = 'RECEIVE_SIGNIN_ERROR',
  RECEIVE_SIGNOUT_SUCCESS = 'RECEIVE_SIGNOUT_SUCCESS',
  RECEIVE_SIGNOUT_ERROR = 'RECEIVE_SIGNOUT_ERROR',
}

/**
 * REQUEST_SIGNIN
 */

export interface RequestSigninAction
  extends Action<AuthActionType.REQUEST_SIGNIN> {}

/**
 * RECEIVE_SIGNIN_SUCCESS
 */

export interface ReceiveSigninSuccessAction
  extends Action<AuthActionType.RECEIVE_SIGNIN_SUCCESS> {
  authenticated: boolean
}

/**
 * RECEIVE_SIGNIN_FAIL
 */

export interface ReceiveSigninFailAction
  extends Action<AuthActionType.RECEIVE_SIGNIN_FAIL> {}

/**
 * RECEIVE_SIGNIN_ERROR
 */

export interface ReceiveSigninErrorAction
  extends Action<AuthActionType.RECEIVE_SIGNIN_ERROR> {
  errors: any[]
}

/**
 * RECEIVE_SIGNOUT_SUCCESS
 */

export interface ReceiveSignoutSuccessAction
  extends Action<AuthActionType.RECEIVE_SIGNOUT_SUCCESS> {
  authenticated: boolean
}

/**
 * RECEIVE_SIGNOUT_ERROR
 */

export interface ReceiveSignoutErrorAction
  extends Action<AuthActionType.RECEIVE_SIGNOUT_ERROR> {
  errors: any[]
}

export enum ResourceActionType {
  DELETE_JOB_SPEC_ERROR = 'DELETE_JOB_SPEC_ERROR',
  RECEIVE_CREATE_ERROR = 'RECEIVE_CREATE_ERROR',
  RECEIVE_CREATE_SUCCESS = 'RECEIVE_CREATE_SUCCESS',
  RECEIVE_DELETE_ERROR = 'RECEIVE_DELETE_ERROR',
  RECEIVE_DELETE_SUCCESS = 'RECEIVE_DELETE_SUCCESS',
  RECEIVE_UPDATE_ERROR = 'RECEIVE_UPDATE_ERROR',
  RECEIVE_UPDATE_SUCCESS = 'RECEIVE_UPDATE_SUCCESS',
  RESPONSE_ACCOUNT_BALANCE = 'RESPONSE_ACCOUNT_BALANCE',
  REQUEST_ACCOUNT_BALANCE = 'REQUEST_ACCOUNT_BALANCE',
  REQUEST_CREATE = 'REQUEST_CREATE',
  REQUEST_DELETE = 'REQUEST_DELETE',
  REQUEST_UPDATE = 'REQUEST_UPDATE',
  UPSERT_ACCOUNT_BALANCE = 'UPSERT_ACCOUNT_BALANCE',
  UPSERT_BRIDGES = 'UPSERT_BRIDGES',
  UPSERT_BRIDGE = 'UPSERT_BRIDGE',
  UPSERT_CONFIGURATION = 'UPSERT_CONFIGURATION',
  UPSERT_JOB = 'UPSERT_JOB',
  UPSERT_JOBS = 'UPSERT_JOBS',
  UPSERT_JOB_RUN = 'UPSERT_JOB_RUN',
  UPSERT_JOB_RUNS = 'UPSERT_JOB_RUNS',
  UPSERT_RECENT_JOB_RUNS = 'UPSERT_RECENT_JOB_RUNS',
  UPSERT_RECENTLY_CREATED_JOBS = 'UPSERT_RECENTLY_CREATED_JOBS',
  UPSERT_TRANSACTION = 'UPSERT_TRANSACTION',
  UPSERT_TRANSACTIONS = 'UPSERT_TRANSACTIONS',
}

/**
 * REQUEST_CREATE
 */

export interface RequestCreateAction
  extends Action<ResourceActionType.REQUEST_CREATE> {}

/**
 * REQUEST_CREATE_SUCCESS
 */

export interface ReceiveCreateSuccessAction
  extends Action<ResourceActionType.RECEIVE_CREATE_SUCCESS> {}

/**
 * REQUEST_CREATE_ERROR
 */

export interface ReceiveCreateErrorAction
  extends Action<ResourceActionType.RECEIVE_CREATE_ERROR> {}

/**
 * REQUEST_DELETE
 */

export interface RequestDeleteAction
  extends Action<ResourceActionType.REQUEST_DELETE> {}

/**
 * RECEIVE_DELETE_SUCCESS
 */

export interface ReceiveDeleteSuccessAction
  extends Action<ResourceActionType.RECEIVE_DELETE_SUCCESS> {
  id: string
}

/**
 * RECEIVE_DELETE_ERROR
 */

export interface ReceiveDeleteErrorAction
  extends Action<ResourceActionType.RECEIVE_DELETE_ERROR> {}

/**
 * REQUEST_UPDATE
 */

export interface RequestUpdateAction
  extends Action<ResourceActionType.REQUEST_UPDATE> {}

/**
 * RECEIVE_UPDATE_SUCCESS
 */

export interface ReceiveUpdateSuccessAction
  extends Action<ResourceActionType.RECEIVE_UPDATE_SUCCESS> {}

/**
 * RECEIVE_UPDATE_ERROR
 */

export interface ReceiveUpdateErrorAction
  extends Action<ResourceActionType.RECEIVE_UPDATE_ERROR> {}

/**
 * REQUEST_ACCOUNT_BALANCE
 */

export interface RequestAccountBalanceAction
  extends Action<ResourceActionType.REQUEST_ACCOUNT_BALANCE> {}

/**
 * UPSERT_ACCOUNT_BALANCE
 */

export interface UpsertAccountBalanceAction
  extends Action<ResourceActionType.UPSERT_ACCOUNT_BALANCE> {
  data: {
    eThKeys: any
  }
}

/**
 * RESPONSE_ACCOUNT_BALANCE
 */

export interface ResponseAccountBalanceAction
  extends Action<ResourceActionType.RESPONSE_ACCOUNT_BALANCE> {}

/**
 * UPSERT_BRIDGES
 */

export interface UpsertBridgesAction
  extends Action<ResourceActionType.UPSERT_BRIDGES> {
  data: {
    bridges: Record<string, any>
    meta: {
      currentPageBridges: {
        data: { id: string }[]
        meta: { count: number }
      }
    }
  }
}

export interface UpsertBridgeAction
  extends Action<ResourceActionType.UPSERT_BRIDGE> {
  data: {
    bridges: Record<string, any>
  }
}

export type ConfigurationAttribute = string | number | null

export interface UpsertConfigurationAction
  extends Action<ResourceActionType.UPSERT_CONFIGURATION> {
  data: {
    configPrinters: Record<string, any>
  }
}

export interface UpsertJobsAction
  extends Action<ResourceActionType.UPSERT_JOBS> {
  data: {
    specs: Record<string, any>
    meta: {
      currentPageJobs: {
        data: { id: string }[]
        meta: { count: number }
      }
    }
  }
}

export interface UpsertRecentlyCreatedJobsAction
  extends Action<ResourceActionType.UPSERT_RECENTLY_CREATED_JOBS> {
  data: {
    specs: Record<string, any>
    meta: {
      recentlyCreatedJobs: {
        data: { id: string }[]
      }
    }
  }
}

export interface UpsertJobAction extends Action<ResourceActionType.UPSERT_JOB> {
  data: any
}

export interface UpsertJobRunsAction
  extends Action<ResourceActionType.UPSERT_JOB_RUNS> {
  data: {
    runs: Record<string, any>
    meta: {
      currentPageJobRuns: {
        data: { id: string }[]
        meta: {
          count: number
        }
      }
    }
  }
}

export interface UpsertRecentJobRunsAction
  extends Action<ResourceActionType.UPSERT_RECENT_JOB_RUNS> {
  data: {
    runs: Record<string, any>
    meta: {
      recentJobRuns: {
        data: { id: string }[]
        meta: {
          count: number
        }
      }
    }
  }
}

export interface UpsertJobRunAction
  extends Action<ResourceActionType.UPSERT_JOB_RUN> {
  data: {
    runs: Record<string, any>
  }
}

export interface UpsertTransactionsAction
  extends Action<ResourceActionType.UPSERT_TRANSACTIONS> {
  data: {
    transactions: Record<string, any>
    meta: {
      currentPageTransactions: {
        data: { id: string }[]
        meta: {
          count: number
        }
      }
    }
  }
}

export interface UpsertTransactionAction
  extends Action<ResourceActionType.UPSERT_TRANSACTION> {
  data: {
    transactions: Record<string, any>
  }
}

export interface DeleteJobSpecError
  extends Action<ResourceActionType.DELETE_JOB_SPEC_ERROR> {
  data: {
    id: string
    jobSpecID: string
  }
}

export type Actions =
  | InitialStateAction
  | RedirectAction
  | MatchRouteAction
  | NotifySuccessAction
  | NotifySuccessMsgAction
  | NotifyErrorAction
  | NotifyErrorMsgAction
  | RequestSigninAction
  | ReceiveSigninSuccessAction
  | ReceiveSigninFailAction
  | ReceiveSigninErrorAction
  | ReceiveSignoutSuccessAction
  | ReceiveSignoutErrorAction
  | RequestCreateAction
  | ReceiveCreateSuccessAction
  | ReceiveCreateErrorAction
  | RequestDeleteAction
  | ReceiveDeleteSuccessAction
  | ReceiveDeleteErrorAction
  | RequestUpdateAction
  | ReceiveUpdateSuccessAction
  | ReceiveUpdateErrorAction
  | RequestAccountBalanceAction
  | UpsertAccountBalanceAction
  | ResponseAccountBalanceAction
  | UpsertBridgesAction
  | UpsertBridgeAction
  | UpsertConfigurationAction
  | UpsertJobsAction
  | UpsertRecentlyCreatedJobsAction
  | UpsertJobAction
  | UpsertJobRunsAction
  | UpsertRecentJobRunsAction
  | UpsertJobRunAction
  | UpsertTransactionsAction
  | UpsertTransactionAction
  | DeleteJobSpecError
