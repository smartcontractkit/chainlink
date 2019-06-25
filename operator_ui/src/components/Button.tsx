import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import MuiButton, {
  ButtonProps as MuiButtonProps
} from '@material-ui/core/Button'
import pick from 'lodash/pick'
import classNames from 'classnames'

const styles = (theme: Theme) =>
  createStyles({
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
  })

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

const buildMuiProps = (props: any) => {
  switch (props.variant) {
    case PRIMARY:
      return PRIMARY_MUI_PROPS
    case SECONDARY:
      return SECONDARY_MUI_PROPS
    case DANGER:
      return DANGER_MUI_PROPS
    default: {
      return DEFAULT_MUI_PROPS
      // return Object.assign({}, DEFAULT_MUI_PROPS, {
      //   TouchRippleProps: {
      //     classes: {
      //       root: props.classes.defaultRipple
      //     }
      //   }
      // })
    }
  }
}

export type ButtonVariant =
  | 'text'
  | 'flat'
  | 'outlined'
  | 'contained'
  | 'raised'
  | 'fab'
  | 'extendedFab'
  | 'danger' // This seems to be a custom variant

interface IProps extends WithStyles<typeof styles> {
  children?: React.ReactNode
  className?: string
  variant?: ButtonVariant
  component?: any
  // component?: React.ReactNode
  // component?: React.ReactType<MuiButtonProps>
  type?: any
  onClick?: any
}

const Button = (props: IProps) => {
  const wprops = pick(props, PROPS_WHITELIST)
  // const muiProps = Object.assign({}, wprops, buildMuiProps(props))
  const muiProps = { type: props.type }
  const variant = props.variant || DEFAULT
  // const className = classNames(
  //   props.classes[variant as keyof typeof props.classes],
  //   props.className
  // )

  // return (
  //   <MuiButton {...muiProps} className={className}>
  //     {props.children}
  //   </MuiButton>
  // )
  return <MuiButton {...muiProps}>{props.children}</MuiButton>
  // return <MuiButton>{props.children}</MuiButton>
}

export default withStyles(styles)(Button)
