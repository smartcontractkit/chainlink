import React from 'react'
import { Graph } from './graph'
import { Info } from './info'

export interface OwnProps {
  config: any
}

interface StateGraphProps {
  latestOraclesState: any
  latestAnswer: any
  latestAnswerTimestamp: any
}

interface StateInfoProps {
  pendingAnswerId: any
  oracleAnswers: any
  oracleList: any
  latestRequestTimestamp: any
  minimumAnswers: any
}

interface DispatchProps {
  fetchJobId?: any
}

export interface Props
  extends OwnProps,
    StateGraphProps,
    StateInfoProps,
    DispatchProps {}

const Vis: React.FC<Props> = ({
  config,
  latestOraclesState,
  latestAnswer,
  latestAnswerTimestamp,
  pendingAnswerId,
  oracleAnswers,
  oracleList,
  latestRequestTimestamp,
  minimumAnswers,
  fetchJobId,
}) => {
  return (
    <>
      <Graph
        latestOraclesState={latestOraclesState}
        latestAnswer={latestAnswer}
        latestAnswerTimestamp={latestAnswerTimestamp}
        config={config}
        fetchJobId={fetchJobId}
      />
      <Info
        latestAnswer={latestAnswer}
        latestAnswerTimestamp={latestAnswerTimestamp}
        latestRequestTimestamp={latestRequestTimestamp}
        minimumAnswers={minimumAnswers}
        oracleAnswers={oracleAnswers}
        config={config}
        pendingRoundId={pendingAnswerId}
        oracleList={oracleList}
      />
    </>
  )
}

export default Vis
