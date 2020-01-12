import React from 'react'
import { connect } from 'react-redux'
import { Redirect } from 'react-router-dom'
import { hot } from 'react-hot-loader'
import { submitSignIn } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { getPersistUrl } from '../utils/storage'
import { SignInForm } from '@chainlink/styleguide'

export const SignIn = props => {
  document.title = 'Sign In'
  const { authenticated, errors, submitSignIn } = props
  return authenticated ? (
    <Redirect to={getPersistUrl() || '/'} />
  ) : (
    <SignInForm
      title="Operator"
      onSubmitOperator={submitSignIn}
      errors={errors}
    />
  )
}

const mapStateToProps = state => ({
  fetching: state.authentication.fetching,
  authenticated: state.authentication.allowed,
  errors: state.notifications.errors,
})

export const ConnectedSignIn = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ submitSignIn }),
)(SignIn)

export default hot(module)(ConnectedSignIn)
