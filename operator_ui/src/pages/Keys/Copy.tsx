import React from 'react'
import FileCopyOutlinedIcon from '@material-ui/icons/FileCopyOutlined'
import { CopyToClipboard } from 'react-copy-to-clipboard'
import Fab from '@material-ui/core/Fab'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import Tooltip from '@material-ui/core/Tooltip'

const styles = (theme: Theme) =>
  createStyles({
    fab: {
      width: 25,
      height: 25,
      minHeight: 25,
      marginLeft: theme.spacing.unit,
      marginRight: theme.spacing.unit,
    },
    icon: {
      fontSize: 18,
    },
  })

export const Copy = withStyles(styles)(
  ({ classes, data }: WithStyles<typeof styles> & { data: string }) => (
    <CopyToClipboard text={data}>
      <Tooltip title="Copy to clipboard">
        <Fab className={classes.fab}>
          <FileCopyOutlinedIcon className={classes.icon} />
        </Fab>
      </Tooltip>
    </CopyToClipboard>
  ),
)
