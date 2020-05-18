import React from 'react'
import { connect, MapStateToProps } from 'react-redux'
import { RouteComponentProps } from '@reach/router'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import { AdminLogo } from '../../components/Logos/Admin'
import Header from '../../components/Header'
import AvatarMenu from '../../components/AvatarMenu'
import AvatarMenuItem from '../../components/AvatarMenuItem'
import Link from '../../components/Link'
import { AppState } from '../../reducers'

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    avatar: {
      float: 'right',
    },
    logo: {
      marginRight: spacing.unit * 2,
      width: 200,
    },
    link: {
      color: palette.common.white,
      textDecoration: 'none',
    },
    menu: {
      display: 'inline-block',
      listStyle: 'none',
      marginRight: spacing.unit * 2,
      float: 'right',
    },
    menuItem: {
      display: 'inline',
      marginRight: spacing.unit * 2,
      '&::after': {
        content: "'|'",
        color: palette.grey[400],
        marginLeft: spacing.unit * 2,
      },
    },
  })

interface OwnProps {
  onHeaderResize: React.ComponentPropsWithoutRef<typeof Header>['onResize']
}

interface StateProps {
  authenticated: boolean
}

interface Props
  extends RouteComponentProps,
    StateProps,
    OwnProps,
    WithStyles<typeof styles> {}

export const AdminHeader: React.FC<Props> = ({ classes, onHeaderResize }) => {
  return (
    <Header onResize={onHeaderResize}>
      <Grid container>
        <Grid item xs={6}>
          <AdminLogo className={classes.logo} width={200} />
        </Grid>
        <Grid item xs={6}>
          <Grid container justify="flex-end">
            <Grid item xs={12}>
              <AvatarMenu className={classes.avatar}>
                <AvatarMenuItem>
                  <a href="/admin/signout" className={classes.link}>
                    Sign Out
                  </a>
                </AvatarMenuItem>
              </AvatarMenu>

              <ul className={classes.menu}>
                <li className={classes.menuItem}>
                  <Link to="/admin/operators">Operators</Link>
                </li>
                <li className={classes.menuItem}>
                  <Link to="/admin/heads">Heads</Link>
                </li>
              </ul>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Header>
  )
}

const mapStateToProps: MapStateToProps<
  StateProps,
  OwnProps,
  AppState
> = state => {
  return {
    authenticated: state.adminAuth.allowed,
  }
}

const ConnectedAdminHeader = connect(mapStateToProps)(AdminHeader)

export default withStyles(styles)(ConnectedAdminHeader)
