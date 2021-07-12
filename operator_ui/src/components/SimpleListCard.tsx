import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import React from 'react'
import { CardTitle } from './CardTitle'

interface Props {
  children: React.ReactNode
  title: string
}

export const SimpleListCard = ({ children, title }: Props) => {
  return (
    <Card>
      <CardTitle>{title}</CardTitle>

      <Table style={{ tableLayout: 'fixed' }}>
        <TableBody>{children}</TableBody>
      </Table>
    </Card>
  )
}
