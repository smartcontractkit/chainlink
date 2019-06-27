import React from 'react'
import PropTypes from 'prop-types'
import Logo from '@chainlink/styleguide/components/Logo'
import src from '../../images/no-activity-icon.svg'

const NoContent = props => {
  return <Logo src={src} alt="No Content" {...props} />
}

NoContent.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number
}

export default NoContent
