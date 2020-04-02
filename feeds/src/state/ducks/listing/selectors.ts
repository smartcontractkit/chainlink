import { createSelector } from 'reselect'
import { FeedConfig } from 'feeds'
import { AppState } from 'state'
import { ListingAnswer } from 'state/ducks/listing/operations'

export interface ListingGroup {
  name: string
  feeds: FeedConfig[]
}

const GROUP_ORDER: string[] = ['USD', 'ETH']

const orderedFeeds = (state: AppState) =>
  state.feeds.order.map(f => state.feeds.items[f])

export const groups = createSelector([orderedFeeds], (feeds: FeedConfig[]) => {
  return GROUP_ORDER.map(groupName => {
    const groupFeeds = feeds.filter(
      f => f.listing && f.pair[1] === groupName && f.listing,
    )
    const group: ListingGroup = { feeds: groupFeeds, name: groupName }

    return group
  })
})

export const answer = (
  state: AppState,
  contractAddress: FeedConfig['contractAddress'],
) => {
  const listingAnswers: ListingAnswer[] = state.listing.answers || []
  return listingAnswers.find(a => a.config.contractAddress === contractAddress)
}
