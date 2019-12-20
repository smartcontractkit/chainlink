import { createMuiTheme } from '@material-ui/core/styles'
import { storiesOf } from '@storybook/react'
import React from 'react'
import { muiTheme } from 'storybook-addon-material-ui'
import { Table, theme } from '../src'

const customTheme = createMuiTheme(theme)
const DEFAULT_HEADERS = ['Node', 'ID', 'Created', 'Status']
const statuses = ['errored', 'pending', 'success']

const randomStatus = statuses[Math.floor(Math.random() * statuses.length)]
const buildNodeCol = () => [
  { type: 'text', text: 'Sample Text A' },
  { type: 'text', text: 'Sample Text B' },
  {
    type: 'time_ago',
    text: new Date(Date.now() - 86400000),
  },
  {
    type: 'status',
    text: randomStatus,
  },
]
const defaultRows = (count: number) => Array(count).fill(buildNodeCol())
console.log(defaultRows(5))

storiesOf('Table', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Filled Table', () => (
    <React.Fragment>
      <Table
        headers={DEFAULT_HEADERS}
        currentPage={0}
        rows={defaultRows(5)}
        count={10}
        onChangePage={() => {}}
        emptyMsg={'Empty Msg'}
      />
    </React.Fragment>
  ))
  .add('Empty Table', () => (
    <React.Fragment>
      <Table
        headers={DEFAULT_HEADERS}
        currentPage={0}
        rows={[]}
        count={10}
        onChangePage={() => {}}
        emptyMsg={'Empty Msg'}
      />
    </React.Fragment>
  ))
