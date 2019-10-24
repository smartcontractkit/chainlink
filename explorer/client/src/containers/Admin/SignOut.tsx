import React, { useEffect } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { Redirect, RouteComponentProps } from '@reach/router'
import { signOut } from '../../actions/adminAuth'
import { State as AppState } from '../../reducers'

interface OwnProps {}

interface StateProps {
  authenticated: boolean
}

interface DispatchProps {
  signOut: () => void
}

interface Props
  extends RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

export const SignOut: React.FC<Props> = ({ authenticated, signOut }) => {
  useEffect(() => {
    signOut()
  }, [signOut])

  return authenticated ? <></> : <Redirect to="/admin/signin" noThrow />
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

const mapDispatchToProps: MapDispatchToProps<DispatchProps, OwnProps> = {
  signOut,
}

export const ConnectedSignOut = connect(
  mapStateToProps,
  mapDispatchToProps,
)(SignOut)

export default ConnectedSignOut
