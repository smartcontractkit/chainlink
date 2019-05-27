import React from 'react'
import PropTypes from 'prop-types'
import pick from 'lodash/pick'
import Image from './Image'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => {
  return {
    text: {
      color: theme.palette.primary.main,
      display: 'inline-block',
      marginLeft: theme.spacing(2),
      paddingTop: theme.spacing(1),
      verticalAlign: 'top'
    }
  }
})

const Logo = props => {
  const classes = useStyles()
  const propsExt = { ...classes, ...props }
  const imageProps = pick(propsExt, ['src', 'width', 'height', 'alt'])
  return <Image {...imageProps} />
}

Logo.propTypes = {
  src: PropTypes.string.isRequired,
  width: PropTypes.number,
  height: PropTypes.number,
  alt: PropTypes.string
}

export default Logo
