import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Icon from '@material-ui/core/Icon'
import classNames from 'classnames'

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
  className?: string
}

const url = (host: string, txHash: string) => `https://${host}/tx/${txHash}`

const EtherscanLink = ({ classes, host, txHash, className }: IProps) => {
  return (
    <a
      href={url(host, txHash)}
      target="_blank"
      rel="noopener noreferrer"
      className={classNames(classes.link, className)}>
      <Typography variant="body1" className={classes.linkText} inline>
        {txHash}
      </Typography>

      <Icon color="action">launch</Icon>
    </a>
  )
}

export default withStyles(styles)(EtherscanLink)
