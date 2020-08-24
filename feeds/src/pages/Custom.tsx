import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect } from 'react'
import { RouteComponentProps } from 'react-router-dom'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { FeedConfig } from 'config'
import { Header } from 'components/header'
import { Aggregator, FluxAggregator } from 'components/aggregator'
import {
  aggregatorActions,
  aggregatorOperations,
} from '../state/ducks/aggregator'
import { useLocation } from 'react-router-dom'
import { parseQuery, uIntFrom } from 'utils'
import { AppState } from '../state/reducers'

interface OwnProps
  extends RouteComponentProps<{ pair: string; network?: string }> {}

interface StateProps {
  config: FeedConfig
}

interface DispatchProps {
  fetchOracleNodes: DispatchBinding<
    typeof aggregatorOperations.fetchOracleNodes
  >
  storeAggregatorConfig: DispatchBinding<
    typeof aggregatorActions.storeAggregatorConfig
  >
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const Page: React.FC<Props> = ({
  config,
  fetchOracleNodes,
  storeAggregatorConfig,
}) => {
  const location = useLocation()

  useEffect(() => {
    storeAggregatorConfig(parseConfig(parseQuery(location.search)))
  }, [storeAggregatorConfig, location.search])

  useEffect(() => {
    fetchOracleNodes()
  }, [fetchOracleNodes])

  let content

  if (config && config.contractVersion === 3) {
    content = <FluxAggregator config={config} />
  } else if (config) {
    content = <Aggregator config={config} />
  } else {
    content = <>There was an error loading the page. Refresh to try again.</>
  }

  return (
    <>
      <div className="page-container-full-width">
        <Header />
      </div>
      <div className="page-wrapper network-page">{content}</div>
    </>
  )
}

/**
 * Hydrate a feed config into its internal representation
 *
 * @param config The config in map format
 */
function parseConfig(config: Record<string, string>): FeedConfig {
  return {
    ...((config as unknown) as FeedConfig),
    networkId: uIntFrom(config.networkId ?? 1),
    contractVersion: uIntFrom(config.contractVersion ?? 2),
    decimalPlaces: uIntFrom(config.decimalPlaces ?? 4),
    heartbeat: uIntFrom(config.heartbeat ?? 0),
    historyDays: uIntFrom(config.historyDays ?? 1),
    formatDecimalPlaces: uIntFrom(config.formatDecimalPlaces ?? 0),
    threshold: uIntFrom(config.threshold ?? 0) ?? null,
    multiply: config.multiply ?? 100000000,
  }
}

const mapStateToProps: MapStateToProps<{}, OwnProps, AppState> = ({
  aggregator: { config },
}: AppState) => ({
  config,
})

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchOracleNodes: aggregatorOperations.fetchOracleNodes,
  storeAggregatorConfig: aggregatorActions.storeAggregatorConfig,
}

export default connect(mapStateToProps, mapDispatchToProps)(Page)
