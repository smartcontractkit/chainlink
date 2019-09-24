import React, { useEffect } from 'react'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import { aggregationOperations } from 'state/ducks/aggregation'

import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'
import { AnswerHistory } from 'components/answerHistory'

const OPTIONS = {
  contractAddress: '0x79fEbF6B9F76853EDBcBc913e6aAE8232cFB9De9',
  name: 'ETH / USD aggregation',
  valuePrefix: '$',
  answerName: 'ETH',
  counter: 300,
  network: 'mainnet',
  history: true
}

const NetworkPage = ({ initContract, clearState }) => {
  useEffect(() => {
    async function init() {
      try {
        await initContract(OPTIONS)
      } catch (error) {
        console.log(error)
      }
    }

    init()
    return () => {
      clearState()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div className="page-wrapper network-page">
      <NetworkGraph options={OPTIONS} />
      <NetworkGraphInfo options={OPTIONS} />
      <AnswerHistory options={OPTIONS} />
    </div>
  )
}

const mapStateToProps = state => ({})

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState
}

export default compose(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )
)(NetworkPage)
