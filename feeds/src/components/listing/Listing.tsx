import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import { Row, Col, Popover } from 'antd'
import { listingOperations, listingSelectors } from '../../state/ducks/listing'

interface Props {
  fetchAnswers: any
  groups: any
  compareOffchain?: boolean
}

export function Listing({ fetchAnswers, groups, compareOffchain }: Props) {
  useEffect(() => {
    fetchAnswers()
  }, [fetchAnswers])
  return (
    <div className="listing">
      {groups.map((group: any) => (
        <div className="listing-grid__group" key={group.name}>
          <h3 className="listing-grid__header">
            Decentralized Price Reference Data for {group.name} Pairs
          </h3>
          <Row gutter={18} className="listing-grid">
            <ListingGrid list={group.list} compareOffchain={compareOffchain} />
          </Row>
        </div>
      ))}
    </div>
  )
}

const Sponsored = ({ data }: any) => {
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

const CompareOffchain = ({ item }: any) => {
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

const GRID = { xs: 24, sm: 12, md: 8 }

const ListingGridRow = ({ item, compareOffchain }: any) => (
  <Col {...GRID}>
    <div className="listing-grid__item">
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

function ListingGrid({ list, compareOffchain }: any) {
  return list.map((item: any) => (
    <ListingGridRow
      item={item}
      compareOffchain={compareOffchain}
      key={item.config.name}
    />
  ))
}

function scrollToTop() {
  window.scrollTo(0, 0)
}

const mapStateToProps = (state: any) => ({
  groups: listingSelectors.groups(state),
})

const mapDispatchToProps = {
  fetchAnswers: listingOperations.fetchAnswers,
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
