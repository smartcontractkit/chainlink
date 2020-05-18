import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom'
import { connect } from 'react-redux'
import { aggregatorOperations } from 'state/ducks/aggregator'
import { AggregatorVis } from 'components/aggregatorVis'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { Header } from 'components/header'
import { OracleTable } from 'components/oracleTable'
import { FeedConfig } from 'config'
import { parseQuery, uIntFrom } from 'utils'

interface OwnProps {}

interface DispatchProps {
  initContract: DispatchBinding<typeof aggregatorOperations.initContract>
}

interface Props extends OwnProps, DispatchProps {}

const Page: React.FC<Props> = ({ initContract }) => {
  const location = useLocation()
  const [config] = useState<FeedConfig>(
    parseConfig(parseQuery(location.search)),
  )

  useEffect(() => {
    try {
      initContract(config)
    } catch (error) {
      console.error('Could not initiate contract:', error)
    }
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

/**
 * Hydrate a feed config into its internal representation
 *
 * @param config The config in map format
 */
function parseConfig(config: Record<string, string>): FeedConfig {
  return {
    ...((config as unknown) as FeedConfig),
    networkId: uIntFrom(config.networkId ?? 0),
    contractVersion: 2,
    decimalPlaces: uIntFrom(config.decimalPlaces ?? 0),
    heartbeat: uIntFrom(config.heartbeat ?? 0) ?? false,
  }
}

export default connect(null, mapDispatchToProps)(Page)
