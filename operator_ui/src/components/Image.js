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
        transform: 'rotate(360deg)'
      }
    }
  }
}

const Image = ({ src, width, height, spin, alt, classes }) => {
  const size = {}
  if (width >= 0) {
    size.width = width
  }
  if (height >= 0) {
    size.height = height
  }

  return (
    <img
      src={src}
      className={spin ? classes.animate : ''}
      alt={alt}
      {...size}
    />
  )
}

Image.propTypes = {
  src: PropTypes.string.isRequired,
  spin: PropTypes.bool.isRequired,
  width: PropTypes.number,
  height: PropTypes.number,
  alt: PropTypes.string
}

Image.defaultProps = {
  spin: false
}

export default withStyles(styles)(Image)
