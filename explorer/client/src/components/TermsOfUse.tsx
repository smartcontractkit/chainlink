import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'

const styles = ({ breakpoints, spacing, palette }: Theme) =>
  createStyles({
    container: {
      paddingLeft: spacing.unit * 2,
      paddingRight: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        paddingLeft: spacing.unit * 3,
        paddingRight: spacing.unit * 3
      },
      textAlign: 'right'
    },
    link: {
      color: palette.grey['500']
    }
  })

interface IProps extends WithStyles<typeof styles> {}

const TermsOfUse = withStyles(styles)(({ classes }: IProps) => {
  return (
    <div className={classes.container}>
      <a href="https://chain.link/terms/" className={classes.link}>
        Terms of Use
      </a>
    </div>
  )
})

export default TermsOfUse
