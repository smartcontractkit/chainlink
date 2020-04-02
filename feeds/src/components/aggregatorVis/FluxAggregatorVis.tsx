import React from 'react'
import { connect } from 'react-redux'
import { AppState } from 'state'
import { Vis } from '../vis'
import { FeedConfig } from 'feeds'
import { aggregatorSelectors } from 'state/ducks/aggregator'

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

const AggregatorVis: React.FC<StateProps> = props => <Vis {...props}></Vis>

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

export default connect(mapStateToProps)(AggregatorVis)
