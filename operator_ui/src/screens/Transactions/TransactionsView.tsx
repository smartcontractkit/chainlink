import React from 'react'

import { useHistory } from 'react-router-dom'
import { gql } from '@apollo/client'

import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

import Content from 'src/components/Content'
import { Heading1 } from 'src/components/Heading/Heading1'
import { Loading } from 'src/components/Feedback/Loading'
import { TransactionRow } from './TransactionRow'
import { TablePagination } from '@material-ui/core'

export const ETH_TRANSACTIONS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment EthTransactionsPayload_ResultsFields on EthTransaction {
    chain {
      id
    }
    from
    hash
    to
    nonce
    sentAt
  }
`

export interface Props {
  data?: FetchEthTransactions
  loading: boolean
  page: number
  pageSize: number
}

export const TransactionsView: React.FC<Props> = ({
  data,
  loading,
  page,
  pageSize,
}) => {
  const history = useHistory()

  return (
    <Content>
      <Grid container spacing={32}>
        <Grid item xs={12}>
          <Heading1>Transactions</Heading1>
        </Grid>

        {loading && <Loading />}

        {data && (
          <Grid item xs={12}>
            <Card>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Txn Hash</TableCell>
                    <TableCell>Chain ID</TableCell>
                    <TableCell>From</TableCell>
                    <TableCell>To</TableCell>
                    <TableCell>Nonce</TableCell>
                    <TableCell>Block</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {data?.ethTransactions.results.length === 0 && (
                    <TableRow>
                      <TableCell component="th" scope="row" colSpan={6}>
                        You haven’t created any transactions yet
                      </TableCell>
                    </TableRow>
                  )}

                  {data?.ethTransactions.results.map((tx, idx) => (
                    <TransactionRow tx={tx} key={idx} />
                  ))}
                </TableBody>
              </Table>
              <TablePagination
                component="div"
                count={data.ethTransactions.metadata.total}
                rowsPerPage={pageSize}
                rowsPerPageOptions={[10, 25, 50, 100]}
                page={page - 1}
                onChangePage={(_, p) => {
                  history.push(`/transactions?page=${p + 1}&per=${pageSize}`)
                }}
                onChangeRowsPerPage={(e) => {
                  history.push(
                    `/transactions?page=${page}&per=${parseInt(
                      e.target.value,
                      10,
                    )}`,
                  )
                }}
                backIconButtonProps={{ 'aria-label': 'prev-page' }}
                nextIconButtonProps={{ 'aria-label': 'next-page' }}
              />
            </Card>
          </Grid>
        )}
      </Grid>
    </Content>
  )
}
