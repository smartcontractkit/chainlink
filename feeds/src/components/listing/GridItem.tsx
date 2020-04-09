import { Col, Popover, Tooltip } from 'antd'
import classNames from 'classnames'
import { FeedConfig } from 'config'
import React, { Fragment, useState } from 'react'
import { connect, MapStateToProps } from 'react-redux'
import { Link } from 'react-router-dom'
import { AppState } from 'state'
import { listingSelectors } from 'state/ducks/listing'
import { HealthCheck, ListingAnswer } from 'state/ducks/listing/reducers'

interface StateProps {
  healthCheck?: HealthCheck
  listingAnswer?: ListingAnswer
}

interface OwnProps {
  feed: FeedConfig
  compareOffchain: boolean
  enableHealth: boolean
}

interface Props extends StateProps, OwnProps {}

const GRID = { xs: 24, sm: 12, md: 8 }

export const GridItem: React.FC<Props> = ({
  feed,
  listingAnswer,
  compareOffchain,
  enableHealth,
  healthCheck,
}) => {
  const status = normalizeStatus(feed, listingAnswer, healthCheck)
  const tooltipErrors = status.errors.join(', ')
  const title = `${status.result}${tooltipErrors}`
  const classes = classNames(
    'listing-grid__item',
    healthClasses(status, enableHealth),
  )
  const sponsors = feed.sponsored ?? []

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
          {listingAnswer && (
            <>
              {feed.valuePrefix} {listingAnswer.answer}
            </>
          )}
        </div>
        {sponsors.length > 0 && (
          <>
            <div className="listing-grid__item--sponsored-title">
              Sponsored by
            </div>
            <div className="listing-grid__item--sponsored">
              <Sponsored data={sponsors} />
            </div>
          </>
        )}
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

interface SponsoredProps {
  data: string[]
}
const Sponsored: React.FC<SponsoredProps> = ({ data }) => {
  const [sliced] = useState(data.slice(0, 2))

  if (data.length <= 2) {
    return (
      <>
        {sliced.map((name, i) => (
          <Fragment key={name}>
            {i > 0 ? ', ' : ''}
            <span key={name}>{name}</span>
          </Fragment>
        ))}
      </>
    )
  }

  return (
    <Popover
      content={data.map(name => (
        <div className="listing-grid__item--sponsored-popover" key={name}>
          {name}
        </div>
      ))}
      title="Sponsored by"
    >
      {sliced.map((name, i) => [i > 0 && ', ', <span key={name}>{name}</span>])}
      , (+{data.length - 2})
    </Popover>
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
  listingAnswer?: ListingAnswer,
  healthCheck?: HealthCheck,
): Status {
  const errors: string[] = []

  if (listingAnswer === undefined || healthCheck === undefined) {
    return { result: 'unknown', errors }
  }

  const answer = parseFloat(listingAnswer.answer ?? '0')
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
  const listingAnswer = listingSelectors.answer(state, contractAddress)
  const healthCheck = state.listing.healthChecks[contractAddress]

  return {
    listingAnswer,
    healthCheck,
  }
}

export default connect(mapStateToProps)(GridItem)
