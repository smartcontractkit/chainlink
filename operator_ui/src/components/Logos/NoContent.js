import { Logo } from 'components/Logo'
import PropTypes from 'prop-types'
import React from 'react'
import src from '../../images/no-activity-icon.svg'

const NoContent = (props) => {
  return <Logo src={src} alt="No Content" {...props} />
}

NoContent.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number,
}

export default NoContent
