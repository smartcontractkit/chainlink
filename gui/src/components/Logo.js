import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Image from './Image'
import logo from '../images/icon-logo-white.svg'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => {
  return {
    text: {
      color: theme.palette.common.white,
      display: 'inline-block',
      marginLeft: theme.spacing.unit,
      paddingTop: theme.spacing.unit,
      verticalAlign: 'top'
    }
  }
}

const Logo = ({width, height, spin, classes}) => (
  <div>
    <Image
      src={logo}
      width={width}
      height={height}
      spin={spin}
      alt='Chainlink'
    />
    <Typography variant='headline' color='inherit' className={classes.text}>
      Chainlink
    </Typography>
  </div>
)

Logo.propTypes = {
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired,
  spin: PropTypes.bool.isRequired
}

export default withStyles(styles)(Logo)
