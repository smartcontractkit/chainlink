import React from 'react'
import BaseLink from './BaseLink'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { grey } from '@material-ui/core/colors'
import { ThemeStyle } from '@material-ui/core/styles/createTypography'
import { PropTypes } from '@material-ui/core'
import classNames from 'classnames'

type Variant = ThemeStyle | 'srOnly'
type Color = PropTypes.Color | 'textPrimary' | 'textSecondary' | 'error'

const styles = () =>
  createStyles({
    link: {
      color: grey[900],
      textDecoration: 'none',
    },
    linkContent: {
      display: 'inline-block',
    },
  })

interface Props extends WithStyles<typeof styles> {
  children: React.ReactNode
  href: string
  variant?: Variant
  color?: Color
  className?: string
}

const Link = ({
  children,
  classes,
  className,
  href,
  variant,
  color,
}: Props) => {
  const v = variant || 'body1'
  const c = color || 'inherit'

  return (
    <BaseLink href={href} className={classNames(classes.link, className)}>
      <Typography variant={v} color={c} className={classes.linkContent}>
        {children}
      </Typography>
    </BaseLink>
  )
}

export default withStyles(styles)(Link)
