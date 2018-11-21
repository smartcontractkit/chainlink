import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'

const styles = theme => ({
  card: {
    paddingBottom: theme.spacing.unit
  }
})

const PaddedTitleCard = ({title, children, classes}) => (
  <Card className={classes.card}>
    <CardContent>
      <Typography variant='headline' component='h2' color='secondary'>
        {title}
      </Typography>
    </CardContent>

    {children}
  </Card>
)

export default withStyles(styles)(PaddedTitleCard)
