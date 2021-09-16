import React from 'react'
import * as models from 'core/store/models'
import { localizedTimestamp, TimeAgo } from 'components/TimeAgo'
import Button from 'components/Button'
import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import Link from 'components/Link'

const styles = (theme: Theme) =>
  createStyles({
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0,
    },
    mainRow: {
      marginBottom: theme.spacing.unit * 2,
    },
    actions: {
      textAlign: 'right',
    },
    regionalNavButton: {
      marginLeft: theme.spacing.unit,
      marginRight: theme.spacing.unit,
    },
    horizontalNav: {
      paddingBottom: 0,
    },
    horizontalNavItem: {
      display: 'inline',
      paddingLeft: 0,
      paddingRight: 0,
    },
    horizontalNavLink: {
      padding: `${theme.spacing.unit * 4}px ${theme.spacing.unit * 4}px`,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: theme.palette.common.white,
      '&:hover': {
        borderBottomColor: theme.palette.primary.main,
      },
    },
    activeNavLink: {
      color: theme.palette.primary.main,
      borderBottomColor: theme.palette.primary.main,
    },
    chainId: {
      overflow: 'hidden',
      textOverflow: 'ellipsis',
    },
  })

interface Chain<T> {
  attributes: T
  id: string
  type: string
}

export type ChainSpecV2 = Chain<models.Chain>

interface Props extends WithStyles<typeof styles> {
  chainId: string
  chain?: ChainSpecV2
}

const RegionalNavComponent = ({ classes, chainId, chain }: Props) => {
  const navOverridesActive = location.pathname.endsWith('/config-overrides')
  const navNodesActive = !navOverridesActive
  return (
    <>
      <Card className={classes.container}>
        <Grid container spacing={0}>
          <Grid item xs={12}>
            <Grid
              container
              spacing={0}
              alignItems="center"
              className={classes.mainRow}
            >
              <Grid item xs={6}>
                {chain && (
                  <Typography
                    variant="h5"
                    color="secondary"
                    className={classes.chainId}
                  >
                    Chain {chain.id || chainId}
                  </Typography>
                )}
              </Grid>
              <Grid item xs={6} className={classes.actions}>
                <Link href={`/chains/${chainId}/nodes/new`}>
                  <Button className={classes.regionalNavButton}>
                    Add Node
                  </Button>
                </Link>
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            {chain?.attributes.createdAt && (
              <Typography variant="subtitle2" color="textSecondary">
                Created{' '}
                <TimeAgo tooltip={false}>{chain.attributes.createdAt}</TimeAgo>{' '}
                ({localizedTimestamp(chain.attributes.createdAt)})
              </Typography>
            )}
          </Grid>
          <Grid item xs={12}>
            <List className={classes.horizontalNav}>
              <ListItem className={classes.horizontalNavItem}>
                <Link
                  href={`/chains/${chainId}`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navNodesActive && classes.activeNavLink,
                  )}
                >
                  Nodes
                </Link>
              </ListItem>
              <ListItem className={classes.horizontalNavItem}>
                <Link
                  href={`/chains/${chainId}/config-overrides`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navOverridesActive && classes.activeNavLink,
                  )}
                >
                  Config Overrides
                </Link>
              </ListItem>
            </List>
          </Grid>
        </Grid>
      </Card>
    </>
  )
}

export default withStyles(styles)(RegionalNavComponent)
