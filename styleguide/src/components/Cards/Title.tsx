import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'
import React from 'react'

interface Props {
  children: string
  divider?: boolean
}

export const CardTitle = ({ children, divider = false }: Props) => {
  return (
    <React.Fragment>
      <CardContent>
        <Typography variant="h5" color="secondary">
          {children}
        </Typography>
      </CardContent>

      {divider && <Divider />}
    </React.Fragment>
  )
}
