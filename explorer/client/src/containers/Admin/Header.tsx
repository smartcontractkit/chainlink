import React from 'react'
import { bindActionCreators, Dispatch } from 'redux'
import { connect } from 'react-redux'
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
import { signOut } from '../../actions/adminAuth'
import { State } from '../../reducers'

const styles = (theme: Theme) =>
  createStyles({
    avatar: {
      float: 'right',
    },
    logo: {
      marginRight: theme.spacing.unit * 2,
      width: 200,
    },
  })

/* eslint-disable-next-line @typescript-eslint/no-empty-interface */
interface OwnProps {}

interface StateProps {
  authenticated: boolean
  errors: string[]
}

interface DispatchProps {
  signOut: () => void
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

interface Props extends RouteComponentProps, WithStyles<typeof styles> {
  onHeaderResize: (width: number, height: number) => void
}

export const AdminHeader = ({ classes, onHeaderResize, signOut }: Props) => {
  return (
    <Header onResize={onHeaderResize}>
      <Grid container>
        <Grid item xs={6}>
          <AdminLogo className={classes.logo} width={200} />
        </Grid>
        <Grid item xs={6}>
          <AvatarMenu className={classes.avatar}>
            <AvatarMenuItem text="Sign Out" onClick={signOut} />
          </AvatarMenu>
        </Grid>
      </Grid>
    </Header>
  )
}
function mapStateToProps(state: State): StateProps {
  return {
    authenticated: state.adminAuth.allowed,
    errors: state.notifications.errors,
  }
}

function mapDispatchToProps(dispatch: Dispatch): DispatchProps {
  return bindActionCreators({ signOut }, dispatch)
}

const ConnectedAdminHeader = connect(
  mapStateToProps,
  mapDispatchToProps,
)(AdminHeader)

export default withStyles(styles)(ConnectedAdminHeader)
