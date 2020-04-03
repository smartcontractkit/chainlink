import { createSelector } from 'reselect'
import { FeedConfig } from 'feeds'
import { AppState } from 'state'
import { ListingAnswer } from 'state/ducks/listing/operations'

export interface ListingGroup {
  name: string
  feeds: FeedConfig[]
}

const FIAT_GROUP_NAME = 'Fiat'
const FIAT_GROUP = ['USD', 'JPY', 'GBP']

const ETH_GROUP_NAME = 'ETH'
const ETH_GROUP = ['ETH']

const GROUPS: Record<string, string[]> = {
  [FIAT_GROUP_NAME]: FIAT_GROUP,
  [ETH_GROUP_NAME]: ETH_GROUP,
}
const GROUP_ORDER: string[] = [FIAT_GROUP_NAME, ETH_GROUP_NAME]

const orderedFeeds = (state: AppState) =>
  state.feeds.order.map(f => state.feeds.items[f])

export const groups = createSelector([orderedFeeds], (feeds: FeedConfig[]) => {
  return GROUP_ORDER.map(groupName => {
    const groupFeeds = feeds.filter(f => {
      if (!f.listing) return false

      const quoteAssets = GROUPS[groupName] || []
      return quoteAssets.includes(f.pair[1])
    })
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
