import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { connect, MapStateToProps } from 'react-redux'
import { Col, Tooltip } from 'antd'
import classNames from 'classnames'
import { FeedConfig } from 'config'
import { AppState } from 'state'
import { listingSelectors, listingOperations } from '../../state/ducks/listing'
import { HealthCheck } from 'state/ducks/listing/reducers'
import Sponsors from './Sponsors'

interface StateProps {
  healthCheck?: HealthCheck
  answer?: string
}

interface OwnProps {
  feed: FeedConfig
  compareOffchain: boolean
  enableHealth: boolean
}

interface DispatchProps {
  fetchAnswer: DispatchBinding<typeof listingOperations.fetchAnswer>
  fetchHealthStatus: DispatchBinding<typeof listingOperations.fetchHealthStatus>
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const GRID = { xs: 24, sm: 12, md: 8 }

export const GridItem: React.FC<Props> = ({
  feed,
  answer,
  compareOffchain,
  enableHealth,
  healthCheck,
  fetchAnswer,
  fetchHealthStatus,
}) => {
  useEffect(() => {
    fetchAnswer(feed)
  }, [fetchAnswer, feed])
  useEffect(() => {
    if (enableHealth) {
      fetchHealthStatus(feed)
    }
  }, [enableHealth, fetchHealthStatus, feed])

  const status = normalizeStatus(feed, answer, healthCheck)
  const tooltipErrors = status.errors.join(', ')
  const title = `${status.result}${tooltipErrors}`
  const classes = classNames(
    'listing-grid__item',
    healthClasses(status, enableHealth),
  )

  const gridItem = (
    <div className={classes}>
      {compareOffchain && <CompareOffchain feed={feed} />}
      <Link
        to={feed.path}
        onClick={scrollToTop}
        className="listing-grid__item--link"
      >
        <div className="listing-grid__item--name">{feed.name}</div>
        <div className="listing-grid__item--answer">
          {answer && (
            <>
              {feed.valuePrefix} {answer}
            </>
          )}
        </div>
        <Sponsors sponsors={feed.sponsored} />
      </Link>
    </div>
  )

  return (
    <Col {...GRID}>
      {enableHealth ? <Tooltip title={title}>{gridItem}</Tooltip> : gridItem}
    </Col>
  )
}

interface CompareOffchainProps {
  feed: FeedConfig
}

function CompareOffchain({ feed }: CompareOffchainProps) {
  const content = feed.compareOffchain ? (
    <a href={feed.compareOffchain} rel="noopener noreferrer">
      Compare Offchain
    </a>
  ) : (
    'No offchain comparison'
  )

  return (
    <div className="listing-grid__item--offchain-comparison">{content}</div>
  )
}

function scrollToTop() {
  window.scrollTo(0, 0)
}

interface Status {
  result: string
  errors: string[]
}

function healthClasses(status: Status, enableHeath: boolean) {
  if (!enableHeath) {
    return
  }
  if (status.result === 'unknown') {
    return 'listing-grid__item--health-unknown'
  }
  if (status.result === 'error') {
    return 'listing-grid__item--health-error'
  }

  return 'listing-grid__item--health-ok'
}

function normalizeStatus(
  feed: FeedConfig,
  rawAnswer?: string,
  healthCheck?: HealthCheck,
): Status {
  const errors: string[] = []

  if (rawAnswer === undefined || healthCheck === undefined) {
    return { result: 'unknown', errors }
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
      `reference contract price is not within threshold ${thresholdMin} - ${thresholdMax}`,
    )
  }

  if (errors.length === 0) {
    return { result: 'ok', errors }
  } else {
    return { result: 'error', errors }
  }
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, AppState> = (
  state,
  ownProps,
) => {
  const contractAddress = ownProps.feed.contractAddress
  const answer = listingSelectors.answer(state, contractAddress)
  const healthCheck = state.listing.healthChecks[contractAddress]

  return {
    answer,
    healthCheck,
  }
}

const mapDispatchToProps = {
  fetchAnswer: listingOperations.fetchAnswer,
  fetchHealthStatus: listingOperations.fetchHealthStatus,
}

export default connect(mapStateToProps, mapDispatchToProps)(GridItem)
