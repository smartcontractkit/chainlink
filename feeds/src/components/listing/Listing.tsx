import { DispatchBinding } from '@chainlink/ts-helpers'
import { Row } from 'antd'
import React, { useEffect } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { AppState } from 'state'
import { listingOperations, listingSelectors } from '../../state/ducks/listing'
import GridItem from './GridItem'

interface OwnProps {
  enableHealth: boolean
  compareOffchain: boolean
}

interface DispatchProps {
  fetchAnswers: DispatchBinding<typeof listingOperations.fetchAnswers>
  fetchHealthStatus: DispatchBinding<typeof listingOperations.fetchHealthStatus>
}

interface StateProps {
  groups: ReturnType<typeof listingSelectors.groups>
}

type Props = OwnProps & DispatchProps & StateProps

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

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => ({
  groups: listingSelectors.groups(state),
})

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchAnswers: listingOperations.fetchAnswers,
  fetchHealthStatus: listingOperations.fetchHealthStatus,
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
