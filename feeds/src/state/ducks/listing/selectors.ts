import { createSelector } from 'reselect'
import { Networks } from '../../../utils'
import feeds from '../../../feeds.json'

const answers = (state: any) => state.listing.answers

const groups = createSelector([answers], (answersList: any) => {
  const groupOrder = ['USD', 'ETH']

  const buildGroups = groupOrder.map(name => {
    const groupName = feeds
      .filter(
        config =>
          config.pair[1] === name && config.networkId === Networks.MAINNET,
      )
      .map(config => {
        if (answersList) {
          return answersList.filter(
            (answers: any) => answers.config.name === config.name,
          )[0]
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
