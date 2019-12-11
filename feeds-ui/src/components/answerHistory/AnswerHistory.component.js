import React, { useEffect, useState, useRef } from 'react'
import HistoryGraphD3 from './HistoryGraph.d3'
import { Icon, Switch, Tooltip } from 'antd'

function AnswerHistory({ answerHistory, options }) {
  const graph = useRef()
  const [bollinger, setBollinger] = useState(options.bollinger)

  useEffect(() => {
    graph.current = new HistoryGraphD3(options)
    graph.current.build()
  }, [options])

  useEffect(() => {
    graph.current.update(answerHistory)
  }, [answerHistory])

  function onChange() {
    setBollinger(!bollinger)
    graph.current.toggleBollinger(!bollinger)
  }

  return (
    <>
      <div className="answer-history">
        <div className="answer-history-header">
          <h2>24h Price history {!answerHistory && <Icon type="loading" />}</h2>
          {options.bollinger && (
            <div className="answer-history-options">
              <Tooltip
                title={
                  <>
                    Statistical chart characterizing the prices and volatility
                    over time.
                    <br />
                    <a
                      style={{ color: 'white', fontWeight: 'bold' }}
                      target="_BLANK"
                      rel="noopener noreferrer"
                      href={`https://en.wikipedia.org/wiki/Bollinger_Bands`}
                    >
                      Read more.
                    </a>
                  </>
                }
              >
                Moving min max averages{' '}
                <Switch defaultChecked={true} onChange={onChange} />
              </Tooltip>
            </div>
          )}
        </div>
        <div className="answer-history-graph"></div>
      </div>
    </>
  )
}

export default AnswerHistory
