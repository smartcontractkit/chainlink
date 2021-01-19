import { Logo } from 'components/Logo'
import PropTypes from 'prop-types'
import React from 'react'
import src from '../../images/chainlink-operator-logo.svg'

const Main = (props) => {
  return <Logo src={src} alt="Chainlink Operator" {...props} />
}

Main.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number,
}

export default Main
