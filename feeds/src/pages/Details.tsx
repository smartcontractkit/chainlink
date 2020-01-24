import React, { useEffect } from 'react'
import { connect } from 'react-redux'
/* import { DispatchBinding } from '@chainlink/ts-helpers' */
import { FeedConfig } from 'feeds'
import { aggregationOperations } from 'state/ducks/aggregation'
import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'
import { Header } from 'components/header'

interface OwnProps {
  config: FeedConfig
}

interface StateProps {}

interface DispatchProps {
  initContract: any
  clearState: any
  /* initContract: DispatchBinding<typeof aggregationOperations.initContract> */
  /* clearState: DispatchBinding<typeof aggregationOperations.clearState> */
}

interface Props extends OwnProps, StateProps, DispatchProps {}

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
        <NetworkGraph options={config} />
        <NetworkGraphInfo options={config} />
        {config.history && <AnswerHistory options={config} />}
        {config.history && <DeviationHistory options={config} />}
        <OracleTable />
      </div>
    </>
  )
}

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState,
}

export default connect(null, mapDispatchToProps)(Page)
