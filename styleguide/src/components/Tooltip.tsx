import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import MuiTooltip from '@material-ui/core/Tooltip'

const styles = ({ palette, shadows, typography }: Theme) => createStyles({
  lightTooltip: {
    background: palette.primary.contrastText,
    color: palette.text.primary,
    boxShadow: shadows[24],
    ...typography.h6
  }
})

interface IProps extends WithStyles<typeof styles> {
  children: React.ReactElement<any>
  title: string
}

const Tooltip = ({ title, children, classes }: IProps) => {
  return (
    <MuiTooltip title={title} classes={{ tooltip: classes.lightTooltip }}>
      {children}
    </MuiTooltip>
  )
}

export default withStyles(styles)(Tooltip)
