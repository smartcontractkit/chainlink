import React from 'react'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'

interface Props {
  errors: ReadonlyArray<string>
}

export const ErrorsCard: React.FC<Props> = ({ errors }) => {
  return (
    <Card>
      <CardHeader title="Errors" />
      <ul>
        {errors.map((error, index) => (
          <li key={error + index}>
            <Typography variant="body1">{error}</Typography>
          </li>
        ))}
      </ul>
    </Card>
  )
}
