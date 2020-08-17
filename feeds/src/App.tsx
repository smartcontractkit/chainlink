import React, { useEffect } from 'react'
import { Route, Router, Switch } from 'react-router-dom'
import { createBrowserHistory } from 'history'
import ReactGA from 'react-ga'
import * as pages from './pages'
import { Footer } from './components/footer'
import { Config } from 'config'

const history = createBrowserHistory()

history.listen(location => {
  ReactGA.pageview(location.pathname + location.search)
})

const allowDevRoutes = Config.devHostnameWhitelist().includes(
  window.location.hostname,
)
const devRoutes = [
  <Route exact path="/create" key="create" component={pages.Create} />,
  <Route exact path="/custom" key="custom" component={pages.Custom} />,
]

const App = () => {
  useEffect(() => {
    ReactGA.pageview(window.location.pathname + window.location.search)
  }, [])

  return (
    <Router history={history}>
      <Switch>
        <Route exact path="/" component={pages.Landing} />
        {allowDevRoutes && devRoutes}
        <Route
          path="/address/:contractAddress"
          component={pages.AggregatorByAddress}
        />
        <Route path="/:network/:pair" component={pages.AggregatorByPair} />
        <Route path="/:pair" component={pages.AggregatorByPair} />
      </Switch>
      <Footer />
    </Router>
  )
}

export default App
