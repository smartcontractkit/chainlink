import React from 'react'
import { connect } from 'react-redux'
import { AppState } from 'state'
import { Vis } from '../vis'
import { FeedConfig } from 'feeds'
import {
  aggregatorSelectors,
  aggregatorOperations,
} from 'state/ducks/aggregator'

interface StateProps {
  config: FeedConfig
  latestOraclesState: any
  latestAnswerTimestamp: any
  latestAnswer: any
  pendingAnswerId: any
  oracleAnswers: any
  oracleList: any
  latestRequestTimestamp: any
  minimumAnswers: any
}

interface DispatchProps {
  fetchJobId: any
}

interface Props extends StateProps, DispatchProps {}

const AggregatorVis: React.FC<Props> = props => <Vis {...props}></Vis>

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
