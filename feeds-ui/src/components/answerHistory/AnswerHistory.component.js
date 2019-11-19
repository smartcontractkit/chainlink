import React, { useEffect, useRef } from 'react'
import HistoryGraphD3 from './HistoryGraph.d3'
import { Icon } from 'antd'

function AnswerHistory({ answerHistory, options }) {
  let graph = useRef()

  useEffect(() => {
    graph.current = new HistoryGraphD3(options)
    graph.current.build()
  }, [options])

  useEffect(() => {
    graph.current.update(answerHistory)
  }, [answerHistory])

  return (
    <>
      <div className="answer-history">
        <h2 className="answer-history-header">
          24h Price history {!answerHistory && <Icon type="loading" />}
        </h2>
        <div className="answer-history-graph"></div>
      </div>
    </>
  )
}

export default AnswerHistory
