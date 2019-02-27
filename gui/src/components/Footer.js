import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import extractBuildInfo from 'utils/extractBuildInfo'

const styles = theme => ({
  style: {
    textAlign: 'center',
    padding: theme.spacing.unit * 2.5,
    position: 'fixed',
    left: '0',
    bottom: '0',
    width: '100%'
  }
})

const { version, sha } = extractBuildInfo()

const Footnote = ({ classes }) => {
  return (
    <Card className={classes.style}>
      <Typography>
        Chainlink Node {version} at commit {sha}
      </Typography>
    </Card>
  )
}

export default withStyles(styles)(Footnote)
