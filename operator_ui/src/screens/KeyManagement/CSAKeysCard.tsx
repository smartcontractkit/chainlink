import React from 'react'

import { gql } from '@apollo/client'

import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

import { CSAKeyRow } from './CSAKeyRow'
import { ErrorRow } from 'src/components/TableRow/ErrorRow'
import { LoadingRow } from 'src/components/TableRow/LoadingRow'
import { NoContentRow } from 'src/components/TableRow/NoContentRow'

export const CSA_KEYS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment CSAKeysPayload_ResultsFields on CSAKey {
    id
    publicKey
  }
`

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

export interface Props extends WithStyles<typeof styles> {
  loading: boolean
  data?: FetchCsaKeys
  errorMsg?: string
  onCreate: () => void
}

export const CSAKeysCard = withStyles(styles)(
  ({ classes, data, errorMsg, loading, onCreate }: Props) => {
    return (
      <Card>
        <CardHeader
          action={
            data?.csaKeys.results?.length === 0 && (
              <Button variant="outlined" color="primary" onClick={onCreate}>
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
                <TableCell>Public Key</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              <LoadingRow visible={loading} />
              <NoContentRow visible={data?.csaKeys.results?.length === 0} />
              <ErrorRow msg={errorMsg} />

              {data?.csaKeys.results?.map((key, idx) => (
                <CSAKeyRow csaKey={key} key={idx} />
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    )
  },
)
