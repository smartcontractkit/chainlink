import React from 'react'
import PropTypes from 'prop-types'
import { receiveSignoutSuccess } from 'actions'
import { connect } from 'react-redux'
import { useHooks, useEffect } from 'use-react-hooks'

export const SignOut = useHooks(props => {
  useEffect(() => {
    document.title = 'Sign Out'
    props.receiveSignoutSuccess()
  }, [])
  return <React.Fragment />
})

SignOut.propTypes = {
  receiveSignoutSuccess: PropTypes.func.isRequired,
}

export default connect(
  null,
  { receiveSignoutSuccess },
)(SignOut)
