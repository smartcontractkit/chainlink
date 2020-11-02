import React from 'react'
import { Route, Redirect } from 'react-router-dom'
import { connect } from 'react-redux'

export class PrivateRoute extends Route {
  constructor(...args) {
    super(...args)
    this.render = this.render.bind(this)
  }

  render(props) {
    return this.props.authenticated ? (
      super.render(props)
    ) : (
      <Redirect to="/signin" />
    )
  }
}

const mapStateToProps = (state) => {
  return {
    authenticated: state.authentication.allowed,
  }
}

export default connect(mapStateToProps)(PrivateRoute)
