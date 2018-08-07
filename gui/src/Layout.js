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
import PrivateRoute from './PrivateRoute'
import Logo from 'components/Logo'
import Loading from 'components/Loading'
import universal from 'react-universal-component'
import { Link, Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { submitSignOut } from 'actions'
import { isFetchingSelector } from 'selectors'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.
const uniOpts = {loading: Loading}
const Bridges = universal(import('./containers/Bridges'), uniOpts)
const BridgeSpec = universal(import('./containers/BridgeSpec'), uniOpts)
const Configuration = universal(import('./containers/Configuration'), uniOpts)
const About = universal(import('./containers/About'), uniOpts)
const Jobs = universal(import('./containers/Jobs'), uniOpts)
const JobSpec = universal(import('./containers/JobSpec'), uniOpts)
const JobSpecRuns = universal(import('./containers/JobSpecRuns'), uniOpts)
const JobSpecRun = universal(import('./containers/JobSpecRun'), uniOpts)
const SignIn = universal(import('./containers/SignIn'), uniOpts)

const appBarHeight = 70
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
    const {classes, errors, isFetching} = this.props
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
            <ListItem button component={Link} to='/create' className={classes.menuitem}>
              <ListItemText primary='Create' />
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
                    <Logo width={39} height={45} spin={isFetching} />
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
                  <PrivateRoute exact path='/about' component={About} />
                  <PrivateRoute exact path='/config' component={Configuration} />
                  <PrivateRoute exact path='/create' component={Create} />
                  <PrivateRoute exact path='/bridges' component={Bridges} />
                  <PrivateRoute exact path='/bridges/:bridgeName' component={BridgeSpec} />
                  <PrivateRoute exact path='/' component={Jobs} />
                  <PrivateRoute exact path='/jobs/page/:jobPage' component={Jobs} />
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
  authenticated: state.authentication.allowed,
  errors: state.errors.messages,
  isFetching: isFetchingSelector(state)
})

const mapDispatchToProps = dispatch => bindActionCreators({submitSignOut}, dispatch)

export const ConnectedLayout = connect(mapStateToProps, mapDispatchToProps)(Layout)

export default hot(module)(withStyles(styles)(ConnectedLayout))
