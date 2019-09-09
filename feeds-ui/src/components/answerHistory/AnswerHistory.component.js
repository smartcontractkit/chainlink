import React, { useEffect } from 'react'
import { createChart, update } from './historyGraph'

function AnswerHistory({ answerHistory }) {
  useEffect(() => {
    createChart()
  }, [])

  useEffect(() => {
    update(answerHistory)
  }, [answerHistory])

  return (
    <>
      <div className="answer-history">
        <h3 className="answer-history-header">24h Price history</h3>
        <div className="answer-history-graph"></div>
      </div>
    </>
  )
}

export default AnswerHistory
