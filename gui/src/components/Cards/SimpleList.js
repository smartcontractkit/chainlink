import React from 'react'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'

const SimpleList = ({children, title}) => {
  return (
    <Card>
      <CardContent>
        <Typography variant='h5' color='secondary'>
          {title}
        </Typography>
      </CardContent>

      <Table>
        <TableBody>
          {children}
        </TableBody>
      </Table>
    </Card>
  )
}

export default SimpleList
