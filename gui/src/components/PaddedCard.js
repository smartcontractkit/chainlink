import React from 'react'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'

const PaddedCard = ({ children, className }) => (
  <Card>
    <CardContent className={className}>{children}</CardContent>
  </Card>
)

export default PaddedCard
