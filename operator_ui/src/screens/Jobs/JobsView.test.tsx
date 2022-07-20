import * as React from 'react'

import { Route } from 'react-router-dom'
import { renderWithRouter, screen } from 'support/test-utils'
import userEvent from '@testing-library/user-event'

import { buildJobs } from 'support/factories/gql/fetchJobs'
import { JobsView, Props as JobsViewProps } from './JobsView'

const { findByText, getAllByRole, getByRole, getByText, queryByText } = screen

function renderComponent(viewProps: JobsViewProps) {
  renderWithRouter(
    <>
      <Route exact path="/jobs">
        <JobsView {...viewProps} />
      </Route>
      <Route exact path="/jobs/new">
        New Jobs Page
      </Route>
    </>,
    { initialEntries: ['/jobs'] },
  )
}

describe('JobsView', () => {
  it('renders the jobs table', () => {
    const jobs = buildJobs()

    renderComponent({
      jobs,
      total: jobs.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    expect(queryByText('ID')).toBeInTheDocument()
    expect(queryByText('Name')).toBeInTheDocument()
    expect(queryByText('Type')).toBeInTheDocument()
    expect(queryByText('External Job ID')).toBeInTheDocument()
    expect(queryByText('Created')).toBeInTheDocument()

    expect(queryByText('1')).toBeInTheDocument()
    expect(queryByText('job 1')).toBeInTheDocument()
    expect(queryByText('Flux Monitor')).toBeInTheDocument()
    expect(
      queryByText('00000000-0000-0000-0000-000000000001'),
    ).toBeInTheDocument()

    expect(queryByText('2')).toBeInTheDocument()
    expect(queryByText('job 2')).toBeInTheDocument()
    expect(queryByText('OCR')).toBeInTheDocument()
    expect(
      queryByText('00000000-0000-0000-0000-000000000002'),
    ).toBeInTheDocument()

    expect(queryByText('1-2 of 2'))
  })

  it('navigates to the new jobs page', async () => {
    const jobs = buildJobs()

    renderComponent({
      jobs,
      total: jobs.length,
      page: 1,
      pageSize: 10,
    })

    userEvent.click(getByText(/new job/i))

    expect(await findByText('New Jobs Page')).toBeInTheDocument()
  })

  it('searches by standard fields', () => {
    const jobs = buildJobs()

    renderComponent({
      jobs,
      total: jobs.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    const searchInput = getByRole('textbox')

    // No match
    userEvent.paste(searchInput, '-1')
    expect(getAllByRole('row')).toHaveLength(1)

    // Type search
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, 'flux')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('Flux Monitor')).toBeInTheDocument()

    userEvent.clear(searchInput)
    userEvent.paste(searchInput, 'ocr')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('OCR')).toBeInTheDocument()

    // External job id search
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, jobs[0].externalJobID)
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('Flux Monitor')).toBeInTheDocument()

    // Name search
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, jobs[0].name)
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('Flux Monitor')).toBeInTheDocument()
  })

  it('searches by OCR fields', () => {
    const jobs = buildJobs()

    renderComponent({
      jobs,
      total: jobs.length,
      page: 1,
      pageSize: 10,
    })

    expect(getAllByRole('row')).toHaveLength(3)

    const searchInput = getByRole('textbox')

    // No match
    userEvent.paste(searchInput, '-1')
    expect(getAllByRole('row')).toHaveLength(1)

    // Contract Address search
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, '0x0000000000000000000000000000000000000001')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('OCR')).toBeInTheDocument()

    // Key bundle id search
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, 'keybundleid')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('OCR')).toBeInTheDocument()

    // transmitter search
    userEvent.clear(searchInput)
    userEvent.paste(searchInput, 'transmitteraddress')
    expect(getAllByRole('row')).toHaveLength(2)
    expect(getByText('OCR')).toBeInTheDocument()
  })
})
