import React from 'react'
import { bindActionCreators, Dispatch } from 'redux'
import { connect } from 'react-redux'
import { Redirect, RouteComponentProps } from '@reach/router'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import { SignIn as SignInForm } from '../../components/Forms/SignIn'
import { signIn } from '../../actions/adminAuth'
import { State } from '../../reducers'

const styles = () => createStyles({})

/* eslint-disable-next-line @typescript-eslint/no-empty-interface */
interface OwnProps {}

interface StateProps {
  authenticated: boolean
  errors: string[]
}

interface DispatchProps {
  signIn: (username: string, password: string) => void
}

interface Props
  extends WithStyles<typeof styles>,
    RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const SignIn = ({ authenticated, errors, signIn }: Props) => {
  return authenticated ? (
    <Redirect to="/admin" noThrow />
  ) : (
    <SignInForm onSubmit={signIn} errors={errors} />
  )
}

function mapStateToProps(state: State): StateProps {
  return {
    authenticated: state.adminAuth.allowed,
    errors: state.notifications.errors,
  }
}

function mapDispatchToProps(dispatch: Dispatch): DispatchProps {
  return bindActionCreators({ signIn }, dispatch)
}

export const ConnectedSignIn = connect(
  mapStateToProps,
  mapDispatchToProps,
)(SignIn)

export default withStyles(styles)(ConnectedSignIn)
