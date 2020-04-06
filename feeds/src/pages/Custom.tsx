import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import { aggregatorOperations } from 'state/ducks/aggregator'
import { AggregatorVis } from 'components/aggregatorVis'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'
import { Header } from 'components/header'
import { parseQuery, uIntFrom } from 'utils'

interface OwnProps {
  history: any
}

interface DispatchProps {
  initContract: any
  clearState: any
}

interface Props extends OwnProps, DispatchProps {}

const Page: React.FC<Props> = ({ initContract, clearState, history }) => {
  const [config] = useState(formatConfig(parseQuery(history.location.search)))

  useEffect(() => {
    initContract(config).catch((error: Error) => {
      console.error('Could not initiate contract:', error)
    })
    return clearState
  }, [initContract, clearState, config])

  return (
    <>
      <div className="page-container-full-width">
        <Header />
      </div>
      <div className="page-wrapper network-page">
        <AggregatorVis config={config} />
        {config && config.history && <AnswerHistory config={config} />}
        {config && config.history && <DeviationHistory config={config} />}
        <OracleTable />
      </div>
    </>
  )
}

const mapDispatchToProps = {
  initContract: aggregatorOperations.initContract,
  clearState: aggregatorOperations.clearState,
}

function formatConfig(config: any) {
  return {
    ...config,
    networkId: uIntFrom(config.networkId ?? 0),
    contractVersion: 2,
    decimalPlaces: uIntFrom(config.decimalPlaces ?? 0),
    heartbeat: uIntFrom(config.heartbeat ?? 0) ?? false,
  }
}

export default connect(null, mapDispatchToProps)(Page)
