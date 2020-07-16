import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect } from 'react'
import { Redirect, RouteComponentProps } from 'react-router-dom'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { FeedConfig } from 'config'
import { Header } from 'components/header'
import { Aggregator, FluxAggregator } from 'components/aggregator'
import { AppState } from 'state'
import { aggregatorOperations } from '../state/ducks/aggregator'

interface OwnProps
  extends RouteComponentProps<{ pair: string; network?: string }> {}

interface StateProps {
  config?: FeedConfig
  loadingFeed: boolean
  errorFeed?: string
}

interface DispatchProps {
  fetchFeedByPair: DispatchBinding<typeof aggregatorOperations.fetchFeedByPair>
  fetchOracleNodes: DispatchBinding<
    typeof aggregatorOperations.fetchOracleNodes
  >
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const Page: React.FC<Props> = ({
  fetchFeedByPair,
  fetchOracleNodes,
  match,
  loadingFeed,
  errorFeed,
  config,
}) => {
  const { pair, network } = match.params
  useEffect(() => {
    fetchFeedByPair(pair, network)
  }, [fetchFeedByPair, pair, network])
  useEffect(() => {
    fetchOracleNodes()
  }, [fetchOracleNodes])

  let content
  if (loadingFeed) {
    content = <>Loading Feed...</>
  } else if (errorFeed && errorFeed === 'Not Found') {
    content = <Redirect to="/" />
  } else if (config && config.contractVersion === 3) {
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

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    config: state.aggregator.config,
    loadingFeed: state.aggregator.loadingFeed,
    errorFeed: state.aggregator.errorFeed,
  }
}

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  fetchFeedByPair: aggregatorOperations.fetchFeedByPair,
  fetchOracleNodes: aggregatorOperations.fetchOracleNodes,
}

export default connect(mapStateToProps, mapDispatchToProps)(Page)
