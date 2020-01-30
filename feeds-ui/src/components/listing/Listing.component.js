import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Row, Col, Popover } from 'antd'

const scrollToTop = () => {
  window.scrollTo(0, 0)
}

const grid = { xs: 24, sm: 12, md: 8 }

const Supported = ({ data }) => {
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
        <div className="listing-grid__item--supported-popover" key={name}>
          {name}
        </div>
      ))}
      title="Supported by"
    >
      {sliced.map((name, i) => [i > 0 && ', ', <span key={name}>{name}</span>])}
      ...
    </Popover>
  )
}

const ListingGridRow = ({ item }) => (
  <Col {...grid}>
    <Link to={item.config.path} onClick={scrollToTop}>
      <div className="listing-grid__item">
        <div className="listing-grid__item--name">{item.config.name}</div>
        <div className="listing-grid__item--answer">
          {item.answer && (
            <>
              {item.config.valuePrefix} {item.answer}
            </>
          )}
        </div>
        {item.config.supported.length > 0 && (
          <>
            <div className="listing-grid__item--supported-title">
              Supported by
            </div>
            <div className="listing-grid__item--supported">
              <Supported data={item.config.supported} />
            </div>
          </>
        )}
      </div>
    </Link>
  </Col>
)

const ListingGrid = ({ list }) => (
  <>
    {list.map(item => (
      <ListingGridRow item={item} key={item.config.name} />
    ))}
  </>
)

const Listing = ({ fetchAnswers, groups }) => {
  useEffect(() => {
    fetchAnswers()
  }, [fetchAnswers])
  return (
    <div className="listing">
      {groups.map(group => (
        <div className="listing-grid__group" key={group.name}>
          <h3 className="listing-grid__header">
            Decentralized Price Reference Data for {group.name} Pairs
          </h3>
          <Row gutter={18} className="listing-grid">
            <ListingGrid list={group.list} />
          </Row>
        </div>
      ))}
    </div>
  )
}

export default Listing
