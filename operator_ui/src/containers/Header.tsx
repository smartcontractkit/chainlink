import AppBar from '@material-ui/core/AppBar'
import MuiDrawer from '@material-ui/core/Drawer'
import Grid from '@material-ui/core/Grid'
import Hidden from '@material-ui/core/Hidden'
import IconButton from '@material-ui/core/IconButton'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import ListItemText from '@material-ui/core/ListItemText'
import Portal from '@material-ui/core/Portal'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import MenuIcon from '@material-ui/icons/Menu'
import classNames from 'classnames'
import React from 'react'
import { connect } from 'react-redux'
import ReactResizeDetector from 'react-resize-detector'
import { bindActionCreators, Dispatch } from 'redux'
import { useHooks, useState } from 'use-react-hooks'
import { submitSignOut } from '../actions'
import AvatarMenu from '../components/AvatarMenu'
import BaseLink from '../components/BaseLink'
import LoadingBar from '../components/LoadingBar'
import MainLogo from '../components/Logos/Main'
import fetchCountSelector from '../selectors/fetchCount'

const SHARED_NAV_ITEMS = [
  ['/jobs', 'Jobs'],
  ['/runs', 'Runs'],
  ['/bridges', 'Bridges'],
  ['/transactions', 'Transactions'],
  ['/config', 'Configuration']
]

const drawerStyles = ({ palette, spacing }: Theme) =>
  createStyles({
    menuitem: {
      padding: spacing.unit * 3,
      display: 'block'
    },
    drawerPaper: {
      backgroundColor: palette.common.white,
      paddingTop: spacing.unit * 7,
      width: drawerWidth
    },
    drawerList: {
      padding: 0
    }
  })

interface IDrawerProps extends WithStyles<typeof drawerStyles> {
  authenticated: boolean
  drawerOpen: boolean
  toggleDrawer: () => void
  submitSignOut: () => void
}

const Drawer = withStyles(drawerStyles)(
  ({
    drawerOpen,
    toggleDrawer,
    authenticated,
    classes,
    submitSignOut
  }: IDrawerProps) => {
    return (
      <MuiDrawer
        anchor="right"
        open={drawerOpen}
        classes={{
          paper: classes.drawerPaper
        }}
        onClose={toggleDrawer}
      >
        <div tabIndex={0} role="button" onClick={toggleDrawer}>
          <List className={classes.drawerList}>
            {SHARED_NAV_ITEMS.map(([href, text]) => (
              <ListItem
                key={href}
                button
                component={() => (
                  <BaseLink href={href}>
                    <ListItemText primary={text} />
                  </BaseLink>
                )}
                className={classes.menuitem}
              />
            ))}
            {authenticated && (
              <ListItem
                button
                onClick={submitSignOut}
                className={classes.menuitem}
              >
                <ListItemText primary="Sign Out" />
              </ListItem>
            )}
          </List>
        </div>
      </MuiDrawer>
    )
  }
)

const navStyles = ({ palette, spacing }: Theme) =>
  createStyles({
    horizontalNav: {
      paddingTop: 0,
      paddingBottom: 0
    },
    horizontalNavItem: {
      display: 'inline'
    },
    horizontalNavLink: {
      color: palette.secondary.main,
      paddingTop: spacing.unit * 3,
      paddingBottom: spacing.unit * 3,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: palette.common.white,
      '&:hover': {
        borderBottomColor: palette.primary.main
      }
    },
    activeNavLink: {
      color: palette.primary.main,
      borderBottomColor: palette.primary.main
    }
  })

const isNavActive = (current?: string, to?: string) =>
  `${to && to.toLowerCase()}` === current

interface INavProps extends WithStyles<typeof navStyles> {
  authenticated: boolean
  url?: string
}

const Nav = withStyles(navStyles)(
  ({ authenticated, url, classes }: INavProps) => {
    return (
      <Typography variant="body1" component="div">
        <List className={classes.horizontalNav}>
          {SHARED_NAV_ITEMS.map(([to, text]) => (
            <ListItem key={to} className={classes.horizontalNavItem}>
              <BaseLink
                href={to}
                className={classNames(
                  classes.horizontalNavLink,
                  isNavActive(to, url) && classes.activeNavLink
                )}
              >
                {text}
              </BaseLink>
            </ListItem>
          ))}
          {authenticated && (
            <ListItem className={classes.horizontalNavItem}>
              <AvatarMenu />
            </ListItem>
          )}
        </List>
      </Typography>
    )
  }
)

const drawerWidth = 240

const styles = ({ palette, spacing, zIndex }: Theme) =>
  createStyles({
    appBar: {
      backgroundColor: palette.common.white,
      zIndex: zIndex.modal - 1
    },
    toolbar: {
      paddingLeft: spacing.unit * 5,
      paddingRight: spacing.unit * 5
    }
  })

interface IProps extends WithStyles<typeof styles> {
  fetchCount: number
  authenticated: boolean
  drawerContainer: Element
  submitSignOut: () => void
  onResize: () => void
  url?: string
}

const Header = useHooks(
  ({
    authenticated,
    classes,
    fetchCount,
    url,
    drawerContainer,
    onResize,
    submitSignOut
  }: IProps) => {
    const [drawerOpen, setDrawerOpen] = useState(false)
    const toggleDrawer = () => setDrawerOpen(!drawerOpen)

    return (
      <AppBar className={classes.appBar} color="default" position="absolute">
        <ReactResizeDetector
          refreshMode="debounce"
          refreshRate={200}
          onResize={onResize}
          handleHeight
        >
          <LoadingBar fetchCount={fetchCount} />

          <Toolbar className={classes.toolbar}>
            <Grid container alignItems="center">
              <Grid item xs={11} sm={6} md={4}>
                <BaseLink href="/">
                  <MainLogo width={200} />
                </BaseLink>
              </Grid>
              <Grid item xs={1} sm={6} md={8}>
                <Grid container justify="flex-end">
                  <Grid item>
                    <Hidden mdUp>
                      <IconButton
                        aria-label="open drawer"
                        onClick={toggleDrawer}
                      >
                        <MenuIcon />
                      </IconButton>
                    </Hidden>
                    <Hidden smDown>
                      <Nav authenticated={authenticated} url={url} />
                    </Hidden>
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </Toolbar>
        </ReactResizeDetector>
        <Portal container={drawerContainer}>
          <Drawer
            toggleDrawer={toggleDrawer}
            drawerOpen={drawerOpen}
            authenticated={authenticated}
            submitSignOut={submitSignOut}
          />
        </Portal>
      </AppBar>
    )
  }
)

const mapStateToProps = (state: any) => ({
  authenticated: state.authentication.allowed,
  fetchCount: fetchCountSelector(state),
  url: state.notifications.currentUrl
})

const mapDispatchToProps = (dispatch: Dispatch) =>
  bindActionCreators({ submitSignOut }, dispatch)

export const ConnectedHeader = connect(
  mapStateToProps,
  mapDispatchToProps
)(Header)

export default withStyles(styles)(ConnectedHeader)
