import React from 'react'
import { Link } from 'react-router-dom'

const UnhandledError = () => (
  <React.Fragment>
    Unhandled error. Please help us by opening a{' '}
    <Link to="https://github.com/smartcontractkit/chainlink/issues/new">
      bug report
    </Link>
  </React.Fragment>
)

export default UnhandledError
