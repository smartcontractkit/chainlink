import { partialAsFull } from '@chainlink/ts-helpers'
import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import { aggregatorOperations } from 'state/ducks/aggregator'
import { AggregatorVis } from 'components/aggregatorVis'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'
import { Header } from 'components/header'
import { parseQuery, uIntFrom } from 'utils'
import { FeedConfig } from 'config'

interface OwnProps {
  history: any
}

interface DispatchProps {
  initContract: any
}

interface Props extends OwnProps, DispatchProps {}

const Page: React.FC<Props> = ({ initContract, history }) => {
  const [config] = useState(formatConfig(parseQuery(history.location.search)))

  useEffect(() => {
    initContract(config).catch(() => {
      console.error('Could not initiate contract')
    })
  }, [initContract, config])

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
}

function formatConfig(queryConfig: Record<string, string>): FeedConfig {
  return partialAsFull<FeedConfig>({
    ...queryConfig,
    networkId: uIntFrom(queryConfig.networkId ?? 0),
    contractVersion: 2,
    decimalPlaces: uIntFrom(queryConfig.decimalPlaces ?? 0),
    heartbeat: uIntFrom(queryConfig.heartbeat ?? 0) ?? false,
  })
}

export default connect(null, mapDispatchToProps)(Page)
