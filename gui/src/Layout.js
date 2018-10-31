import React, { Component } from 'react'
import Routes from 'react-static-routes'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import CssBaseline from '@material-ui/core/CssBaseline'
import Drawer from '@material-ui/core/Drawer'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import ListItemText from '@material-ui/core/ListItemText'
import Grid from '@material-ui/core/Grid'
import IconButton from '@material-ui/core/IconButton'
import MenuIcon from '@material-ui/icons/Menu'
import Typography from '@material-ui/core/Typography'
import Hidden from '@material-ui/core/Hidden'
import PrivateRoute from './PrivateRoute'
import Logo from 'components/Logo'
import LoadingBar from 'components/LoadingBar'
import Loading from 'components/Loading'
import Notifications from 'containers/Notifications'
import universal from 'react-universal-component'
import { Redirect } from 'react-router'
import { Link, Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { submitSignOut } from 'actions'
import fetchCountSelector from 'selectors/fetchCount'

// Asynchronously load routes that are chunked via code-splitting
// 'import' as a function must take a string. It can't take a variable.
const uniOpts = {loading: Loading}
const DashboardsIndex = universal(import('./containers/Dashboards/Index'), uniOpts)
const JobsIndex = universal(import('./containers/Jobs/Index'), uniOpts)
const JobsShow = universal(import('./containers/Jobs/Show'), uniOpts)
const JobsNew = universal(import('./containers/Jobs/New'), uniOpts)
const BridgesIndex = universal(import('./containers/Bridges/Index'), uniOpts)
const BridgesNew = universal(import('./containers/Bridges/New'), uniOpts)
const BridgesShow = universal(import('./containers/Bridges/Show'), uniOpts)
const BridgesEdit = universal(import('./containers/Bridges/Edit'), uniOpts)
const JobRunsIndex = universal(import('./containers/JobRuns/Index'), uniOpts)
const JobRunsShow = universal(import('./containers/JobRuns/Show'), uniOpts)
const Configuration = universal(import('./containers/Configuration'), uniOpts)
const About = universal(import('./containers/About'), uniOpts)
const SignIn = universal(import('./containers/SignIn'), uniOpts)
const SignOut = universal(import('./containers/SignOut'), uniOpts)

const drawerWidth = 240

// Custom styles
const styles = theme => {
  return {
    appBar: {
      backgroundColor: theme.palette.common.white,
      zIndex: theme.zIndex.modal + 1
    },
    toolbar: {
      paddingLeft: theme.spacing.unit * 5,
      paddingRight: theme.spacing.unit * 5
    },
    main: {
      paddingTop: 97
    },
    content: {
      margin: theme.spacing.unit * 5,
      marginTop: 0
    },
    menuitem: {
      padding: theme.spacing.unit * 3,
      display: 'block'
    },
    horizontalNav: {
      paddingBottom: 0
    },
    horizontalNavItem: {
      display: 'inline'
    },
    horizontalNavLink: {
      color: theme.palette.text.primary,
      paddingTop: theme.spacing.unit * 4,
      paddingBottom: theme.spacing.unit * 4,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: theme.palette.common.white,
      '&:hover': {
        color: theme.palette.primary.main,
        borderBottomColor: theme.palette.primary.main
      }
    },
    drawerPaper: {
      backgroundColor: theme.palette.common.white,
      paddingTop: theme.spacing.unit * 7,
      width: drawerWidth
    },
    drawerList: {
      padding: 0
    }
  }
}

class Layout extends Component {
  state = {drawerOpen: false}

  toggleDrawer = () => {
    this.setState({drawerOpen: !this.state.drawerOpen})
  }

  signOut = () => {
    this.props.submitSignOut()
  }

  render () {
    const {classes, fetchCount, redirectTo} = this.props
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
        <div
          tabIndex={0}
          role='button'
          onClick={this.toggleDrawer}
        >
          <List className={classes.drawerList}>
            <ListItem button component={Link} to='/jobs' className={classes.menuitem}>
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
            {this.props.authenticated &&
              <ListItem button onClick={this.signOut} className={classes.menuitem}>
                <ListItemText primary='Sign Out' />
              </ListItem>
            }
          </List>
        </div>
      </Drawer>
    )

    const nav = (
      <Typography variant='body1' component='div'>
        <List className={classes.horizontalNav}>
          <ListItem className={classes.horizontalNavItem}>
            <Link to='/jobs' className={classes.horizontalNavLink}>Jobs</Link>
          </ListItem>
          <ListItem className={classes.horizontalNavItem}>
            <Link to='/bridges' className={classes.horizontalNavLink}>Bridges</Link>
          </ListItem>
          <ListItem className={classes.horizontalNavItem}>
            <Link to='/config' className={classes.horizontalNavLink}>Configuration</Link>
          </ListItem>
          <ListItem className={classes.horizontalNavItem}>
            <Link to='/about' className={classes.horizontalNavLink}>About</Link>
          </ListItem>
          {this.props.authenticated &&
            <ListItem className={classes.horizontalNavItem}>
              <Link to='/signout' className={classes.horizontalNavLink}>Sign Out</Link>
            </ListItem>
          }
        </List>
      </Typography>
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
              <LoadingBar fetchCount={fetchCount} />

              <Toolbar className={classes.toolbar}>
                <Grid container alignItems='center' className={classes.appBarContent}>
                  <Grid item xs={6} md={4}>
                    <Link to='/'>
                      <Logo width={40} height={50} />
                    </Link>
                  </Grid>
                  <Grid item xs={6} md={8}>
                    <Grid container justify='flex-end'>
                      <Grid item>
                        <Hidden mdUp>
                          <IconButton aria-label='open drawer' onClick={this.toggleDrawer}>
                            <MenuIcon />
                          </IconButton>
                        </Hidden>
                        <Hidden smDown>
                          {nav}
                        </Hidden>
                      </Grid>
                    </Grid>
                  </Grid>
                </Grid>
              </Toolbar>
            </AppBar>

            <main className={classes.main}>
              <Notifications />

              <div className={classes.content}>
                <Switch>
                  <Route exact path='/signin' component={SignIn} />
                  <PrivateRoute exact path='/signout' component={SignOut} />
                  {redirectTo && <Redirect to={redirectTo} />}
                  <PrivateRoute
                    exact
                    path='/'
                    render={props => <DashboardsIndex {...props} recentlyCreatedPageSize={4} />}
                  />
                  <PrivateRoute exact path='/jobs' component={JobsIndex} />
                  <PrivateRoute exact path='/jobs/page/:jobPage' component={JobsIndex} />
                  <PrivateRoute exact path='/jobs/new' component={JobsNew} />
                  <PrivateRoute exact path='/jobs/:jobSpecId' component={JobsShow} />
                  <PrivateRoute exact path='/jobs/:jobSpecId/runs' component={JobRunsIndex} />
                  <PrivateRoute exact path='/jobs/:jobSpecId/runs/page/:jobRunsPage' component={JobRunsIndex} />
                  <PrivateRoute exact path='/jobs/:jobSpecId/runs/id/:jobRunId' component={JobRunsShow} />
                  <PrivateRoute exact path='/bridges' component={BridgesIndex} />
                  <PrivateRoute exact path='/bridges/page/:bridgePage' component={BridgesIndex} />
                  <PrivateRoute exact path='/bridges/new' component={BridgesNew} />
                  <PrivateRoute exact path='/bridges/:bridgeId' component={BridgesShow} />
                  <PrivateRoute exact path='/bridges/:bridgeId/edit' component={BridgesEdit} />
                  <PrivateRoute exact path='/about' component={About} />
                  <PrivateRoute exact path='/config' component={Configuration} />
                  <Routes />
                </Switch>
              </div>
            </main>

            {drawer}
          </Grid>
        </Grid>
      </Router>
    )
  }
}

const mapStateToProps = state => ({
  authenticated: state.authentication.allowed,
  fetchCount: fetchCountSelector(state),
  redirectTo: state.redirect.to
})

const mapDispatchToProps = dispatch => bindActionCreators({submitSignOut}, dispatch)

export const ConnectedLayout = connect(mapStateToProps, mapDispatchToProps)(Layout)

export default hot(module)(withStyles(styles)(ConnectedLayout))
