import React, { PureComponent } from 'react'
import Routes from 'react-static-routes'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import AppBar from '@material-ui/core/AppBar'
import Typography from '@material-ui/core/Typography'
import universal from 'react-universal-component'
import createStore from 'connectors/redux'
import { Link, Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import { Provider } from 'react-redux'
import logoImg from './logo.svg'

// Use universal-react-component for code-splitting non-static routes
const Bridges = universal(import('./containers/Bridges'))
const Configuration = universal(import('./containers/Configuration'))
const JobSpec = universal(import('./containers/JobSpec'))
const JobSpecRuns = universal(import('./containers/JobSpecRuns'))
const JobSpecRun = universal(import('./containers/JobSpecRun'))

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
    },
    menu: {
      marginTop: theme.spacing.unit * 2,
      textDecoration: 'none'
    },
    menuitem: {
      color: theme.palette.common.white,
      paddingRight: theme.spacing.unit * 5,
      display: 'inline-block',
      textDecoration: 'none'
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
                <Grid container spacing={40}>
                  <Grid item xs={9}>
                    <Link to='/'>
                      <img src={logoImg} alt='Chainlink' width={121} height={44} />
                    </Link>
                  </Grid>
                  <Grid item xs={3}>
                    <div  align='right' className={classes.menu}>
                      <Link to='/bridges' className={classes.menuitem}>
                        <Typography color='inherit'>
                          Bridges
                        </Typography>
                      </Link>
                      <Link to='/config' className={classes.menuitem}>
                        <Typography color='inherit'>
                          Configuration
                        </Typography>
                      </Link>
                    </div>
                  </Grid>
                </Grid>
              </AppBar>

              <div className={classes.content}>
                <Switch>
                  <Route exact path='/job_specs/:jobSpecId' component={JobSpec} />
                  <Route exact path='/job_specs/:jobSpecId/runs' component={JobSpecRuns} />
                  <Route exact path='/job_specs/:jobSpecId/runs/:jobRunId' component={JobSpecRun} />
                  <Route exact path='/config' component={Configuration} />
                  <Route exact path='/bridges' component={Bridges} />
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
