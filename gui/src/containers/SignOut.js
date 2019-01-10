import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import { receiveSignoutSuccess } from 'actions'
import { connect } from 'react-redux'

export const SignOut = props => {
  useEffect(() => { props.receiveSignoutSuccess() }, [])
  return <React.Fragment />
}

SignOut.propTypes = {
  receiveSignoutSuccess: PropTypes.func.isRequired
}

export default connect(
  null,
  { receiveSignoutSuccess }
)(SignOut)
