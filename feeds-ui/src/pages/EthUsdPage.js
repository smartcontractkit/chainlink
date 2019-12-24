import React, { useEffect } from 'react'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import { aggregationOperations } from 'state/ducks/aggregation'

import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'

const OPTIONS = {
  contractAddress: '0x79fEbF6B9F76853EDBcBc913e6aAE8232cFB9De9',
  name: 'ETH / USD aggregation',
  valuePrefix: '$',
  answerName: 'ETH',
  counter: 600,
  network: 'mainnet',
  decimalPlaces: 3,
  multiply: '100000000',
  history: true,
  bollinger: false,
}

const NetworkPage = ({ initContract, clearState }) => {
  useEffect(() => {
    async function init() {
      try {
        await initContract(OPTIONS)
      } catch (error) {
        //
      }
    }

    init()
    return () => {
      clearState()
    }
  }, [initContract, clearState])

  return (
    <div className="page-wrapper network-page">
      <NetworkGraph options={OPTIONS} />
      <NetworkGraphInfo options={OPTIONS} />
      <AnswerHistory options={OPTIONS} />
      <DeviationHistory options={OPTIONS} />
      <OracleTable />
    </div>
  )
}

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState,
}

export default compose(connect(null, mapDispatchToProps))(NetworkPage)
