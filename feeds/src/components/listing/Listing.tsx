import { DispatchBinding } from '@chainlink/ts-helpers'
import { Row } from 'antd'
import React, { useEffect } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { AppState } from 'state'
import { listingOperations, listingSelectors } from '../../state/ducks/listing'
import GridItem from './GridItem'

interface OwnProps {
  enableDetails: boolean
}

interface StateProps {
  loadingFeeds: boolean
  feedGroups: listingSelectors.ListingGroup[]
}

interface DispatchProps {
  fetchFeeds: DispatchBinding<typeof listingOperations.fetchFeeds>
}

interface Props extends OwnProps, StateProps, DispatchProps {}

export const Listing: React.FC<Props> = ({
  fetchFeeds,
  loadingFeeds,
  feedGroups,
  enableDetails,
}) => {
  useEffect(() => {
    fetchFeeds()
  }, [fetchFeeds])

  let content
  if (loadingFeeds) {
    content = <>Loading Feeds...</>
  } else {
    content = (
      <div className="listing">
        {feedGroups.map(g => (
          <div className="listing-grid__group" key={g.name}>
            <h3 className="listing-grid__header">
              Decentralized Price Reference Data for {g.name} Pairs
            </h3>

            <Row gutter={18} className="listing-grid">
              {g.feeds.map(f => (
                <GridItem key={f.name} feed={f} enableDetails={enableDetails} />
              ))}
            </Row>
          </div>
        ))}
      </div>
    )
  }

  return content
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    loadingFeeds: state.listing.loadingFeeds,
    feedGroups: listingSelectors.feedGroups(state),
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchFeeds: listingOperations.fetchFeeds,
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
