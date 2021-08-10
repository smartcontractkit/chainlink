import React from 'react'
import Grid from '@material-ui/core/Grid'
import { v2 } from 'api'
import * as jsonapi from 'utils/json-api-client'
import * as models from 'core/store/models'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import Button from 'components/Button'

import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'
import { TimeAgo } from 'components/TimeAgo'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import { Copy } from './Copy'

const styles = () => {
  return createStyles({
    cardContent: {
      padding: 0,
      '&:last-child': {
        padding: 0,
      },
    },
  })
}

interface Props extends WithStyles<typeof styles> {}

export const CSAKeys = withStyles(styles)(({ classes }: Props) => {
  const [csaKeys, setCSAKeys] = React.useState<
    jsonapi.ApiResponse<models.CSAKey[]>['data']
  >()
  const { error, setError } = useErrorHandler()
  const [isFetching, setIsFetching] = React.useState<boolean>(true)
  const { isLoading } = useLoadingPlaceholder(!error && !csaKeys)

  // Load the CSA Keys
  React.useEffect(() => {
    let isMounted = true

    async function fetch() {
      try {
        const res = await v2.csaKeys.getCSAKeys()

        if (isMounted) {
          setCSAKeys(res.data)
        }
      } catch (e) {
        setError(e)
      } finally {
        setIsFetching(false)
      }
    }

    if (isFetching) {
      fetch()
    }

    return () => {
      isMounted = false
    }
  }, [isFetching, setCSAKeys, setError])

  const handleCreate = React.useCallback(async () => {
    try {
      await v2.csaKeys.createCSAKey()

      setIsFetching(true)
    } catch (e) {
      setError(e)
    }
  }, [setError])

  return (
    <Grid item xs={12}>
      <Card>
        <CardHeader
          action={
            !isLoading &&
            csaKeys?.length === 0 && (
              <Button
                data-testid="keys-create"
                variant="secondary"
                onClick={() => handleCreate()}
              >
                New CSA Key
              </Button>
            )
          }
          title="CSA Key"
          subheader="Manage your CSA Key"
        />
        <CardContent className={classes.cardContent}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <Typography variant="body1" color="textSecondary">
                    Public Key
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
              {csaKeys?.length === 0 && (
                <TableRow>
                  <TableCell component="th" scope="row" colSpan={5}>
                    No entries to show.
                  </TableCell>
                </TableRow>
              )}

              {csaKeys?.map((key) => (
                <TableRow hover key={key.id}>
                  <TableCell>
                    <Typography variant="body1">
                      {key.attributes.publicKey}{' '}
                      <Copy data={key.attributes.publicKey} />
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body1">
                      <TimeAgo tooltip>
                        {key.attributes.createdAt || ''}
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
})
