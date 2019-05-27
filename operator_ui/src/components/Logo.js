import React from 'react'
import PropTypes from 'prop-types'
import pick from 'lodash/pick'
import Image from './Image'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => {
  return {
    text: {
      color: theme.palette.primary.main,
      display: 'inline-block',
      marginLeft: theme.spacing(2),
      paddingTop: theme.spacing(1),
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
