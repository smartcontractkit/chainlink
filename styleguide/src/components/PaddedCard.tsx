import React from 'react'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'

interface IProps {
  children: React.ReactNode
  className?: string
}

const PaddedCard = ({ children, className }: IProps) => (
  <Card>
    <CardContent className={className}>{children}</CardContent>
  </Card>
)

export default PaddedCard
