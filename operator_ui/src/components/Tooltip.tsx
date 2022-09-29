import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import MuiTooltip from '@material-ui/core/Tooltip'
import React from 'react'

const styles = ({ palette, shadows, typography }: Theme) =>
  createStyles({
    lightTooltip: {
      background: palette.primary.contrastText,
      // @ts-expect-error
      color: palette.text.primary,
      boxShadow: shadows[24],
      ...typography.h6,
    },
  })

interface Props extends WithStyles<typeof styles> {
  children: React.ReactElement<any>
  title: string
}

const UnstyledTooltip = ({ title, children, classes }: Props) => {
  return (
    <MuiTooltip title={title} classes={{ tooltip: classes.lightTooltip }}>
      {children}
    </MuiTooltip>
  )
}

export const Tooltip = withStyles(styles)(UnstyledTooltip)
