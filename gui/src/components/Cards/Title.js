import React from 'react'
import PropTypes from 'prop-types'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'

const Title = ({ children, divider }) => (
  <React.Fragment>
    <CardContent>
      <Typography variant="h5" color="secondary">
        {children}
      </Typography>
    </CardContent>

    {divider && <Divider />}
  </React.Fragment>
)

Title.propTypes = {
  divider: PropTypes.bool.isRequired
}

Title.defaultProps = {
  divider: false
}

export default Title
