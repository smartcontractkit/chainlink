import { Logo } from 'components/Logo'
import PropTypes from 'prop-types'
import React from 'react'
import src from '../../images/icon-logo-blue.svg'

const Hexagon = (props) => {
  return <Logo src={src} alt="Chainlink Operator" {...props} />
}

Hexagon.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number,
}

export default Hexagon
