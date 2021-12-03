import React from 'react'

import Avatar from '@material-ui/core/Avatar'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import ListItemIcon from '@material-ui/core/ListItemIcon'
import ListItemText from '@material-ui/core/ListItemText'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import VpnKeyRoundedIcon from '@material-ui/icons/VpnKeyRounded'

const styles = (theme: Theme) =>
  createStyles({
    listItemPrimary: {
      marginBottom: theme.spacing.unit,
    },
    listItemSecondary: {
      color: theme.palette.grey[600],
    },
    avatar: {
      backgroundColor: theme.palette.primary.main,
    },
  })

interface Props extends WithStyles<typeof styles> {
  primary: React.ReactNode
  secondary: React.ReactNode[]
}

// KeyBundle provides a component which styles the key bundle information in a
// cell.
export const KeyBundle = withStyles(styles)(
  ({ classes, primary, secondary }: Props) => {
    return (
      <List dense={true}>
        <ListItem>
          <ListItemIcon>
            <Avatar className={classes.avatar}>
              <VpnKeyRoundedIcon />
            </Avatar>
          </ListItemIcon>
          <ListItemText
            primary={
              <Typography
                className={classes.listItemPrimary}
                variant="body1"
                component="span"
              >
                {primary}
              </Typography>
            }
            secondary={secondary.map((secondaryItem: any, index: number) => (
              <Typography
                key={index}
                className={classes.listItemSecondary}
                variant="subtitle2"
                component="span"
              >
                {secondaryItem}
              </Typography>
            ))}
          />
        </ListItem>
      </List>
    )
  },
)
