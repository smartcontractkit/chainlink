import { FeedConfig } from 'config'
import React from 'react'
import { humanizeUnixTimestamp } from 'utils'
import TooltipQuestion from '../../shared/TooltipQuestion'
import Heartbeat from './Heartbeat'
import Legend from './Legend'
import Percent from './Percent'

interface OwnProps {
  config: FeedConfig
}

interface StateProps {
  latestAnswer: string
  latestRequestTimestamp: number
  minimumAnswers: number
  oracleAnswers: any
  oracleList: Array<string>
  latestAnswerTimestamp: any
  pendingRoundId: any
}

export interface Props extends OwnProps, StateProps {}

const Info: React.FC<Props> = ({
  latestAnswer,
  latestAnswerTimestamp,
  latestRequestTimestamp,
  minimumAnswers,
  oracleAnswers,
  oracleList,
  config,
  pendingRoundId,
}) => {
  const updateTime = latestAnswerTimestamp
    ? humanizeUnixTimestamp(latestAnswerTimestamp, 'h:mm A')
    : '...'

  const updateDate = latestAnswerTimestamp
    ? humanizeUnixTimestamp(latestAnswerTimestamp, 'MMM Do YYYY')
    : '...'

  const getCurrentResponses = () => {
    if (!oracleAnswers) {
      return '...'
    }

    const responded = oracleAnswers.filter((r: any) => {
      return r.answerId >= pendingRoundId
    })

    return oracleList && oracleList.length
      ? `${responded.length} / ${oracleList.length}`
      : '...'
  }

  return (
    <div className="network-graph-info__wrapper">
      <div className="network-graph-info__title">
        <h4 className="network-graph-info__title--address">
          {config.proxyAddress && (
            <>
              Feed: {config.proxyAddress}{' '}
              <TooltipQuestion title={'Ethereum contract address'} />
              <br />
            </>
          )}
          Aggregator: {config.contractAddress}{' '}
          <TooltipQuestion title={'Ethereum contract address'} />
        </h4>
        <h1 className="network-graph-info__title--name">
          {config.name} aggregation
        </h1>
      </div>
      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Latest and trusted answer{' '}
          <TooltipQuestion
            title={`Answers are calculated in smart contract by Quickselect algorithm based on minimum ${minimumAnswers} oracle answers`}
          />
        </div>
        <h2 className="network-graph-info__item--value">
          {config.valuePrefix || ''} {latestAnswer || '...'}
        </h2>
      </div>
      {config.threshold && (
        <div className="network-graph-info__item">
          <div className="network-graph-info__item--label">
            Primary Aggregation Parameter{' '}
            <TooltipQuestion
              title={`A new trusted answer is written when the off-chain price moves more than the deviation threshold`}
            />
          </div>
          <h2 className="network-graph-info__item--value">
            Deviation Threshold: <Percent value={config.threshold} />
          </h2>
        </div>
      )}
      {config.heartbeat ? (
        <div className="network-graph-info__item">
          <div className="network-graph-info__item--label">
            Secondary Aggregation Parameter{' '}
            <TooltipQuestion
              title={`Every ${config.heartbeat} seconds, aggregator smart contract calls oracles to get the new trusted answer`}
            />
          </div>
          <h2 className="network-graph-info__item--value">
            Heartbeat:{' '}
            <Heartbeat
              latestRequestTimestamp={latestRequestTimestamp}
              heartbeat={config.heartbeat}
            />
          </h2>
        </div>
      ) : null}
      <div className="network-graph-info__item">
        <div className="network-graph-info__item--label">
          Oracle responses (minimum {minimumAnswers || '...'}){' '}
          <TooltipQuestion
            title={`Smart contract is connected to ${oracleList &&
              oracleList.length} oracles. Each aggregation requires at least ${minimumAnswers} oracle responses to be able to calculate trusted answer`}
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

export default Info
