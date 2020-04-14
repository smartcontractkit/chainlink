import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { Row } from 'antd'
import { FeedConfig } from 'config'
import GridItem from './GridItem'
import { AppState } from 'state'
import { listingOperations, listingSelectors } from '../../state/ducks/listing'

interface Props {
  groups: listingSelectors.ListingGroup[]
  feeds: FeedConfig[]
  fetchAnswers: any
  fetchHealthStatus: any
  enableHealth: boolean
  compareOffchain: boolean
}

export const Listing: React.FC<Props> = ({
  feeds,
  fetchAnswers,
  fetchHealthStatus,
  groups,
  compareOffchain,
  enableHealth,
}) => {
  useEffect(() => {
    fetchAnswers(feeds)
  }, [fetchAnswers, feeds])
  useEffect(() => {
    if (enableHealth) {
      fetchHealthStatus(groups)
    }
  }, [enableHealth, fetchHealthStatus, groups])

  return (
    <div className="listing">
      {groups.map(group => (
        <div className="listing-grid__group" key={group.name}>
          <h3 className="listing-grid__header">
            Decentralized Price Reference Data for {group.name} Pairs
          </h3>

          <Row gutter={18} className="listing-grid">
            {group.feeds.map(f => (
              <GridItem
                key={f.name}
                feed={f}
                compareOffchain={compareOffchain}
                enableHealth={enableHealth}
              />
            ))}
          </Row>
        </div>
      ))}
    </div>
  )
}

const mapStateToProps = (state: AppState) => {
  const groups = listingSelectors.groups(state)
  const feeds = groups.flatMap(g => g.feeds)

  return { feeds, groups }
}

const mapDispatchToProps = {
  fetchAnswers: listingOperations.fetchAnswers,
  fetchHealthStatus: listingOperations.fetchHealthStatus,
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
