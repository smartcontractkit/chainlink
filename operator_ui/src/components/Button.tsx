import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import MuiButton, {
  ButtonProps as MuiButtonProps,
} from '@material-ui/core/Button'
import classNames from 'classnames'

const styles = ({ palette }: Theme) =>
  createStyles({
    default: {
      borderColor: '#BECAD6',
      '&:hover': {
        backgroundColor: palette.common.white,
        borderColor: '#BECAD6',
        boxShadow:
          '0 2px 4px 0 rgba(0,123,255,0.06), 0 2px 2px 0 rgba(0,0,0,0.06)',
      },
    },
    primary: {
      boxShadow: '0 0',
      backgroundColor: palette.primary.main,
      color: palette.common.white,
      '&:hover': {
        backgroundColor: palette.primary.main,
        boxShadow:
          '0 2px 4px 0 rgba(0,123,255,0.19), 0 2px 2px 0 rgba(0,0,0,0.15)',
      },
    },
    secondary: {
      '&:hover': {
        backgroundColor: palette.common.white,
      },
    },
    danger: {
      borderColor: palette.error.main,
      color: palette.error.main,
      '&:hover': {
        backgroundColor: palette.common.white,
        borderColor: palette.error.main,
        boxShadow:
          '0 2px 4px 0 rgba(0,123,255,0.06), 0 2px 2px 0 rgba(0,0,0,0.06)',
      },
    },
    defaultRipple: {
      color: palette.text.secondary,
    },
  })

// Unfortunately @material-ui/core does not export the type so we have to redefine it.
type MuiButtonVariant =
  | 'text'
  | 'flat'
  | 'outlined'
  | 'contained'
  | 'raised'
  | 'fab'
  | 'extendedFab'
  | undefined

export type ButtonVariant =
  | MuiButtonVariant
  | 'primary'
  | 'secondary'
  | 'danger'
  | 'default'

const muiProps = (variant: ButtonVariant, classes: any): MuiButtonProps => {
  switch (variant) {
    case 'primary':
      return { variant: 'contained' }
    case 'secondary':
      return { variant: 'outlined', color: 'primary' }
    case 'danger':
      return { variant: 'outlined', color: 'primary' }
    default: {
      return {
        variant: 'outlined',
        color: 'secondary',
        TouchRippleProps: {
          classes: {
            root: classes.defaultRipple,
          },
        },
      }
    }
  }
}

interface Props extends WithStyles<typeof styles> {
  children: React.ReactNode
  component?: React.ReactNode
  onClick?: React.MouseEventHandler<JSX.Element>
  type?: string
  disabled?: boolean
  className?: string
  variant?: ButtonVariant
  size?: 'small' | 'medium' | 'large'
  onMouseLeave?: () => void
  // Ideally this would be typed as below. However the MuiButton type annotations
  // don't allow an object to be passed through.
  //
  // href?:
  //   | string
  //   | {
  //       pathname: string
  //       state: { definition: object }
  //     }
  href?: any
}

const Button = ({
  variant = 'default',
  disabled,
  type,
  component,
  href,
  classes,
  className,
  children,
  onClick,
  onMouseLeave,
  size,
}: Props) => {
  const curryProps = Object.assign(
    { component, disabled, href, onClick, type, size, onMouseLeave },
    muiProps(variant, classes),
  )
  const cn = classNames(classes[variant as keyof typeof classes], className)

  return (
    <MuiButton {...curryProps} className={cn}>
      {children}
    </MuiButton>
  )
}

export default withStyles(styles)(Button)
