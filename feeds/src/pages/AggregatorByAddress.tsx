import { DispatchBinding } from '@chainlink/ts-helpers'
import React, { useEffect } from 'react'
import { Redirect } from 'react-router-dom'
import { connect, MapStateToProps } from 'react-redux'
import { FeedConfig } from 'config'
import { Aggregator } from 'components/aggregator'
import { Header } from 'components/header'
import { AppState } from 'state'
import { aggregatorOperations } from '../state/ducks/aggregator'

interface OwnProps {
  match: {
    params: {
      contractAddress: string
    }
  }
}

interface StateProps {
  config?: FeedConfig
  loadingFeed: boolean
  errorFeed?: string
}

interface DispatchProps {
  fetchFeedByAddress: DispatchBinding<
    typeof aggregatorOperations.fetchFeedByAddress
  >
  fetchOracleNodes: DispatchBinding<
    typeof aggregatorOperations.fetchOracleNodes
  >
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const Page: React.FC<Props> = ({
  fetchFeedByAddress,
  fetchOracleNodes,
  match,
  loadingFeed,
  errorFeed,
  config,
}) => {
  const contractAddress = match.params.contractAddress
  useEffect(() => {
    fetchFeedByAddress(contractAddress)
  }, [fetchFeedByAddress, contractAddress])
  useEffect(() => {
    fetchOracleNodes()
  }, [fetchOracleNodes])

  let content
  if (config) {
    content = <Aggregator config={config} />
  } else if (loadingFeed) {
    content = <>Loading Feed...</>
  } else if (errorFeed && errorFeed === 'Not Found') {
    content = <Redirect to="/" />
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

const mapDispatchToProps = {
  fetchFeedByAddress: aggregatorOperations.fetchFeedByAddress,
  fetchOracleNodes: aggregatorOperations.fetchOracleNodes,
}

export default connect(mapStateToProps, mapDispatchToProps)(Page)
