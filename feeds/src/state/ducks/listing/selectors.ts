import { createSelector } from 'reselect'
import feeds from '../../../feeds.json'

const answers = (state: any) => state.listing.answers

const groups = createSelector([answers], (answersList: any) => {
  const groupOrder = ['USD', 'ETH']

  const buildGroups = groupOrder.map(name => {
    const groupName = feeds
      .filter(config => config.pair[1] === name && config.listing)
      .map(config => {
        if (answersList) {
          const addAnswer = answersList.find(
            (answers: any) => answers.config.name === config.name,
          )
          return addAnswer || { config }
        }

        return { config }
      })

    return {
      list: groupName,
      name,
    }
  })

  return buildGroups
})

export { groups }
