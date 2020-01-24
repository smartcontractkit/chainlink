import React from 'react'
import { Icon } from 'antd'
import CountDown from './CountDown.component'
import Percent from './Percent.component'
import Legend from './Legend.component'
import TooltipQuestion from '../shared/TooltipQuestion'
import { humanizeUnixTimestamp, networkName, Networks } from 'utils'
import { FeedConfig } from 'feeds'

interface OwnProps {
  options: FeedConfig
}

interface StateProps {
  currentAnswer: any
  requestTime: any
  minimumResponses: any
  oracleResponse: any
  oracles: any
  updateHeight: any
  pendingAnswerId: any
}

export interface Props extends OwnProps, StateProps {}

const NetworkGraphInfo: React.FC<Props> = ({
  currentAnswer,
  requestTime,
  minimumResponses,
  oracleResponse,
  oracles,
  updateHeight,
  options,
  pendingAnswerId,
}) => {
  const updateTime = updateHeight
    ? humanizeUnixTimestamp(updateHeight, 'h:mm A')
    : '...'

  const updateDate = updateHeight
    ? humanizeUnixTimestamp(updateHeight, 'MMM Do YYYY')
    : '...'

  const getCurrentResponses = () => {
    if (!oracleResponse) {
      return '...'
    }

    const responended = oracleResponse.filter((r: any) => {
      return r.answerId >= pendingAnswerId
    })

    return `${responended.length} / ${oracles && oracles.length}`
  }

  return (
    <div className="network-graph-info__wrapper">
      <div className="network-graph-info__title">
        <h4 className="network-graph-info__title--address">
          {options.networkId !== Networks.MAINNET && (
            <div style={{ color: '#ff6300' }}>
              <Icon type="warning" />{' '}
              {networkName(options.networkId).toUpperCase()} NETWORK
            </div>
          )}
          {options.contractAddress}{' '}
          <TooltipQuestion title={'Ethereum contract address'} />
        </h4>
        <h1 className="network-graph-info__title--name">
          {options.name} aggregation
        </h1>
      </div>

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Latest and trusted answer{' '}
          <TooltipQuestion
            title={`Answers are calculated in smart contract by Quickselect algorithm based on minimum ${minimumResponses} oracle answers`}
          />
        </div>
        <h2 className="network-graph-info__item--value">
          {options.valuePrefix || ''} {currentAnswer || '...'}
        </h2>
      </div>

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Deviation Threshold{' '}
          <TooltipQuestion
            title={`A new trusted answer is written when the off-chain price moves more than the deviation threshold`}
          />
        </div>
        <h2 className="network-graph-info__item--value">
          <Percent value={options.threshold} />
        </h2>
      </div>

      {options.counter && (
        <div className="network-graph-info__item">
          <div className="network-graph-info__item--label">
            Next aggregation starts in{' '}
            <TooltipQuestion
              title={`Every ${options.counter} seconds, aggregator smart contract calls oracles to get the new trusted answer`}
            />
          </div>
          <h2 className="network-graph-info__item--value">
            <CountDown requestTime={requestTime} counter={options.counter} />
          </h2>
        </div>
      )}

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Oracle responses (minimum {minimumResponses || '...'}){' '}
          <TooltipQuestion
            title={`Smart contract is connected to ${oracles &&
              oracles.length} oracles. Each aggregation requires at least ${minimumResponses} oracle responses to be able to calculate trusted answer`}
          />
        </div>
        <h2 className="network-graph-info__item--value">
          {getCurrentResponses()}
        </h2>
      </div>

      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Update date {updateDate}{' '}
          <TooltipQuestion
            title={`Date of updated smart contract with new trusted answer`}
          />
        </div>
        <h2 className="network-graph-info__item--value">{updateTime}</h2>
      </div>
      <Legend />
    </div>
  )
}

export default NetworkGraphInfo
