import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { FeedConfig } from 'feeds'
import { fluxAggregatorOperations } from 'state/ducks/aggregator'
import { FluxAggregatorVis } from 'components/aggregatorVis'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'
import { Header } from 'components/header'

interface OwnProps {
  config: FeedConfig
}

interface DispatchProps {
  initContract: any
  clearState: any
}

interface Props extends OwnProps, DispatchProps {}

const Page: React.FC<Props> = ({ initContract, clearState, config }) => {
  useEffect(() => {
    initContract(config).catch(() => {
      console.error('Could not initiate contract')
    })
    return () => {
      clearState()
    }
  }, [initContract, clearState, config])

  return (
    <>
      <div className="page-container-full-width">
        <Header />
      </div>
      <div className="page-wrapper network-page">
        <FluxAggregatorVis config={config} />
        {config.history && <AnswerHistory config={config} />}
        {config.history && <DeviationHistory config={config} />}
        <OracleTable />
      </div>
    </>
  )
}

const mapDispatchToProps = {
  initContract: fluxAggregatorOperations.initContract,
  clearState: fluxAggregatorOperations.clearState,
}

export default connect(null, mapDispatchToProps)(Page)
