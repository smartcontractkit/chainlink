import React from 'react'
import PropTypes from 'prop-types'
import Logo from '@chainlink/styleguide/components/Logo'
import src from '../../images/icon-logo-blue.svg'

const Hexagon = props => {
  return <Logo src={src} alt="Chainlink Operator" {...props} />
}

Hexagon.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number
}

export default Hexagon
