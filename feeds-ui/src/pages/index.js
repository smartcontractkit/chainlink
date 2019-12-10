import React from 'react'
import { Route, BrowserRouter } from 'react-router-dom'
import { Header } from 'components/header'

import EthUsdPage from './EthUsdPage'
import Testnet from './Testnet'
import CreatePage from './Create'
import CustomPage from './Custom'
import RopstenPage from './Ropsten'
import MainnetPage from './Mainnet'

const AppRoutes = () => {
  return (
    <BrowserRouter>
      <Header />
      <Route exact path="/" component={EthUsdPage} />
      <Route exact path="/testnet" component={Testnet} />
      <Route exact path="/create" component={CreatePage} />
      <Route exact path="/custom" component={CustomPage} />
      <Route path="/ropsten/:address" component={RopstenPage} />
      <Route path="/mainnet/:address" component={MainnetPage} />
    </BrowserRouter>
  )
}

export default AppRoutes
