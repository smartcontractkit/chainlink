import React from 'react'
import { mount } from 'enzyme'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableRow from '@material-ui/core/TableRow'
import TableCell, {
  LinkColumn,
  TimeAgoColumn,
  TextColumn,
} from '../../../components/Table/TableCell'

const mountWithinTableRow = (component: React.ReactNode) => {
  return mount(
    <Table>
      <TableBody>
        <TableRow>{component}</TableRow>
      </TableBody>
    </Table>,
  )
}

describe('components/Table/TableCell', () => {
  it('can render text columns', () => {
    const col: TextColumn = { type: 'text', text: 'Hello' }
    const wrapper = mountWithinTableRow(<TableCell column={col} />)

    expect(wrapper.text()).toContain('Hello')
  })

  it('can render link columns', () => {
    const col: LinkColumn = {
      type: 'link',
      text: 'Hello World',
      to: '/world',
    }
    const wrapper = mountWithinTableRow(<TableCell column={col} />)

    expect(wrapper.text()).toContain('Hello World')
    expect(wrapper.find('a').props().href).toEqual('/world')
  })

  it('can render time ago columns', () => {
    const col: TimeAgoColumn = {
      type: 'time_ago',
      text: new Date().toISOString(),
    }
    const wrapper = mountWithinTableRow(<TableCell column={col} />)

    expect(wrapper.text()).toContain('just now')
  })
})
