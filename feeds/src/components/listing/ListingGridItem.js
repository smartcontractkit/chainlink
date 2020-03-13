import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import { Col, Popover } from 'antd'
import classNames from 'classnames'

const GRID = { xs: 24, sm: 12, md: 8 }

function ListingGridItem({ item, compareOffchain, enableHealth, healthPrice }) {
  const [statusResult, statusErrors] = status(item, healthPrice)
  const tooltipErrors = `${
    statusResult === 'error' ? ':' : ''
  } ${statusErrors.join(', ')}`
  const tooltip = `${statusResult}${tooltipErrors}`
  const classes = classNames(
    'listing-grid__item',
    healthClasses(item, statusResult, enableHealth),
  )

  return (
    <Col {...GRID} title={tooltip}>
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
    </Col>
  )
}

function CompareOffchain({ item }) {
  let content = 'No offchain comparison'

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

function Sponsored({ data }) {
  const [sliced] = useState(data.slice(0, 2))

  if (data.length <= 2) {
    return sliced.map((name, i) => [
      i > 0 && ', ',
      <span key={name}>{name}</span>,
    ])
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

function healthClasses(item, statusResult, enableHeath) {
  if (!enableHeath) {
    return
  }
  if (statusResult === 'unknown') {
    return 'listing-grid__item--health-unknown'
  }
  if (statusResult === 'error') {
    return 'listing-grid__item--health-error'
  }

  return 'listing-grid__item--health-ok'
}

function status(item, healthPrice) {
  const thresholdDiff = healthPrice * (item.config.threshold / 100)
  const thresholdMin = Math.max(healthPrice - thresholdDiff, 0)
  const thresholdMax = healthPrice + thresholdDiff
  const withinThreshold =
    item.answer >= thresholdMin && item.answer < thresholdMax

  const errors = []

  if (item.answer === undefined || healthPrice === undefined) {
    return ['unknown', errors]
  }
  if (item.answer === 0) {
    errors.push('answer price is 0')
  }
  if (!withinThreshold) {
    errors.push('reference contract price is not within threshold')
  }

  if (errors.length === 0) {
    return ['ok', errors]
  } else {
    return ['error', errors]
  }
}

const mapStateToProps = (state, ownProps) => {
  const contractAddress = ownProps.item.config.contractAddress
  const healthPrice = state.listing.healthPrices[contractAddress]
  return { healthPrice }
}

export default connect(mapStateToProps)(ListingGridItem)
