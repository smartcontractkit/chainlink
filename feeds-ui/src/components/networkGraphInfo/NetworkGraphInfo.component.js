import React from 'react'
import moment from 'moment'

import CountDown from './CountDown.component'

function NetworkGraphInfo({
  currentAnswer,
  requestTime,
  nextAnswerId,
  minimumResponses,
  oracleResponse,
  oracles,
  updateHeight
}) {
  const updateBlock = updateHeight && updateHeight.block | '...'
  const updateTime =
    updateHeight && updateHeight.timestamp
      ? moment.unix(updateHeight.timestamp).format('hh:mm:ss A')
      : '...'

  const currentAnswerId = nextAnswerId && nextAnswerId - 1
  const responses =
    (oracleResponse &&
      oracles &&
      `${oracleResponse && oracleResponse.length} / ${oracles &&
        oracles.length}`) ||
    '...'

  return (
    <div className="network-graph-info__wrapper">
      <div className="network-graph-info">
        <div className="network-graph-info__label">Current ETH price</div>
        <h2 className="network-graph-info__value">
          $ {currentAnswer || '...'}
        </h2>
      </div>

      <div className="network-graph-info">
        <div className="network-graph-info__label">
          Next aggregation starts in
        </div>
        <h2 className="network-graph-info__value">
          <CountDown requestTime={requestTime} />
        </h2>
      </div>

      <div className="network-graph-info">
        <div className="network-graph-info__label">Current aggregation</div>
        <h2 className="network-graph-info__value">
          {currentAnswerId || '...'}
        </h2>
      </div>

      <div className="network-graph-info">
        <div className="network-graph-info__label">
          Oracle responses (minimum {minimumResponses || '...'})
        </div>
        <h2 className="network-graph-info__value">{responses}</h2>
      </div>

      <div className="network-graph-info">
        <div className="network-graph-info__label">
          Update time ({updateBlock} block)
        </div>
        <h2 className="network-graph-info__value">{updateTime}</h2>
      </div>
    </div>
  )
}

export default NetworkGraphInfo
