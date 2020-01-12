import React from 'react'
import {
  createStyles,
  withStyles,
  Theme,
  WithStyles,
} from '@material-ui/core/styles'
import MenuItem from '@material-ui/core/MenuItem'

const styles = ({ palette }: Theme) =>
  createStyles({
    menuItem: {
      color: palette.common.white,
      '&:hover': {
        backgroundColor: 'transparent',
      },
    },
  })

interface Props extends WithStyles<typeof styles> {}

const UnstyledAvatarMenuItem: React.FC<Props> = ({ classes, children }) => {
  return <MenuItem className={classes.menuItem}>{children}</MenuItem>
}

export const AvatarMenuItem = withStyles(styles)(UnstyledAvatarMenuItem)
