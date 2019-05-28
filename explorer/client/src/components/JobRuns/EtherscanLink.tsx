import React from 'react'
import { Link } from '@reach/router'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Icon from '@material-ui/core/Icon'

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    link: {
      textDecoration: 'none',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'flex-end'
    },
    linkText: {
      color: palette.primary.main,
      marginRight: spacing.unit * 4
    },
    bottomCol: {
      borderBottom: 'none'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  txHash: string
  host: string
}

const url = (host: string, txHash: string) => `https://${host}/tx/${txHash}`

const Details = ({ classes, host, txHash }: IProps) => {
  return (
    <a
      href={url(host, txHash)}
      target="_blank"
      rel="noopener noreferrer"
      className={classes.link}>
      <Typography variant="body1" className={classes.linkText} inline>
        {txHash}
      </Typography>

      <Icon color="action">launch</Icon>
    </a>
  )
}

export default withStyles(styles)(Details)
