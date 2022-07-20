import React from 'react'

import { CopyToClipboard } from 'react-copy-to-clipboard'

import FileCopyOutlinedIcon from '@material-ui/icons/FileCopyOutlined'
import IconButton from '@material-ui/core/IconButton'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import Tooltip from '@material-ui/core/Tooltip'

const styles = (theme: Theme) =>
  createStyles({
    button: {
      width: 27,
      height: 27,
      minHeight: 27,
      marginLeft: theme.spacing.unit,
      marginRight: theme.spacing.unit,
    },
    icon: {
      fontSize: 18,
    },
  })

export const CopyIconButton = withStyles(styles)(
  ({ classes, data }: WithStyles<typeof styles> & { data: string }) => {
    const [copied, setCopied] = React.useState(false)

    return (
      <CopyToClipboard text={data} onCopy={() => setCopied(true)}>
        <Tooltip title="Copied!" open={copied} onClose={() => setCopied(false)}>
          <IconButton className={classes.button}>
            <FileCopyOutlinedIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      </CopyToClipboard>
    )
  },
)
