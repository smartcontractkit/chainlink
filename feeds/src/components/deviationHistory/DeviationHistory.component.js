import React, { useEffect, useRef } from 'react'
import DeviationHistoryD3 from './DeviationGraph.d3'
import { Icon } from 'antd'

function DeviationHistory({ answerHistory, options }) {
  const graph = useRef()

  useEffect(() => {
    graph.current = new DeviationHistoryD3(options)
    graph.current.build()
  }, [options])

  useEffect(() => {
    graph.current.update(answerHistory)
  }, [answerHistory])

  return (
    <>
      <div className="deviation-history">
        <h2 className="deviation-history-header">
          24h Volatility {!answerHistory && <Icon type="loading" />}
        </h2>
        <div className="deviation-history-graph"></div>
      </div>
    </>
  )
}

export default DeviationHistory
