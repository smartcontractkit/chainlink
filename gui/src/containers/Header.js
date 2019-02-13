import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
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
import MainLogo from 'components/Logos/Main'
import AvatarMenu from 'components/AvatarMenu'
import { submitSignOut } from 'actions'
import fetchCountSelector from 'selectors/fetchCount'
import ReactResizeDetector from 'react-resize-detector'

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
      paddingTop: 0,
      paddingBottom: 0
    },
    horizontalNavItem: {
      display: 'inline'
    },
    horizontalNavLink: {
      color: theme.palette.secondary.main,
      paddingTop: theme.spacing.unit * 3,
      paddingBottom: theme.spacing.unit * 3,
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

const SHARED_NAV_ITEMS = [
  ['/jobs', 'Jobs'],
  ['/bridges', 'Bridges'],
  ['/config', 'Configuration']
]

class Header extends React.Component {
  constructor (props) {
    super(props)
    this.state = { drawerOpen: false }
  }

  setDrawerOpen (isOpen) {
    this.setState({
      drawerOpen: isOpen
    })
  }

  render () {
    const toggleDrawer = () => this.setDrawerOpen(!this.state.drawerOpen)
    const signOut = () => this.props.submitSignOut()
    const { classes, fetchCount } = this.props

    const drawer = (<Drawer
      anchor='right'
      open={this.state.drawerOpen}
      classes={{
        paper: classes.drawerPaper
      }}
      onClose={toggleDrawer}
    >
      <div
        tabIndex={0}
        role='button'
        onClick={toggleDrawer}
      >
        <List className={classes.drawerList}>
          {SHARED_NAV_ITEMS.map(([to, text]) => (
            <ListItem key={to} button component={Link} to={to} className={classes.menuitem}>
              <ListItemText primary={text} />
            </ListItem>
          ))}
          {this.props.authenticated &&
            <ListItem button onClick={signOut} className={classes.menuitem}>
              <ListItemText primary='Sign Out' />
            </ListItem>
          }
        </List>
      </div>
    </Drawer>)

    const nav = (
      <Typography variant='body1' component='div'>
        <List className={classes.horizontalNav}>
          {SHARED_NAV_ITEMS.map(([to, text]) => (
            <ListItem key={to} className={classes.horizontalNavItem}>
              <Link to={to} className={classes.horizontalNavLink}>{text}</Link>
            </ListItem>
          ))}
          {this.props.authenticated &&
            <ListItem className={classes.horizontalNavItem}>
              <AvatarMenu />
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
        <ReactResizeDetector
          refreshMode='debounce'
          refreshRate={200}
          onResize={this.props.onResize}
          handleHeight
        >
          <LoadingBar fetchCount={fetchCount} />

          <Toolbar className={classes.toolbar}>
            <Grid container alignItems='center'>
              <Grid item xs={11} sm={6} md={4}>
                <Link to='/'>
                  <MainLogo width={200} />
                </Link>
              </Grid>
              <Grid item xs={1} sm={6} md={8}>
                <Grid container justify='flex-end'>
                  <Grid item>
                    <Hidden mdUp>
                      <IconButton aria-label='open drawer' onClick={toggleDrawer}>
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
  { submitSignOut },
  dispatch
)

export const ConnectedHeader = connect(mapStateToProps, mapDispatchToProps)(Header)

export default withStyles(styles)(ConnectedHeader)
