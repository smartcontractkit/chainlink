import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => ({
  card: {
    paddingTop: theme.spacing.unit * 2,
    paddingBottom: theme.spacing.unit * 2,
    paddingLeft: theme.spacing.unit * 3,
    paddingRight: theme.spacing.unit * 3
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
