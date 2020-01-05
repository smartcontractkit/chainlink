import React from 'react'
import { Button } from '@chainlink/styleguide'
import { CopyToClipboard } from 'react-copy-to-clipboard'

const Copy = ({ data, buttonText }) => (
  <CopyToClipboard text={data}>
    <Button>{buttonText}</Button>
  </CopyToClipboard>
)

export default Copy
