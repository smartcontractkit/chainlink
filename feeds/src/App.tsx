import React, { useEffect } from 'react'
import { Route, Router, Switch, RouteComponentProps } from 'react-router-dom'
import { createBrowserHistory } from 'history'
import ReactGA from 'react-ga'
import { FeedConfig } from 'feeds'
import * as pages from './pages'
import WithFeedConfig from './enhancers/WithFeedConfig'
import { Footer } from './components/footer'
import { Networks } from './utils'

const history = createBrowserHistory()

history.listen(location => {
  ReactGA.pageview(location.pathname + location.search)
})

const injectFeedConfig = (networkId?: Networks) => (
  props: RouteComponentProps<any>,
) => (
  <WithFeedConfig
    networkId={networkId}
    {...props}
    render={(config: FeedConfig) => {
      if (config.contractVersion === 3) {
        return <pages.FluxAggregator config={config} />
      } else {
        return <pages.Aggregator config={config} />
      }
    }}
  />
)

const App = () => {
  useEffect(() => {
    ReactGA.pageview(window.location.pathname + window.location.search)
  }, [])

  return (
    <Router history={history}>
      <Switch>
        <Route exact path="/" component={pages.Landing} />
        <Route exact path="/create" component={pages.Create} />
        <Route exact path="/custom" component={pages.Custom} />
        <Route
          path="/ropsten/:pair"
          component={injectFeedConfig(Networks.ROPSTEN)}
        />
        <Route path="/address/:address" component={injectFeedConfig()} />
        <Route path="/:pair" component={injectFeedConfig(Networks.MAINNET)} />
      </Switch>
      <Footer />
    </Router>
  )
}

export default App
