import React, { useEffect, useRef } from 'react'
import HistoryGraphD3 from './HistoryGraph.d3'
import { Icon } from 'antd'

function AnswerHistory({ answerHistory, options }) {
  let graph = useRef()

  useEffect(() => {
    graph.current = new HistoryGraphD3(options)
    graph.current.build()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    graph.current.update(answerHistory)
  }, [answerHistory])

  return (
    <>
      <div className="answer-history">
        <h3 className="answer-history-header">
          24h Price history {!answerHistory && <Icon type="loading" />}
        </h3>
        <div className="answer-history-graph"></div>
      </div>
    </>
  )
}

export default AnswerHistory
