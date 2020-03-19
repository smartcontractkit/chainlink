import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { connect, MapStateToProps } from 'react-redux'
import { Col, Popover, Tooltip } from 'antd'
import classNames from 'classnames'
import { AppState } from 'state'

interface StateProps {
  healthCheck: any
}

interface OwnProps {
  item: any
  compareOffchain?: boolean
  enableHealth: boolean
}

interface Props extends StateProps, OwnProps {}

interface Status {
  result: string
  errors: string[]
}

const GRID = { xs: 24, sm: 12, md: 8 }

const GridItem: React.FC<Props> = ({
  item,
  compareOffchain,
  enableHealth,
  healthCheck,
}) => {
  const status = normalizeStatus(item, healthCheck)
  const tooltipErrors = status.errors.join(', ')
  const title = `${status.result}${tooltipErrors}`
  const classes = classNames(
    'listing-grid__item',
    healthClasses(status, enableHealth),
  )
  const gridItem = (
    <div className={classes}>
      {compareOffchain && <CompareOffchain item={item} />}
      <Link
        to={item.config.path}
        onClick={scrollToTop}
        className="listing-grid__item--link"
      >
        <div className="listing-grid__item--name">{item.config.name}</div>
        <div className="listing-grid__item--answer">
          {item.answer && (
            <>
              {item.config.valuePrefix} {item.answer}
            </>
          )}
        </div>
        {item.config.sponsored.length > 0 && (
          <>
            <div className="listing-grid__item--sponsored-title">
              Sponsored by
            </div>
            <div className="listing-grid__item--sponsored">
              <Sponsored data={item.config.sponsored} />
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

function CompareOffchain({ item }: any) {
  let content: any = 'No offchain comparison'

  if (item.config.compare_offchain) {
    content = (
      <a href={item.config.compare_offchain} rel="noopener noreferrer">
        Compare Offchain
      </a>
    )
  }

  return (
    <div className="listing-grid__item--offchain-comparison">{content}</div>
  )
}

function Sponsored({ data }: any) {
  const [sliced] = useState(data.slice(0, 2))

  if (data.length <= 2) {
    return sliced.map((name: any, i: number) => [
      i > 0 && ', ',
      <span key={name}>{name}</span>,
    ])
  }

  return (
    <Popover
      content={data.map((name: any) => (
        <div className="listing-grid__item--sponsored-popover" key={name}>
          {name}
        </div>
      ))}
      title="Sponsored by"
    >
      {sliced.map((name: any, i: number) => [
        i > 0 && ', ',
        <span key={name}>{name}</span>,
      ])}
      , (+{data.length - 2})
    </Popover>
  )
}

function scrollToTop() {
  window.scrollTo(0, 0)
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

function normalizeStatus(item: any, healthCheck: any): Status {
  const errors: string[] = []

  if (item.answer === undefined || healthCheck === undefined) {
    return { result: 'unknown', errors }
  }

  const thresholdDiff = healthCheck.currentPrice * (item.config.threshold / 100)
  const thresholdMin = Math.max(healthCheck.currentPrice - thresholdDiff, 0)
  const thresholdMax = healthCheck.currentPrice + thresholdDiff
  const withinThreshold =
    item.answer >= thresholdMin && item.answer < thresholdMax

  if (item.answer === 0) {
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
  const contractAddress = ownProps.item.config.contractAddress
  const healthCheck = state.listing.healthChecks[contractAddress]
  return { healthCheck }
}

export default connect(mapStateToProps)(GridItem)
