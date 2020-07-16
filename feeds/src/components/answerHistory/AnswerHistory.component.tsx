import { Icon } from 'antd'
import { FeedConfig } from 'config'
import React, { useEffect, useRef } from 'react'
import HistoryGraphD3 from './HistoryGraph.d3'

interface StateProps {
  answerHistory: any
}

interface OwnProps {
  config: FeedConfig
}

export interface Props extends StateProps, OwnProps {}

const AnswerHistory: React.FC<Props> = ({ answerHistory, config }) => {
  const graph = useRef<any>()

  useEffect(() => {
    graph.current = new HistoryGraphD3(config)
    graph.current.build()
  }, [config])

  useEffect(() => {
    graph.current.update(answerHistory)
  }, [answerHistory])

  return (
    <>
      <div className="answer-history">
        <div className="answer-history-header">
          <h2>24h Price history {!answerHistory && <Icon type="loading" />}</h2>
        </div>
        <div className="answer-history-graph"></div>
      </div>
    </>
  )
}

export default AnswerHistory
