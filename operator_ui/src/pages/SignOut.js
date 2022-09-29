import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import { receiveSignoutSuccess } from 'actionCreators'
import { connect } from 'react-redux'

export const SignOut = ({ receiveSignoutSuccess }) => {
  useEffect(() => {
    receiveSignoutSuccess()
  }, [receiveSignoutSuccess])
  return <React.Fragment />
}

SignOut.propTypes = {
  receiveSignoutSuccess: PropTypes.func.isRequired,
}

export default connect(null, { receiveSignoutSuccess })(SignOut)
