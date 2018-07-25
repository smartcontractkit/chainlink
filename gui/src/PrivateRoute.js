import React from 'react'
import { Route, Redirect } from 'react-router'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

export class PrivateRoute extends Route {
  render = (props) => (
    this.props.authenticated === true
      ? super.render(props)
      : <Redirect to='/signin' />
  )
}

const mapStateToProps = state => {
  return {
    authenticated: state.session.authenticated
  }
}

const mapDispatchToProps = (dispatch) => {
  return bindActionCreators({}, dispatch)
}

export default connect(mapStateToProps, mapDispatchToProps)(PrivateRoute)
