import React from 'react'
import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'
import MenuItem from '@material-ui/core/MenuItem'
import Typography from '@material-ui/core/Typography'
import { grey } from '@material-ui/core/colors'

const styles = ({ palette }: Theme) =>
  createStyles({
    menuItem: {
      '&:hover': {
        backgroundColor: 'transparent',
      },
    },
    link: {
      color: palette.common.white,
      textDecoration: 'none',
      '&:hover': {
        color: grey[200],
      },
    },
  })

interface Props extends WithStyles<typeof styles> {
  text: string
  onClick: (e: React.MouseEvent) => void
}

const AvatarMenuItem = ({ classes, onClick, text }: Props) => {
  return (
    <MenuItem className={classes.menuItem} onClick={onClick}>
      <Typography variant="body1" className={classes.link}>
        {text}
      </Typography>
    </MenuItem>
  )
}

export default withStyles(styles)(AvatarMenuItem)
