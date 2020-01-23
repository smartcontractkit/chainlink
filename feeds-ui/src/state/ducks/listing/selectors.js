import { createSelector } from 'reselect'
import feeds from 'feeds.json'
import { MAINNET_ID } from 'utils'

const answers = state => state.listing.answers

const groups = createSelector([answers], answersList => {
  const groupOrder = ['USD', 'ETH']

  const buildGroups = groupOrder.map(name => {
    const groupName = feeds
      .filter(
        config => config.pair[1] === name && config.networkId === MAINNET_ID,
      )
      .map(config => {
        if (answersList) {
          return answersList.filter(
            answers => answers.config.name === config.name,
          )[0]
        }
        return {
          config,
        }
      })

    return {
      list: groupName,
      name,
    }
  })

  return buildGroups
})

export { groups }
