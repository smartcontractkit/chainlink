import React from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { Redirect, RouteComponentProps } from '@reach/router'
import { SignIn as SignInForm } from '../../components/Forms/SignIn'
import { signIn } from '../../actions/adminAuth'
import { AppState } from '../../reducers'

interface OwnProps {}

interface StateProps {
  authenticated: boolean
  errors: string[]
}

interface DispatchProps {
  signIn: (...args: Parameters<typeof signIn>) => void
}

interface Props
  extends RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const SignIn: React.FC<Props> = ({ authenticated, errors, signIn }) => {
  return authenticated ? (
    <Redirect to="/admin" noThrow />
  ) : (
    <SignInForm onSubmit={signIn} errors={errors} />
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

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  signIn,
}

export const ConnectedSignIn = connect(
  mapStateToProps,
  mapDispatchToProps,
)(SignIn)

export default ConnectedSignIn
