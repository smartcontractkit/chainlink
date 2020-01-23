import React from 'react'
import { Redirect } from 'react-router-dom'
import feeds from 'feeds.json'

/**
 * withConfig enhancer
 */

function addrCmp(a, b) {
  return a.toLowerCase() === b.toLowerCase()
}

const WithConfig = ({ render, match, networkId }) => {
  const config = feeds.find(contractConfig => {
    if (match.params.pair) {
      return (
        addrCmp(contractConfig.path, match.params.pair) &&
        contractConfig.networkId === networkId
      )
    } else if (match.params.address) {
      return addrCmp(contractConfig.contractAddress, match.params.address)
    } else {
      return false
    }
  })

  return config ? render(config) : <Redirect to={'/'} />
}

export default WithConfig
