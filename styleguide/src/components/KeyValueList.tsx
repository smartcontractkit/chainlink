import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import { titleCase } from 'change-case'
import React from 'react'
import { CardTitle } from './Cards/Title'

const renderKey = (k: string, titleize: boolean) =>
  titleize ? titleCase(k) : k

const renderEntries = (entries: Array<Array<string>>, titleize: boolean) =>
  entries.map(([k, v]) => (
    <TableRow key={k}>
      <Col>{renderKey(k, titleize)}</Col>
      <Col>{v}</Col>
    </TableRow>
  ))

const renderBody = (
  entries: Array<Array<string>>,
  error: string,
  titleize: boolean
) => {
  if (error) {
    return <ErrorRow>{error}</ErrorRow>
  } else if (entries.length === 0) {
    return <FetchingRow />
  } else {
    return renderEntries(entries, titleize)
  }
}

interface ISpanRowProps {
  children: React.ReactNode
}

const SpanRow = ({ children }: ISpanRowProps) => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={3}>
      {children}
    </TableCell>
  </TableRow>
)

const FetchingRow = () => <SpanRow>...</SpanRow>

interface IErrorRowProps {
  children: React.ReactNode
}

const ErrorRow = ({ children }: IErrorRowProps) => <SpanRow>{children}</SpanRow>

interface IColProps {
  children: React.ReactNode
}

const Col = ({ children }: IColProps) => (
  <TableCell>
    <Typography variant="body1">
      <span>{children}</span>
    </Typography>
  </TableCell>
)

interface IHeadColProps {
  children: React.ReactNode
}

const HeadCol = ({ children }: IHeadColProps) => (
  <TableCell>
    <Typography variant="body1" color="textSecondary">
      {children}
    </Typography>
  </TableCell>
)

interface IProps {
  entries: Array<Array<string>>
  titleize: boolean
  showHead: boolean
  title?: string
  error?: string
}

export const KeyValueList = ({
  entries,
  error = '',
  showHead = false,
  title,
  titleize = false
}: IProps) => (
  <Card>
    {title && <CardTitle divider>{title}</CardTitle>}

    <Table>
      {showHead && (
        <TableHead>
          <TableRow>
            <HeadCol>Key</HeadCol>
            <HeadCol>Value</HeadCol>
          </TableRow>
        </TableHead>
      )}
      <TableBody>{renderBody(entries, error, titleize)}</TableBody>
    </Table>
  </Card>
)
