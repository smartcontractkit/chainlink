import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import React from 'react'

interface Props {
  children: React.ReactNode
  className?: string
}

export const PaddedCard = ({ children, className }: Props) => (
  <Card>
    <CardContent className={className}>{children}</CardContent>
  </Card>
)
