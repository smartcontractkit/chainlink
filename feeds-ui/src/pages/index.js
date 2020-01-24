import React, { useEffect } from 'react'
import { Route, Router, Switch } from 'react-router-dom'
import Landing from './Landing'
import CreatePage from './Create'
import CustomPage from './Custom'
import DetailsPage from './Details'
import WithConfig from 'enhancers/withConfig'
import { ROPSTEN_ID, MAINNET_ID } from 'utils'
import { createBrowserHistory } from 'history'
import ReactGA from 'react-ga'

const history = createBrowserHistory()

history.listen(location => {
  ReactGA.pageview(location.pathname + location.search)
})

const configWrapper = networkId => props => (
  <WithConfig
    networkId={networkId}
    {...props}
    render={config => <DetailsPage config={config} />}
  />
)

const AppRoutes = () => {
  useEffect(() => {
    ReactGA.pageview(window.location.pathname + window.location.search)
  }, [])

  return (
    <Router history={history}>
      <Switch>
        <Route exact path="/" component={Landing} />
        <Route exact path="/create" component={CreatePage} />
        <Route exact path="/custom" component={CustomPage} />
        <Route path="/ropsten/:pair" component={configWrapper(ROPSTEN_ID)} />
        <Route path="/address/:address" component={configWrapper()} />
        <Route path="/:pair" component={configWrapper(MAINNET_ID)} />
      </Switch>
    </Router>
  )
}

export default AppRoutes
