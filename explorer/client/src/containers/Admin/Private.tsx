import React from 'react'
import { RouteComponentProps, Redirect } from '@reach/router'
import { connect, MapStateToProps } from 'react-redux'
import { AppState } from '../../reducers'

interface OwnProps {}

interface StateProps {
  authenticated: boolean
}

interface DispatchProps {}

interface Props
  extends RouteComponentProps,
    StateProps,
    DispatchProps,
    OwnProps {}

const Private: React.FC<Props> = ({ authenticated }) => {
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

const ConnectedPrivate = connect(mapStateToProps)(Private)

export default ConnectedPrivate
