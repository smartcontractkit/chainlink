import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => {
  return {
    animate: {
      animation: 'spin 4s linear infinite'
    },
    '@keyframes spin': {
      '100%': {
        'transform': 'rotate(360deg)'
      }
    }
  }
}

const Image = ({src, width, height, spin, alt, classes}) => (
  <img
    src={src}
    width={width}
    height={height}
    className={spin ? classes.animate : ''}
    alt={alt}
  />
)

Image.propTypes = {
  src: PropTypes.string.isRequired,
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  spin: PropTypes.bool.isRequired,
  alt: PropTypes.string
}

Image.defaultProps = {
  spin: false
}

export default withStyles(styles)(Image)
