import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
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

const MetaInfo = ({title, value, classes}) => (
  <Card className={classes.card}>
    <Typography gutterBottom variant='headline' component='h2'>
      {title}
    </Typography>
    <Typography variant='display2' color='inherit'>
      {value}
    </Typography>
  </Card>
)

MetaInfo.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ]),
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(MetaInfo)
