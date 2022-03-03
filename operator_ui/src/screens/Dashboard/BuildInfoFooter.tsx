import React from 'react'

import Paper from '@material-ui/core/Paper'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import { extractBuildInfo } from 'utils/extractBuildInfo'

const styles = (theme: Theme) =>
  createStyles({
    style: {
      textAlign: 'center',
      padding: theme.spacing.unit * 2.5,
      position: 'fixed',
      left: '0',
      bottom: '0',
      width: '100%',
      borderRadius: 0,
    },
    bareAnchor: {
      color: theme.palette.common.black,
      textDecoration: 'none',
    },
  })

interface Props extends WithStyles<typeof styles> {}

export const BuildInfoFooter = withStyles(styles)(({ classes }: Props) => {
  const { version, sha } = extractBuildInfo()

  return (
    <Paper className={classes.style}>
      <Typography>
        Chainlink Node {version} at commit{' '}
        <a
          target="_blank"
          rel="noopener noreferrer"
          href={`https://github.com/smartcontractkit/chainlink/commit/${sha}`}
          className={classes.bareAnchor}
        >
          {sha}
        </a>
      </Typography>
    </Paper>
  )
})
