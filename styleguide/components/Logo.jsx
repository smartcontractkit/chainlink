import React from 'react'
import PropTypes from 'prop-types'
import pick from 'lodash/pick'
import { withStyles } from '@material-ui/core/styles'
import Image from './Image'

const styles = theme => {
  return {
    text: {
      color: theme.palette.primary.main,
      display: 'inline-block',
      marginLeft: theme.spacing.unit * 2,
      paddingTop: theme.spacing.unit,
      verticalAlign: 'top'
    }
  }
}

const Logo = props => {
  const imageProps = pick(props, ['src', 'width', 'height', 'alt'])
  return <Image {...imageProps} />
}

Logo.propTypes = {
  src: PropTypes.string.isRequired,
  width: PropTypes.number,
  height: PropTypes.number,
  alt: PropTypes.string
}

export default withStyles(styles)(Logo)
