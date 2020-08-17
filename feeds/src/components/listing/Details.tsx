import React from 'react'
import { Popover, Icon } from 'antd'
import { FeedConfig } from 'config'
import { humanizeUnixTimestamp } from '../../utils'
import { Status } from './GridItem'

interface Props {
  feed: FeedConfig
  answer?: string
  healthCheckPrice?: number
  healthCheckStatus: Status
  answerTimestamp?: number
  healthClasses?: string
}

function stopPropagation(e: React.MouseEvent) {
  e.stopPropagation()
}

type DeviationIconProps = {
  status?: string
}

const DeviationIcon: React.FC<DeviationIconProps> = ({ status }) => {
  switch (status) {
    case 'error':
      return (
        <Icon
          className="listing-grid__item--health-icon-error"
          type="close-circle"
        />
      )
    case 'ok':
      return (
        <Icon
          className="listing-grid__item--health-icon-ok"
          type="check-circle"
        />
      )
    default:
      return (
        <Icon
          className="listing-grid__item--health-icon-unknown"
          type="warning"
        />
      )
  }
}

export const DetailsContent: React.FC<Props> = ({
  feed,
  answer,
  healthCheckPrice,
  healthCheckStatus,
  answerTimestamp,
  healthClasses,
}) => (
  <div onClick={stopPropagation} className="listing-grid__item--details">
    <p>
      <span>Updated time:</span>
      {answerTimestamp && humanizeUnixTimestamp(answerTimestamp, 'LLL')}
    </p>
    <p>
      <span>On-chain answer:</span>
      {answer ? (
        <>
          {feed.valuePrefix} {answer}
        </>
      ) : (
        'Unavailable'
      )}
    </p>
    <p>
      <span>Off-chain answer:</span>
      {healthCheckPrice ? (
        <>
          {feed.valuePrefix} {healthCheckPrice}
        </>
      ) : (
        'Unavailable'
      )}
    </p>
    <p>
      <span>Deviation:</span>
      <DeviationIcon status={healthClasses} /> {healthCheckStatus.result}
      {healthCheckStatus.errors.length > 0 && (
        <div>{healthCheckStatus.errors.join(', ')}</div>
      )}
    </p>
    <CompareOffchain feed={feed} />
  </div>
)

export const Details: React.FC<Props> = props => {
  return (
    <Popover
      content={<DetailsContent {...props} />}
      title={`${props.feed.name} Details`}
    >
      <span>
        <Icon type="info-circle" />
      </span>
    </Popover>
  )
}

interface CompareOffchainProps {
  feed: FeedConfig
}

function CompareOffchain({ feed }: CompareOffchainProps) {
  const content = feed.compareOffchain ? (
    <a href={feed.compareOffchain} rel="noopener noreferrer">
      <Icon type="link" /> Compare Offchain
    </a>
  ) : (
    'No offchain comparison'
  )

  return (
    <div className="listing-grid__item--offchain-comparison">{content}</div>
  )
}
