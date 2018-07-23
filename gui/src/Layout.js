import React, { Component } from 'react'
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
import universal from 'react-universal-component'
import { Link, Router, Route, Switch } from 'react-static'
import { hot } from 'react-hot-loader'
import { withStyles } from '@material-ui/core/styles'
import logoImg from './logo.svg'

// Use universal-react-component for code-splitting non-static routes
const Bridges = universal(import('./containers/Bridges'))
const Configuration = universal(import('./containers/Configuration'))
const JobSpec = universal(import('./containers/JobSpec'))
const JobSpecRuns = universal(import('./containers/JobSpecRuns'))
const JobSpecRun = universal(import('./containers/JobSpecRun'))

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
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
      display: 'inline-block',
      width: 'inherit',
      textDecoration: 'none'
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
    }
  }
}

class Layout extends Component {
  constructor (props) {
    super(props)
    this.state = {drawerOpen: false}
    this.toggleDrawer = this.toggleDrawer.bind(this)
  }

  toggleDrawer () {
    this.setState({drawerOpen: !this.state.drawerOpen})
  }

  render () {
    const {classes} = this.props
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
            <ListItem button>
              <Link to='/' className={classes.menuitem}>
                <ListItemText primary='Jobs' />
              </Link>
            </ListItem>
            <ListItem button>
              <Link to='/bridges' className={classes.menuitem}>
                <ListItemText primary='Bridges' />
              </Link>
            </ListItem>
            <ListItem button>
              <Link to='/config' className={classes.menuitem}>
                <ListItemText primary='Configuration' />
              </Link>
            </ListItem>
            <ListItem button>
              <Link to='/about' className={classes.menuitem}>
                <ListItemText primary='About' />
              </Link>
            </ListItem>
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

            <div className={classes.content}>
              <div className={classes.toolbar} />
              <Switch>
                <Route exact path='/job_specs/:jobSpecId' component={JobSpec} />
                <Route exact path='/job_specs/:jobSpecId/runs' component={JobSpecRuns} />
                <Route exact path='/job_specs/:jobSpecId/runs/page/:jobRunsPage' component={JobSpecRuns} />
                <Route exact path='/job_specs/:jobSpecId/runs/id/:jobRunId' component={JobSpecRun} />
                 <Route exact path='/config' component={Configuration} />
                <Route exact path='/bridges' component={Bridges} />
                <Routes />
              </Switch>
            </div>

            {drawer}
          </Grid>
        </Grid>
      </Router>
    )
  }
}

const LayoutWithStyles = withStyles(styles)(Layout)

export default hot(module)(LayoutWithStyles)
