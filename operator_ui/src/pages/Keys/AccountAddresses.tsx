import React from 'react'
import Grid from '@material-ui/core/Grid'
import { v2 } from 'api'
import * as jsonapi from '@chainlink/json-api-client'
import * as presenters from 'core/store/presenters'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'
import { TimeAgo } from '@chainlink/styleguide'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'
import { Copy } from './Copy'

const styles = (theme: Theme) =>
  createStyles({
    card: {
      padding: theme.spacing.unit,
      marginBottom: theme.spacing.unit * 3,
    },
  })

export const AccountAddresses = withStyles(styles)(
  ({ classes }: WithStyles<typeof styles>) => {
    const [accountBalances, setAccountBalances] = React.useState<
      jsonapi.ApiResponse<presenters.AccountBalance[]>['data']
    >()
    const { error, ErrorComponent, setError } = useErrorHandler()
    const { isLoading, LoadingPlaceholder } = useLoadingPlaceholder(
      !error && !accountBalances,
    )

    const handleFetchIndex = React.useCallback(() => {
      v2.user.balances
        .getAccountBalances()
        .then(({ data }) => {
          setAccountBalances(data)
        })
        .catch(setError)
    }, [setError])

    React.useEffect(() => {
      handleFetchIndex()
    }, [handleFetchIndex])

    return (
      <Grid item xs={12}>
        <ErrorComponent />

        <Card className={classes.card}>
          <CardHeader title="Account addresses" />

          <CardContent>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                    <Typography variant="body1" color="textSecondary">
                      Address
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body1" color="textSecondary">
                      Type
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body1" color="textSecondary">
                      Link balance
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body1" color="textSecondary">
                      ETH balance
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body1" color="textSecondary">
                      Created
                    </Typography>
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {isLoading && (
                  <TableRow>
                    <TableCell component="th" scope="row" colSpan={5}>
                      <LoadingPlaceholder />
                    </TableCell>
                  </TableRow>
                )}

                {accountBalances?.length === 0 && (
                  <TableRow>
                    <TableCell component="th" scope="row" colSpan={5}>
                      No entries to show.
                    </TableCell>
                  </TableRow>
                )}
                {accountBalances?.map((balance) => (
                  <TableRow hover key={balance.id}>
                    <TableCell>
                      <Typography variant="body1">
                        {balance.attributes.address}{' '}
                        <Copy data={balance.attributes.address} />
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        {balance.attributes.isFunding
                          ? 'Emergency funding'
                          : 'Regular'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        {balance.attributes.linkBalance}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        {balance.attributes.ethBalance}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        <TimeAgo tooltip>
                          {balance.attributes.createdAt || ''}
                        </TimeAgo>
                      </Typography>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </Grid>
    )
  },
)
