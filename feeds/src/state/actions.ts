import {
  FetchFeedsBeginAction,
  FetchFeedsSuccessAction,
  FetchFeedsErrorAction,
  FetchAnswerSuccessAction,
  FetchHealthPriceSuccessAction,
} from './ducks/listing/actions'
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

export interface InitialStateAction {
  type: 'INITIAL_STATE'
}

export type Actions =
  | InitialStateAction
  | FetchFeedsBeginAction
  | FetchFeedsSuccessAction
  | FetchFeedsErrorAction
  | FetchAnswerSuccessAction
  | FetchHealthPriceSuccessAction
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
