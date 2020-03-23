import React, { useEffect, useRef } from 'react'
import { Icon } from 'antd'
import DeviationHistoryD3 from './DeviationGraph.d3'

function DeviationHistory({ answerHistory, config }) {
  const graph = useRef()

  useEffect(() => {
    graph.current = new DeviationHistoryD3(config)
    graph.current.build()
  }, [config])

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
