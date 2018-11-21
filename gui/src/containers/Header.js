import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-static'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import ReactResizeDetector from 'react-resize-detector'
import { withStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Grid from '@material-ui/core/Grid'
import Hidden from '@material-ui/core/Hidden'
import Typography from '@material-ui/core/Typography'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import ListItemText from '@material-ui/core/ListItemText'
import Drawer from '@material-ui/core/Drawer'
import IconButton from '@material-ui/core/IconButton'
import MenuIcon from '@material-ui/icons/Menu'
import Portal from '@material-ui/core/Portal'
import LoadingBar from 'components/LoadingBar'
import Logo from 'components/Logo'
import { submitSignOut } from 'actions'
import fetchCountSelector from 'selectors/fetchCount'

const drawerWidth = 240

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
      color: theme.palette.secondary.main,
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

class Header extends Component {
  state = {drawerOpen: false}

  toggleDrawer = () => {
    this.setState({drawerOpen: !this.state.drawerOpen})
  }

  signOut = () => {
    this.props.submitSignOut()
  }

  render () {
    const {classes, fetchCount} = this.props
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
      <AppBar
        className={classes.appBar}
        color='default'
        position='absolute'
      >
        <ReactResizeDetector handleHeight onResize={this.props.onResize}>
          <LoadingBar fetchCount={fetchCount} />

          <Toolbar className={classes.toolbar}>
            <Grid container alignItems='center' className={classes.appBarContent}>
              <Grid item xs={11} sm={6} md={4}>
                <Link to='/'>
                  <Logo width={40} height={50} />
                </Link>
              </Grid>
              <Grid item xs={1} sm={6} md={8}>
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
        </ReactResizeDetector>
        <Portal container={this.props.drawerContainer}>
          {drawer}
        </Portal>
      </AppBar>
    )
  }
}

Header.propTypes = {
  onResize: PropTypes.func.isRequired,
  drawerContainer: PropTypes.object
}

const mapStateToProps = state => ({
  authenticated: state.authentication.allowed,
  fetchCount: fetchCountSelector(state)
})

const mapDispatchToProps = dispatch => bindActionCreators(
  {submitSignOut},
  dispatch
)

export const ConnectedHeader = connect(mapStateToProps, mapDispatchToProps)(Header)

export default withStyles(styles)(ConnectedHeader)
