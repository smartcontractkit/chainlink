import React from 'react'
import { Vis } from '../vis'
import { FeedConfig } from 'feeds'

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

export interface Props extends StateProps, DispatchProps {}

const AggregatorVis: React.FC<Props> = props => <Vis {...props} />
export default AggregatorVis
