import React, { PureComponent } from 'react'
import { Router } from 'react-static'
import { hot } from 'react-hot-loader'
import Routes from 'react-static-routes'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import AppBar from '@material-ui/core/AppBar'
import { withStyles } from '@material-ui/core/styles'
import logoImg from './logo.svg'

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
              <img src={logoImg} alt='Chainlink' width={121} height={44} />
            </AppBar>

            <div className={classes.content}>
              <Routes />
            </div>
          </Grid>
        </Grid>
      </Router>
    )
  }
}

const AppWithStyles = withStyles(styles)(App)

export default hot(module)(AppWithStyles)
