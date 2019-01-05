import React from 'react'
import { Route, Redirect } from 'react-router-dom'
import { connect } from 'react-redux'

export class PrivateRoute extends Route {
  render = props => (
    this.props.authenticated ? super.render(props) : <Redirect to='/signin' />
  )
}

const mapStateToProps = state => {
  return {
    authenticated: state.authentication.allowed
  }
}

export default connect(mapStateToProps)(PrivateRoute)
