import { SetAnswersAction } from './ducks/listing/actions'
import {
  SetOraclesAction,
  SetCurrentAnswerAction,
  SetLatestCompletedAnswerIdAction,
  SetPendingAnswerIdAction,
  SetNextAnswerIdAction,
  SetOracleResponseAction,
  SetRequestTimeAction,
  SetMinimumResponsesAction,
  SetUpdateHeightAction,
  SetAnswersHistoryAction,
  SetCurrentAddressAction,
  SetOptionsAction,
  SetClearStateAction,
  SetEthGasPriceAction,
} from './ducks/aggregation/actions'
import { SetTooltipAction, SetDrawerAction } from './ducks/networkGraph/actions'

export type Actions =
  | SetAnswersAction
  | SetOraclesAction
  | SetCurrentAnswerAction
  | SetLatestCompletedAnswerIdAction
  | SetPendingAnswerIdAction
  | SetNextAnswerIdAction
  | SetOracleResponseAction
  | SetRequestTimeAction
  | SetMinimumResponsesAction
  | SetUpdateHeightAction
  | SetAnswersHistoryAction
  | SetCurrentAddressAction
  | SetOptionsAction
  | SetClearStateAction
  | SetEthGasPriceAction
  | SetTooltipAction
  | SetDrawerAction
