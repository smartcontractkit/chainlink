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
  WithStyles,
} from '@material-ui/core/styles'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import MenuIcon from '@material-ui/icons/Menu'
import classNames from 'classnames'
import React from 'react'
import { connect } from 'react-redux'
import ReactResizeDetector from 'react-resize-detector'
import { useLocation } from 'react-router-dom'
import { bindActionCreators, Dispatch } from 'redux'
import { submitSignOut } from 'actionCreators'
import { AccountMenu } from '../components/AccountMenu'
import { SettingsMenu } from '../components/SettingsMenu'
import BaseLink from '../components/BaseLink'
import LoadingBar from '../components/LoadingBar'
import MainLogo from '../components/Logos/Main'
import fetchCountSelector from '../selectors/fetchCount'
import { Feature, useFeatureFlag } from 'src/hooks/useFeatureFlag'

const NAV_ITEMS = [
  ['/jobs', 'Jobs'],
  ['/runs', 'Runs'],
  ['/chains', 'Chains'],
  ['/nodes', 'Nodes'],
  ['/bridges', 'Bridges'],
  ['/transactions', 'Transactions'],
]

const DRAWER_NAV_ITEMS = [
  ...NAV_ITEMS,
  ['/keys', 'Keys'],
  ['/config', 'Configuration'],
]

const drawerWidth = 240

const drawerStyles = ({ palette, spacing }: Theme) =>
  createStyles({
    menuitem: {
      padding: spacing.unit * 3,
      display: 'block',
    },
    drawerPaper: {
      backgroundColor: palette.common.white,
      paddingTop: spacing.unit * 7,
      width: drawerWidth,
    },
    drawerList: {
      padding: 0,
    },
  })

interface DrawerProps extends WithStyles<typeof drawerStyles> {
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
    submitSignOut,
  }: DrawerProps) => {
    return (
      <MuiDrawer
        anchor="right"
        open={drawerOpen}
        classes={{
          paper: classes.drawerPaper,
        }}
        onClose={toggleDrawer}
      >
        <div tabIndex={0} role="button" onClick={toggleDrawer}>
          <List className={classes.drawerList}>
            {DRAWER_NAV_ITEMS.map(([href, text]) => (
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
  },
)

const navStyles = ({ palette, spacing }: Theme) =>
  createStyles({
    horizontalNav: {
      paddingTop: 0,
      paddingBottom: 0,
    },
    horizontalNavItem: {
      display: 'inline',
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
        borderBottomColor: palette.primary.main,
      },
    },
    activeNavLink: {
      color: palette.primary.main,
      borderBottomColor: palette.primary.main,
    },
  })

interface NavProps extends WithStyles<typeof navStyles> {
  authenticated: boolean
}

const Nav = withStyles(navStyles)(({ authenticated, classes }: NavProps) => {
  const { pathname } = useLocation()
  const isFeedsManagerFeatureEnabled = useFeatureFlag(Feature.FeedsManager)

  return (
    <Typography variant="body1" component="div">
      <List className={classes.horizontalNav}>
        {NAV_ITEMS.map(([navItemPath, text]) => (
          <ListItem key={navItemPath} className={classes.horizontalNavItem}>
            <BaseLink
              key={navItemPath}
              href={navItemPath}
              className={classNames(
                classes.horizontalNavLink,
                pathname.startsWith(navItemPath) && classes.activeNavLink,
              )}
            >
              {text}
            </BaseLink>
          </ListItem>
        ))}
        {/* Feeds Manager link hidden behind a feature flag. This is temporary until we
        enable this for everyone */}
        {isFeedsManagerFeatureEnabled && (
          <ListItem className={classes.horizontalNavItem}>
            <BaseLink
              href={'/feeds_manager'}
              className={classNames(
                classes.horizontalNavLink,
                pathname.includes('/feeds_manager') && classes.activeNavLink,
              )}
            >
              Feeds Manager
            </BaseLink>
          </ListItem>
        )}
        {authenticated && (
          <>
            <ListItem className={classes.horizontalNavItem}>
              <SettingsMenu />
              <AccountMenu />
            </ListItem>
          </>
        )}
      </List>
    </Typography>
  )
})

const styles = ({ palette, spacing, zIndex }: Theme) =>
  createStyles({
    appBar: {
      backgroundColor: palette.common.white,
      zIndex: zIndex.modal - 1,
    },
    toolbar: {
      paddingLeft: spacing.unit * 5,
      paddingRight: spacing.unit * 5,
    },
  })

interface Props extends WithStyles<typeof styles> {
  fetchCount: number
  authenticated: boolean
  drawerContainer: HTMLElement | null
  submitSignOut: () => void
  onResize: (width: number, height: number) => void
}

interface State {
  drawerOpen: boolean
}

class Header extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = {
      drawerOpen: false,
    }
    this.toggleDrawer = this.toggleDrawer.bind(this)
  }

  toggleDrawer() {
    this.setState({
      drawerOpen: !this.state.drawerOpen,
    })
  }

  render() {
    const {
      authenticated,
      classes,
      fetchCount,
      drawerContainer,
      onResize,
      submitSignOut,
    } = this.props

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
                        onClick={this.toggleDrawer}
                      >
                        <MenuIcon />
                      </IconButton>
                    </Hidden>
                    <Hidden smDown>
                      <Nav authenticated={authenticated} />
                    </Hidden>
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </Toolbar>
        </ReactResizeDetector>
        <Portal container={drawerContainer}>
          <Drawer
            toggleDrawer={this.toggleDrawer}
            drawerOpen={this.state.drawerOpen}
            authenticated={authenticated}
            submitSignOut={submitSignOut}
          />
        </Portal>
      </AppBar>
    )
  }
}

const mapStateToProps = (state: any) => ({
  authenticated: state.authentication.allowed,
  fetchCount: fetchCountSelector(state),
  url: state.notifications.currentUrl,
})

const mapDispatchToProps = (dispatch: Dispatch) =>
  bindActionCreators({ submitSignOut }, dispatch)

const ConnectedHeader = connect(mapStateToProps, mapDispatchToProps)(Header)

export default withStyles(styles)(ConnectedHeader)
