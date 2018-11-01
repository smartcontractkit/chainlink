import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Image from './Image'
import logo from '../images/icon-logo-blue.svg'
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

const Logo = ({width, height, classes}) => (
  <div>
    <Image
      src={logo}
      width={width}
      height={height}
      alt='Chainlink Operator'
    />
    <Typography variant='headline' color='inherit' className={classes.text}>
      Chainlink Operator
    </Typography>
  </div>
)

Logo.propTypes = {
  width: PropTypes.number.isRequired,
  height: PropTypes.number.isRequired
}

export default withStyles(styles)(Logo)
