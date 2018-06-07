import React, { PureComponent } from 'react'
import Routes from 'react-static-routes'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import AppBar from '@material-ui/core/AppBar'
import universal from 'react-universal-component'
import createStore from 'connectors/redux'
import { Link, Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import { Provider } from 'react-redux'
import logoImg from './logo.svg'

// Use universal-react-component for code-splitting non-static routes
const JobSpec = universal(import('./containers/JobSpec'))

// Custom styles
const styles = theme => {
  return {
    appBar: {
      backgroundColor: theme.palette.background.appBar,
      paddingTop: theme.spacing.unit * 3,
      paddingBottom: theme.spacing.unit * 3,
      paddingLeft: theme.spacing.unit * 5,
      paddingRight: theme.spacing.unit * 5
    },
    content: {
      margin: theme.spacing.unit * 5,
      marginTop: 0
    }
  }
}

class App extends PureComponent {
  // Remove the server-side injected CSS.
  componentDidMount () {
    const jssStyles = document.getElementById('jss-server-side')
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles)
    }
  }

  render () {
    const { classes } = this.props

    return (
      <Provider store={createStore()}>
        <Router>
          <Grid container>
            <CssBaseline />
            <Grid item xs={12}>
              <AppBar
                className={classes.appBar}
                elevation={0}
                color='default'
                position='static'
              >
                <Link to='/'>
                  <img src={logoImg} alt='Chainlink' width={121} height={44} />
                </Link>
              </AppBar>

              <div className={classes.content}>
                <Switch>
                  <Route path='/job_specs/:jobSpecId' component={JobSpec} />
                  <Routes />
                </Switch>
              </div>
            </Grid>
          </Grid>
        </Router>
      </Provider>
    )
  }
}

const AppWithStyles = withStyles(styles)(App)

export default hot(module)(AppWithStyles)
