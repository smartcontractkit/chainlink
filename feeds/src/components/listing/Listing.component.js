import React, { useEffect } from 'react'
import { Row } from 'antd'
import ListingGridItem from './ListingGridItem'

const ListingGrid = ({ list, compareOffchain, health }) => (
  <>
    {list.map(item => (
      <ListingGridItem
        key={item.config.name}
        item={item}
        compareOffchain={compareOffchain}
        enableHealth={health}
      />
    ))}
  </>
)

const Listing = ({
  fetchAnswers,
  fetchHealthStatus,
  groups,
  compareOffchain,
  health,
}) => {
  useEffect(() => {
    fetchAnswers()
    fetchHealthStatus(groups)
  }, [fetchAnswers, fetchHealthStatus, groups])

  return (
    <div className="listing">
      {groups.map(group => (
        <div className="listing-grid__group" key={group.name}>
          <h3 className="listing-grid__header">
            Decentralized Price Reference Data for {group.name} Pairs
          </h3>
          <Row gutter={18} className="listing-grid">
            <ListingGrid
              list={group.list}
              compareOffchain={compareOffchain}
              health={health}
            />
          </Row>
        </div>
      ))}
    </div>
  )
}

export default Listing
