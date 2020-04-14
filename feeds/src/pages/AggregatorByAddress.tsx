import React, { useEffect, useState } from 'react'
import { Redirect } from 'react-router-dom'
import { connect, MapStateToProps } from 'react-redux'
import { FeedConfig } from 'config'
import { Aggregator } from 'components/aggregator'
import { Header } from 'components/header'
import { AppState } from 'state'

interface OwnProps {
  match: {
    params: {
      contractAddress: string
    }
  }
}

interface StateProps {
  config?: FeedConfig
}

interface DispatchProps {}

interface Props extends OwnProps, StateProps, DispatchProps {}

const Page: React.FC<Props> = ({ config }) => {
  const [loaded, setLoaded] = useState<boolean>(false)
  let content

  useEffect(() => {
    setLoaded(true)
  }, [loaded, setLoaded])

  if (config) {
    content = <Aggregator config={config} />
  } else if (loaded) {
    content = <Redirect to="/" />
  } else {
    content = <>Loading Feed...</>
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

function selectFeedConfig(
  { feeds }: AppState,
  contractAddress: string,
): FeedConfig | undefined {
  return feeds.items[contractAddress]
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, AppState> = (
  state,
  ownProps,
) => {
  const contractAddress = ownProps.match.params.contractAddress
  const config = selectFeedConfig(state, contractAddress)

  return {
    config,
  }
}

export default connect(mapStateToProps)(Page)
