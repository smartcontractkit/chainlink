import React from 'react'
import PropTypes from 'prop-types'
import Image from './Image'
import logo from '../images/chainlink-operator-logo.svg'
import { withStyles } from '@material-ui/core/styles'

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

const Logo = ({width, height, classes}) => {
  const size = {width, height}

  return (
    <Image
      src={logo}
      alt='Chainlink Operator'
      {...size}
    />
  )
}

Logo.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number
}

export default withStyles(styles)(Logo)
