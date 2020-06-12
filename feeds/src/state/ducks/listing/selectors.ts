import { FeedConfig } from 'config'
import { createSelector } from 'reselect'
import { AppState } from 'state'

/**
 * feed groups
 */
export interface ListingGroup {
  name: string
  feeds: FeedConfig[]
}

const FIAT_GROUP_NAME = 'Fiat'
const FIAT_GROUP = ['USD', 'JPY', 'GBP']

const ETH_GROUP_NAME = 'ETH'
const ETH_GROUP = ['ETH', 'Wei']

const GROUPS: Record<string, string[]> = {
  [FIAT_GROUP_NAME]: FIAT_GROUP,
  [ETH_GROUP_NAME]: ETH_GROUP,
}
const GROUP_ORDER: string[] = [FIAT_GROUP_NAME, ETH_GROUP_NAME]

export const feedGroups = createSelector<
  AppState,
  FeedConfig[],
  ListingGroup[]
>([orderedFeeds], listedFeeds => {
  return GROUP_ORDER.map(groupName => {
    const groupFeeds = listedFeeds.filter(f => {
      const quoteAssets = GROUPS[groupName] || []
      return quoteAssets.includes(f.pair[1])
    })

    return { feeds: groupFeeds, name: groupName }
  })
})

function feedsItems(state: AppState) {
  return state.listing.feedItems
}

function feedsOrder(state: AppState) {
  return state.listing.feedOrder
}

function orderedFeeds(state: AppState) {
  return createSelector([feedsItems, feedsOrder], (items, order) =>
    order.map(f => items[f]),
  )(state)
}

/**
 * answers
 */
export function answer(
  state: AppState,
  contractAddress: FeedConfig['contractAddress'],
) {
  return createSelector<
    AppState,
    AppState['listing']['answers'],
    string | undefined
  >(
    [listingAnswers],
    answers => answers[contractAddress],
  )(state)
}

function listingAnswers(state: AppState) {
  return state.listing.answers
}
