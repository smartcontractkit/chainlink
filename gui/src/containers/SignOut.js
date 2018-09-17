import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { receiveSignoutSuccess } from 'actions'
import { connect } from 'react-redux'

export class SignOut extends Component {
  componentWillMount () {
    this.props.receiveSignoutSuccess()
  }

  render () {
    return <React.Fragment />
  }
}

SignOut.propTypes = {
  receiveSignoutSuccess: PropTypes.func.isRequired
}

export default connect(
  null,
  {receiveSignoutSuccess}
)(SignOut)
