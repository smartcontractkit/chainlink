import { SetAnswersAction, SetHealthPriceAction } from './ducks/listing/actions'
import {
  SetOracleListAction,
  SetLatestAnswerAction,
  SetLatestCompletedAnswerIdAction,
  SetPendingAnswerIdAction,
  SetNextAnswerIdAction,
  SetOracleAnswersAction,
  SetLatestRequestTimestampAction,
  SetMinimumAnswersAction,
  SetLatestAnswerTimestampAction,
  SetAnswersHistoryAction,
  SetCurrentAddressAction,
  SetConfigAction,
  SetClearStateAction,
  SetEthGasPriceAction,
} from './ducks/aggregator/actions'
import { SetTooltipAction, SetDrawerAction } from './ducks/networkGraph/actions'

export type Actions =
  | { type: 'initial_state' }
  | SetAnswersAction
  | SetHealthPriceAction
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
  | SetConfigAction
  | SetClearStateAction
  | SetEthGasPriceAction
  | SetTooltipAction
  | SetDrawerAction
