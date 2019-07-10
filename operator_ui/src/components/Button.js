import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import MuiButton from '@material-ui/core/Button'
import pick from 'lodash/pick'
import classNames from 'classnames'

const styles = theme => {
  return {
    default: {
      borderColor: '#BECAD6',
      '&:hover': {
        backgroundColor: theme.palette.common.white,
        borderColor: '#BECAD6',
        boxShadow:
          '0 2px 4px 0 rgba(0,123,255,0.06), 0 2px 2px 0 rgba(0,0,0,0.06)'
      }
    },
    primary: {
      boxShadow: '0 0',
      backgroundColor: theme.palette.primary.main,
      color: theme.palette.common.white,
      '&:hover': {
        backgroundColor: theme.palette.primary.main,
        boxShadow:
          '0 2px 4px 0 rgba(0,123,255,0.19), 0 2px 2px 0 rgba(0,0,0,0.15)'
      }
    },
    secondary: {
      '&:hover': {
        backgroundColor: theme.palette.common.white
      }
    },
    danger: {
      borderColor: theme.palette.error.main,
      color: theme.palette.error.main,
      '&:hover': {
        backgroundColor: theme.palette.common.white,
        borderColor: theme.palette.error.main,
        boxShadow:
          '0 2px 4px 0 rgba(0,123,255,0.06), 0 2px 2px 0 rgba(0,0,0,0.06)'
      }
    },
    defaultRipple: {
      color: '#818EA3'
    }
  }
}

const PROPS_WHITELIST = ['component', 'disabled', 'onClick', 'type', 'href']
const DEFAULT = 'default'
const PRIMARY = 'primary'
const SECONDARY = 'secondary'
const DANGER = 'danger'
const VARIANTS = [DEFAULT, PRIMARY, SECONDARY, DANGER]

const PRIMARY_MUI_PROPS = { variant: 'contained' }
const SECONDARY_MUI_PROPS = {
  variant: 'outlined',
  color: 'primary'
}
const DANGER_MUI_PROPS = {
  variant: 'outlined',
  color: 'primary'
}
const DEFAULT_MUI_PROPS = {
  variant: 'outlined',
  color: 'secondary'
}

const buildMuiProps = props => {
  switch (props.variant) {
    case PRIMARY:
      return PRIMARY_MUI_PROPS
    case SECONDARY:
      return SECONDARY_MUI_PROPS
    case DANGER:
      return DANGER_MUI_PROPS
    default: {
      return Object.assign({}, DEFAULT_MUI_PROPS, {
        TouchRippleProps: {
          classes: {
            root: props.classes.defaultRipple
          }
        }
      })
    }
  }
}

const Button = props => {
  const wprops = pick(props, PROPS_WHITELIST)
  const muiProps = Object.assign({}, wprops, buildMuiProps(props))
  const className = classNames(props.classes[props.variant], props.className)

  return (
    <MuiButton {...muiProps} className={className}>
      {props.children}
    </MuiButton>
  )
}

Button.propTypes = {
  variant: PropTypes.oneOf(VARIANTS)
}

Button.defaultProps = {
  variant: DEFAULT
}

export default withStyles(styles)(Button)
