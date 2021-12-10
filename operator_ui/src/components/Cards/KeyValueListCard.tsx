import React from 'react'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

const renderEntries = (entries: Array<Array<string>>) =>
  entries.map(([k, v]) => (
    <TableRow key={k}>
      <TableCell>{k}</TableCell>
      <TableCell>{String(v)}</TableCell>
    </TableRow>
  ))

const renderBody = (
  entries: Array<Array<string>>,
  loading: boolean,
  error: string,
) => {
  if (error) {
    return <ErrorRow>{error}</ErrorRow>
  }

  if (loading) {
    return <FetchingRow />
  }

  return renderEntries(entries)
}

const SpanRow: React.FC = ({ children }) => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={3}>
      {children}
    </TableCell>
  </TableRow>
)

const FetchingRow = () => <SpanRow>...</SpanRow>

const ErrorRow: React.FC = ({ children }) => <SpanRow>{children}</SpanRow>

export interface Props {
  entries: Array<Array<any>>
  loading: boolean
  showHead?: boolean
  title?: string
  error?: string
}

export const KeyValueListCard = ({
  loading,
  entries,
  error = '',
  showHead = false,
  title,
}: Props) => (
  <Card>
    {title && <CardHeader title={title} />}

    <Table>
      {showHead && (
        <TableHead>
          <TableRow>
            <TableCell>Key</TableCell>
            <TableCell>Value</TableCell>
          </TableRow>
        </TableHead>
      )}
      <TableBody>{renderBody(entries, loading, error)}</TableBody>
    </Table>
  </Card>
)
