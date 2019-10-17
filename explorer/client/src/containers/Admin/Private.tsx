import React from 'react'
import { RouteComponentProps, Redirect } from '@reach/router'
import { connect } from 'react-redux'
import { State } from '../../reducers'

interface Props extends RouteComponentProps {
  authenticated: boolean
}

const Private: React.FC<Props> = ({ authenticated }) => {
  return authenticated ? <></> : <Redirect to="/admin/signin" noThrow />
}

const mapStateToProps = (state: State) => {
  return {
    authenticated: state.adminAuth.allowed,
  }
}

const ConnectedPrivate = connect(mapStateToProps)(Private)

export default ConnectedPrivate
