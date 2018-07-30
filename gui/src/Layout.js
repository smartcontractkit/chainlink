import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Routes from 'react-static-routes'
import AppBar from '@material-ui/core/AppBar'
import CssBaseline from '@material-ui/core/CssBaseline'
import Drawer from '@material-ui/core/Drawer'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import ListItemText from '@material-ui/core/ListItemText'
import Grid from '@material-ui/core/Grid'
import IconButton from '@material-ui/core/IconButton'
import MenuIcon from '@material-ui/icons/Menu'
import Flash from 'components/Flash'
import universal from 'react-universal-component'
import PrivateRoute from './PrivateRoute'
import { Link, Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { submitSignOut } from 'actions'

import logoImg from './logo.svg'

// Use universal-react-component for code-splitting non-static routes
const Bridges = universal(import('./containers/Bridges'))
const BridgeSpec = universal(import('./containers/BridgeSpec'))
const Configuration = universal(import('./containers/Configuration'))
const Jobs = universal(import('./containers/Jobs'))
const JobSpec = universal(import('./containers/JobSpec'))
const JobSpecRuns = universal(import('./containers/JobSpecRuns'))
const JobSpecRun = universal(import('./containers/JobSpecRun'))
const SignIn = universal(import('./containers/SignIn'))

const appBarHeight = 64
const drawerWidth = 240

// Custom styles
const styles = theme => {
  return {
    appBar: {
      backgroundColor: theme.palette.background.appBar,
      paddingLeft: theme.spacing.unit * 5,
      paddingRight: theme.spacing.unit * 5,
      zIndex: theme.zIndex.modal + 1
    },
    appBarContent: {
      height: appBarHeight
    },
    content: {
      margin: theme.spacing.unit * 5,
      marginTop: 0
    },
    menuButton: {
      color: theme.palette.common.white
    },
    menuitem: {
      padding: theme.spacing.unit * 3,
      display: 'block'
    },
    drawerPaper: {
      backgroundColor: theme.palette.common.white,
      width: drawerWidth
    },
    drawerList: {
      padding: 0
    },
    toolbar: {
      minHeight: appBarHeight
    },
    flash: {
      textAlign: 'center'
    }
  }
}

class Layout extends Component {
  constructor (props) {
    super(props)
    this.state = {drawerOpen: false}
    this.toggleDrawer = this.toggleDrawer.bind(this)
    this.signOut = this.signOut.bind(this)
  }

  toggleDrawer () {
    this.setState({drawerOpen: !this.state.drawerOpen})
  }

  signOut () {
    this.props.submitSignOut()
  }

  render () {
    const {classes, errors} = this.props
    const {drawerOpen} = this.state

    const drawer = (
      <Drawer
        anchor='right'
        open={drawerOpen}
        classes={{
          paper: classes.drawerPaper
        }}
        onClose={this.toggleDrawer}
      >
        <div className={classes.toolbar} />
        <div
          tabIndex={0}
          role='button'
          onClick={this.toggleDrawer}
        >
          <List className={classes.drawerList}>
            <ListItem button component={Link} to='/' className={classes.menuitem}>
              <ListItemText primary='Jobs' />
            </ListItem>
            <ListItem button component={Link} to='/bridges' className={classes.menuitem}>
              <ListItemText primary='Bridges' />
            </ListItem>
            <ListItem button component={Link} to='/config' className={classes.menuitem}>
              <ListItemText primary='Configuration' />
            </ListItem>
            <ListItem button component={Link} to='/about' className={classes.menuitem}>
              <ListItemText primary='About' />
            </ListItem>
            { this.props.authenticated &&
            <ListItem button onClick={this.signOut} className={classes.menuitem}>
              <ListItemText primary='Sign Out' />
            </ListItem>
            }
          </List>
        </div>
      </Drawer>
    )

    return (
      <Router>
        <Grid container>
          <CssBaseline />
          <Grid item xs={12}>
            <AppBar
              className={classes.appBar}
              color='default'
              position='absolute'
            >
              <Grid container alignItems='center' className={classes.appBarContent}>
                <Grid item xs={9}>
                  <Link to='/'>
                    <img src={logoImg} alt='Chainlink' width={121} height={44} />
                  </Link>
                </Grid>
                <Grid item xs={3}>
                  <div align='right'>
                    <IconButton
                      aria-label='open drawer'
                      onClick={this.toggleDrawer}
                      className={classes.menuButton}
                    >
                      <MenuIcon />
                    </IconButton>
                  </div>
                </Grid>
              </Grid>
            </AppBar>

            <div>
              <div className={classes.toolbar} />

              {
                errors.length > 0 &&
                <Flash error className={classes.flash}>
                  {errors.map((msg, i) => <p key={i}>{msg}</p>)}
                </Flash>
              }

              <div className={classes.content}>
                <Switch>
                  <Route exact path='/signin' component={SignIn} />
                  <PrivateRoute exact path='/job_specs/:jobSpecId' component={JobSpec} />
                  <PrivateRoute exact path='/job_specs/:jobSpecId/runs' component={JobSpecRuns} />
                  <PrivateRoute exact path='/job_specs/:jobSpecId/runs/page/:jobRunsPage' component={JobSpecRuns} />
                  <PrivateRoute exact path='/job_specs/:jobSpecId/runs/id/:jobRunId' component={JobSpecRun} />
                  <PrivateRoute exact path='/config' component={Configuration} />
                  <PrivateRoute exact path='/bridges' component={Bridges} />
                  <PrivateRoute exact path='/bridges/:bridgeName' component={BridgeSpec} />
                  <PrivateRoute exact path='/' component={Jobs} />
                  <Routes />
                </Switch>
              </div>
            </div>

            {drawer}
          </Grid>
        </Grid>
      </Router>
    )
  }
}

Layout.propTypes = {
  errors: PropTypes.array
}

Layout.defaultProps = {
  errors: []
}

const mapStateToProps = state => ({
  authenticated: state.session.authenticated,
  errors: state.errors.messages
})

const mapDispatchToProps = dispatch => bindActionCreators({submitSignOut}, dispatch)

export const ConnectedLayout = connect(mapStateToProps, mapDispatchToProps)(Layout)

export default hot(module)(withStyles(styles)(ConnectedLayout))
