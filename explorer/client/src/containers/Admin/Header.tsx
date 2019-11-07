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
  })

interface OwnProps {
  onHeaderResize: React.ComponentPropsWithoutRef<typeof Header>['onResize']
}

interface StateProps {
  authenticated: boolean
  errors: string[]
}

interface DispatchProps {}

interface Props
  extends RouteComponentProps,
    StateProps,
    DispatchProps,
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
          <AvatarMenu className={classes.avatar}>
            <AvatarMenuItem>
              <a href="/admin/signout" className={classes.link}>
                Sign Out
              </a>
            </AvatarMenuItem>
          </AvatarMenu>
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
    errors: state.notifications.errors,
  }
}

const ConnectedAdminHeader = connect(mapStateToProps)(AdminHeader)

export default withStyles(styles)(ConnectedAdminHeader)
