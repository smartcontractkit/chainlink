import React from 'react'

import IconButton from '@material-ui/core/IconButton'
import Menu from '@material-ui/core/Menu'
import SettingsIcon from '@material-ui/icons/Settings'
import { Theme, withStyles, WithStyles } from '@material-ui/core/styles'

import { MenuItemLink } from 'components/MenuItemLink'

const styles = (theme: Theme) => {
  return {
    iconButton: {
      fontSize: 32,
      color: theme.palette.primary.main,
    },
    menuList: {
      paddingTop: 0,
      paddingBottom: 0,
    },
  }
}

interface Props extends WithStyles<typeof styles> {}

export const SettingsMenu = withStyles(styles)(({ classes }: Props) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const handleOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  return (
    <React.Fragment>
      <IconButton disableRipple onClick={handleOpen}>
        <SettingsIcon className={classes.iconButton} />
      </IconButton>

      <Menu
        id="settings-menu"
        anchorEl={anchorEl}
        getContentAnchorEl={null}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        transformOrigin={{ vertical: 8, horizontal: 128 }}
        open={Boolean(anchorEl)}
        onClose={handleClose}
        disableAutoFocusItem
        MenuListProps={{
          className: classes.menuList,
        }}
      >
        <MenuItemLink onClick={handleClose} to="/keys">
          Key Management
        </MenuItemLink>
        <MenuItemLink onClick={handleClose} to="/config">
          Configuration
        </MenuItemLink>
      </Menu>
    </React.Fragment>
  )
})
