import React, { useEffect } from 'react'
import { Route, Router, Switch } from 'react-router-dom'
import { createBrowserHistory } from 'history'
import ReactGA from 'react-ga'
import * as pages from './pages'
import { Footer } from './components/footer'

const history = createBrowserHistory()

history.listen(location => {
  ReactGA.pageview(location.pathname + location.search)
})

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
