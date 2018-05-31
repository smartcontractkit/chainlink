import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => ({
  card: {
    paddingTop: theme.spacing.unit * 5,
    paddingBottom: theme.spacing.unit * 5,
    paddingLeft: theme.spacing.unit * 4,
    paddingRight: theme.spacing.unit * 4
  }
})

const PaddedCard = ({children, classes}) => (
  <Card className={classes.card}>
    {children}
  </Card>
)

PaddedCard.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(PaddedCard)
