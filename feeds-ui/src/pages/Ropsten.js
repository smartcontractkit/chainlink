import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { aggregationOperations } from 'state/ducks/aggregation'
import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import withRopsten from 'enhancers/withRopsten'
import { OracleTable } from 'components/oracleTable'

const NetworkPage = ({ initContract, clearState, options }) => {
  useEffect(() => {
    async function init() {
      await initContract(options).catch(() => {})
    }

    init()
    return () => {
      clearState()
    }
  }, [initContract, clearState, options])

  return (
    <div className="page-wrapper network-page">
      <NetworkGraph options={options} />
      <NetworkGraphInfo options={options} />
      <AnswerHistory options={options} />
      <DeviationHistory options={options} />
      <OracleTable />
    </div>
  )
}

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState,
}

export default connect(null, mapDispatchToProps)(withRopsten(NetworkPage))
