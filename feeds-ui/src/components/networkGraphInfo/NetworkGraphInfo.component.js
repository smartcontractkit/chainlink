import React from 'react'
import moment from 'moment'
import { Icon } from 'antd'

import CountDown from './CountDown.component'

function NetworkGraphInfo({
  currentAnswer,
  requestTime,
  nextAnswerId,
  minimumResponses,
  oracleResponse,
  oracles,
  updateHeight,
  options,
  pendingAnswerId
}) {
  const updateTime =
    updateHeight && updateHeight.timestamp
      ? moment.unix(updateHeight.timestamp).format('hh:mm:ss A')
      : '...'

  const updateDate =
    updateHeight && updateHeight.timestamp
      ? moment.unix(updateHeight.timestamp).format('MMM Do YYYY')
      : '...'

  const getCurrentResponses = () => {
    if (!oracleResponse) {
      return '...'
    }

    const responended = oracleResponse.filter(r => {
      return r.answerId >= pendingAnswerId
    })

    return `${responended.length} / ${oracles && oracles.length}`
  }

  return (
    <div className="network-graph-info__wrapper">
      <div className="network-graph-info__title">
        <h4 className="network-graph-info__title--address">
          {options.network !== 'mainnet' && (
            <div style={{ color: '#ff6300' }}>
              <Icon type="warning" /> {options.network.toUpperCase()} NETWORK
            </div>
          )}

          {options.contractAddress}
        </h4>
        <h2 className="network-graph-info__title--name">{options.name}</h2>
      </div>

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">Latest answer</div>
        <h2 className="network-graph-info__item--value">
          {options.valuePrefix || ''} {currentAnswer || '...'}
        </h2>
      </div>

      {options.counter && (
        <div className="network-graph-info__item">
          <div className="network-graph-info__item--label">
            Next aggregation starts in
          </div>
          <h2 className="network-graph-info__item--value">
            <CountDown requestTime={requestTime} counter={options.counter} />
          </h2>
        </div>
      )}

      {/* <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Current aggregation
        </div>
        <h2 className="network-graph-info__item--value">
          {currentAnswerId || '...'}
        </h2>
      </div> */}

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Oracle responses (minimum {minimumResponses || '...'})
        </div>
        <h2 className="network-graph-info__item--value">
          {/* {responses} */}
          {getCurrentResponses()}
        </h2>
      </div>

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Update date {updateDate}
        </div>
        <h2 className="network-graph-info__item--value">{updateTime}</h2>
      </div>
    </div>
  )
}

export default NetworkGraphInfo
