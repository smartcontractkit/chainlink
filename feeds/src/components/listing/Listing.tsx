import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { Row } from 'antd'
import { AppState } from 'state'
import { listingOperations, listingSelectors } from '../../state/ducks/listing'
import GridItem from './GridItem'

interface Props {
  groups: any[]
  fetchAnswers: any
  fetchHealthStatus: any
  enableHealth: boolean
  compareOffchain?: boolean
}

export const Listing: React.FC<Props> = ({
  fetchAnswers,
  fetchHealthStatus,
  groups,
  compareOffchain,
  enableHealth,
}) => {
  useEffect(() => {
    fetchAnswers()
  }, [fetchAnswers])
  useEffect(() => {
    if (enableHealth) {
      fetchHealthStatus(groups)
    }
  }, [enableHealth, fetchHealthStatus, groups])

  return (
    <div className="listing">
      {groups.map((group: any) => (
        <div className="listing-grid__group" key={group.name}>
          <h3 className="listing-grid__header">
            Decentralized Price Reference Data for {group.name} Pairs
          </h3>

          <Row gutter={18} className="listing-grid">
            {group.list.map((item: any) => (
              <GridItem
                key={item.config.name}
                item={item}
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

const mapStateToProps = (state: AppState) => ({
  groups: listingSelectors.groups(state),
})

const mapDispatchToProps = {
  fetchAnswers: listingOperations.fetchAnswers,
  fetchHealthStatus: listingOperations.fetchHealthStatus,
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
