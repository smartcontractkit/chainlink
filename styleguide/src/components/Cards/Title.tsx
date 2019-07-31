import React from 'react'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'

interface IProps {
  children: string
  divider?: boolean
}

const Title = ({ children, divider = false }: IProps) => {
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

export default Title
