import React from 'react'
import * as jsonapi from '@chainlink/json-api-client'

export interface RedirectAction {
  type: 'REDIRECT'
}

export interface MatchRouteAction {
  type: 'MATCH_ROUTE'
  match?: {
    url: string
  }
}

export interface NotifySuccessAction {
  type: 'NOTIFY_SUCCESS'
  component: React.FC<any>
  props: any
}

export interface NotifySuccessMsgAction {
  type: 'NOTIFY_SUCCESS_MSG'
  msg: string
}

export interface NotifyErrorAction {
  type: 'NOTIFY_ERROR'
  component: React.FC<any>
  error: {
    errors: jsonapi.ErrorItem[]
  }
}

export interface NotifyErrorMsgAction {
  type: 'NOTIFY_ERROR_MSG'
  msg: string
}

export interface RequestSigninAction {
  type: 'REQUEST_SIGNIN'
}

export interface ReceiveSigninSuccessAction {
  type: 'RECEIVE_SIGNIN_SUCCESS'
}

export interface ReceiveSigninFailAction {
  type: 'RECEIVE_SIGNIN_FAIL'
}

export interface ReceiveSigninErrorAction {
  type: 'RECEIVE_SIGNIN_ERROR'
}

export interface RequestSignoutAction {
  type: 'REQUEST_SIGNOUT'
}

export interface ReceiveSignoutSuccessAction {
  type: 'RECEIVE_SIGNOUT_SUCCESS'
}

export interface ReceiveSignoutErrorAction {
  type: 'RECEIVE_SIGNOUT_ERROR'
}

export interface RequestCreateAction {
  type: 'REQUEST_CREATE'
}

export interface ReceiveCreateSuccessAction {
  type: 'RECEIVE_CREATE_SUCCESS'
}

export interface ReceiveCreateErrorAction {
  type: 'RECEIVE_CREATE_ERROR'
}

export interface RequestDeleteAction {
  type: 'REQUEST_DELETE'
}

export interface ReceiveDeleteSuccessAction {
  type: 'RECEIVE_DELETE_SUCCESS'
}

export interface ReceiveDeleteErrorAction {
  type: 'RECEIVE_DELETE_ERROR'
}

export interface RequestUpdateAction {
  type: 'REQUEST_UPDATE'
}

export interface ReceiveUpdateSuccessAction {
  type: 'RECEIVE_UPDATE_SUCCESS'
}

export interface ReceiveUpdateErrorAction {
  type: 'RECEIVE_UPDATE_ERROR'
}

export interface RequestAccountBalanceAction {
  type: 'REQUEST_ACCOUNT_BALANCE'
}

export interface UpsertAccountBalanceAction {
  type: 'UPSERT_ACCOUNT_BALANCE'
}

export interface ResponseAccountBalanceAction {
  type: 'RESPONSE_ACCOUNT_BALANCE'
}

export type Actions =
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
  | RequestSignoutAction
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
