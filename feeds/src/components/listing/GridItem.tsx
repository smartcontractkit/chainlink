import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { connect, MapStateToProps } from 'react-redux'
import { Col } from 'antd'
import classNames from 'classnames'
import { FeedConfig } from 'config'
import { AppState } from 'state'
import { listingSelectors, listingOperations } from '../../state/ducks/listing'
import { HealthCheck } from 'state/ducks/listing/reducers'
import Sponsors from './Sponsors'
import { Details } from './Details'
import { humanizeUnixTimestamp } from '../../utils'

interface StateProps {
  healthCheck?: HealthCheck
  answer?: string
  answerTimestamp?: number
}

interface OwnProps {
  feed: FeedConfig
  enableDetails: boolean
}

interface DispatchProps {
  fetchLatestData: DispatchBinding<typeof listingOperations.fetchLatestData>
  fetchHealthStatus: DispatchBinding<typeof listingOperations.fetchHealthStatus>
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const GRID = { xs: 24, sm: 12, md: 8 }

export const GridItem: React.FC<Props> = ({
  feed,
  answer,
  answerTimestamp,
  healthCheck,
  enableDetails,
  fetchLatestData,
  fetchHealthStatus,
}) => {
  useEffect(() => {
    fetchLatestData(feed)
  }, [fetchLatestData, feed])

  useEffect(() => {
    if (enableDetails) {
      fetchHealthStatus(feed)
    }
  }, [enableDetails, fetchHealthStatus, feed])

  const healthCheckStatus = normalizeStatus(feed, answer, healthCheck)

  const classes = classNames('listing-grid__item', {
    [`listing-grid__item--health listing-grid__item--health-${healthClasses(
      healthCheckStatus,
    )}`]: enableDetails,
  })

  const gridItem = (
    <div className={classes}>
      <Link
        to={feed.path}
        onClick={scrollToTop}
        className="listing-grid__item--link"
      >
        <div className="listing-grid__item--details-icon">
          {enableDetails && (
            <Details
              feed={feed}
              healthCheckPrice={healthCheck?.currentPrice}
              healthCheckStatus={healthCheckStatus}
              answer={answer}
              answerTimestamp={answerTimestamp}
              healthClasses={healthClasses(healthCheckStatus)}
            />
          )}
        </div>
        <div className="listing-grid__item--name">{feed.name}</div>
        <div className="listing-grid__item--answer">
          {answer && (
            <>
              {feed.valuePrefix} {answer}
            </>
          )}
          {enableDetails && answerTimestamp && (
            <div> {humanizeUnixTimestamp(answerTimestamp, 'LLL')}</div>
          )}
        </div>
        <Sponsors sponsors={feed.sponsored} />
      </Link>
    </div>
  )

  return <Col {...GRID}>{gridItem}</Col>
}

function scrollToTop() {
  window.scrollTo(0, 0)
}

export interface Status {
  result: string
  errors: string[]
}

function healthClasses(status: Status) {
  if (status.result === 'Unknown') {
    return 'unknown'
  }
  if (status.result === 'Error') {
    return 'error'
  }

  return 'ok'
}

function normalizeStatus(
  feed: FeedConfig,
  rawAnswer?: string,
  healthCheck?: HealthCheck,
): Status {
  const errors: string[] = []

  if (rawAnswer === undefined || healthCheck === undefined) {
    return { result: 'Unknown', errors }
  }

  const answer = parseFloat(rawAnswer ?? '0')
  const thresholdDiff = healthCheck.currentPrice * (feed.threshold / 100)
  const thresholdMin = Math.max(healthCheck.currentPrice - thresholdDiff, 0)
  const thresholdMax = healthCheck.currentPrice + thresholdDiff
  const withinThreshold = answer >= thresholdMin && answer < thresholdMax

  if (answer === 0) {
    errors.push('answer price is 0')
  }
  if (!withinThreshold) {
    errors.push(
      `Reference contract price is not within threshold ${thresholdMin} - ${thresholdMax}`,
    )
  }

  if (errors.length === 0) {
    return { result: `OK. Within ${feed.threshold}% threshold`, errors }
  } else {
    return { result: 'Error', errors }
  }
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, AppState> = (
  state,
  ownProps,
) => {
  const contractAddress = ownProps.feed.contractAddress
  const answer = listingSelectors.answer(state, contractAddress)
  const answerTimestamp = listingSelectors.answerTimestamp(
    state,
    contractAddress,
  )
  const healthCheck = state.listing.healthChecks[contractAddress]

  return {
    answer,
    answerTimestamp,
    healthCheck,
  }
}

const mapDispatchToProps = {
  fetchLatestData: listingOperations.fetchLatestData,
  fetchHealthStatus: listingOperations.fetchHealthStatus,
}

export default connect(mapStateToProps, mapDispatchToProps)(GridItem)
