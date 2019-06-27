import React from 'react'
import PropTypes from 'prop-types'
import Logo from '@chainlink/styleguide/components/Logo'
import src from '../../images/chainlink-operator-logo.svg'

const Main = props => {
  return <Logo src={src} alt="Chainlink Operator" {...props} />
}

Main.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number
}

export default Main
