import React from 'react'
import Button from 'components/Button'
import { CopyToClipboard } from 'react-copy-to-clipboard'

const Copy = ({ data, buttonText }) => (
  <CopyToClipboard text={data}>
    <Button>{buttonText}</Button>
  </CopyToClipboard>
)

export default Copy
