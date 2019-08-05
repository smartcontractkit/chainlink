import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import React from 'react'

interface IProps {
  children: React.ReactNode
  className?: string
}

export const PaddedCard = ({ children, className }: IProps) => (
  <Card>
    <CardContent className={className}>{children}</CardContent>
  </Card>
)
