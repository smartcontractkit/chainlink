import {
  FetchFeedsBeginAction,
  FetchFeedsSuccessAction,
  FetchFeedsErrorAction,
  FetchAnswerSuccessAction,
  FetchHealthPriceSuccessAction,
} from './ducks/listing/actions'
import {
  ClearStateAction,
  FetchFeedByPairBeginAction,
  FetchFeedByPairSuccessAction,
  FetchFeedByPairErrorAction,
  FetchFeedByAddressBeginAction,
  FetchFeedByAddressSuccessAction,
  FetchFeedByAddressErrorAction,
  FetchOracleNodesBeginAction,
  FetchOracleNodesSuccessAction,
  FetchOracleNodesErrorAction,
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
  SetEthGasPriceAction,
} from './ducks/aggregator/actions'
import { SetTooltipAction, SetDrawerAction } from './ducks/networkGraph/actions'

export interface InitialStateAction {
  type: 'INITIAL_STATE'
}

export type Actions =
  | InitialStateAction
  | ClearStateAction
  | FetchFeedsBeginAction
  | FetchFeedsSuccessAction
  | FetchFeedsErrorAction
  | FetchAnswerSuccessAction
  | FetchHealthPriceSuccessAction
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
  | SetTooltipAction
  | SetDrawerAction
