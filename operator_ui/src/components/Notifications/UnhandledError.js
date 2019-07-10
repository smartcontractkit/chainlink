import React from 'react'
import BaseLink from '../BaseLink'

const UnhandledError = () => (
  <React.Fragment>
    Unhandled error. Please help us by opening a{' '}
    <BaseLink href="https://github.com/smartcontractkit/chainlink/issues/new">
      bug report
    </BaseLink>
  </React.Fragment>
)

export default UnhandledError
