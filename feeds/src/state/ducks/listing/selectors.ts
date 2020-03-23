import { createSelector } from 'reselect'
import feeds from '../../../feeds.json'

export interface ListingGroup {
  list: any[]
  name: string
}

const GROUP_ORDER: string[] = ['USD', 'ETH']

const answers = (state: any) => state.listing.answers

export const groups = createSelector([answers], (answersList: any[]) => {
  return GROUP_ORDER.map(name => {
    const list = feeds
      .filter(config => config.pair[1] === name && config.listing)
      .map(config => {
        if (answersList) {
          const addAnswer = answersList.find(
            (a: any) => a.config.name === config.name,
          )
          return addAnswer || { config }
        }

        return { config }
      })

    return {
      list,
      name,
    }
  })
})
