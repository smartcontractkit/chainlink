import AggregatorVis from './AggregatorVis.component'
import { connect } from 'react-redux'
import { AppState } from 'state'

import {
  aggregatorSelectors,
  aggregatorOperations,
} from 'state/ducks/aggregator'

const mapStateToProps = (state: AppState) => ({
  latestOraclesState: aggregatorSelectors.latestOraclesState(state),
  latestAnswer: state.aggregator.latestAnswer,
  latestAnswerTimestamp: state.aggregator.latestAnswerTimestamp,
  pendingAnswerId: state.aggregator.pendingAnswerId,
  oracleAnswers: state.aggregator.oracleAnswers,
  oracleList: state.aggregator.oracleList,
  latestRequestTimestamp: state.aggregator.latestRequestTimestamp,
  minimumAnswers: state.aggregator.minimumAnswers,
})

const mapDispatchToProps = {
  fetchJobId: aggregatorOperations.fetchJobId,
}

export default connect(mapStateToProps, mapDispatchToProps)(AggregatorVis)
