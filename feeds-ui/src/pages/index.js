import React from 'react'
import { Route, BrowserRouter, Switch } from 'react-router-dom'
import { Header } from 'components/header'

import EthUsdPage from './EthUsdPage'
import CreatePage from './Create'
import CustomPage from './Custom'
import RopstenPage from './Ropsten'
import MainnetPage from './Mainnet'

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Header />
      <Switch>
        <Route exact path="/" component={EthUsdPage} />
        <Route exact path="/create" component={CreatePage} />
        <Route exact path="/custom" component={CustomPage} />
        <Route path="/ropsten/:pair" component={RopstenPage} />
        <Route path="/:pair" component={MainnetPage} />
      </Switch>
    </BrowserRouter>
  )
}

export default AppRoutes
