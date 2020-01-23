import React from 'react'
import { Route, BrowserRouter, Switch } from 'react-router-dom'
import { Footer } from 'components/footer'
import Landing from './Landing'
import CreatePage from './Create'
import CustomPage from './Custom'
import DetailsPage from './Details'
import WithConfig from 'enhancers/withConfig'
import { ROPSTEN_ID, MAINNET_ID } from 'utils'

const configWrapper = networkId => props => (
  <WithConfig
    networkId={networkId}
    {...props}
    render={config => <DetailsPage config={config} />}
  />
)

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path="/" component={Landing} />
        <Route exact path="/create" component={CreatePage} />
        <Route exact path="/custom" component={CustomPage} />
        <Route path="/ropsten/:pair" component={configWrapper(ROPSTEN_ID)} />
        <Route path="/address/:address" component={configWrapper()} />
        <Route path="/:pair" component={configWrapper(MAINNET_ID)} />
      </Switch>
      <Footer />
    </BrowserRouter>
  )
}

export default AppRoutes
