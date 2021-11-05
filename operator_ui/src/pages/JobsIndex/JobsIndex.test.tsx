/* eslint-env jest */
import React from 'react'
import { Route } from 'react-router-dom'
import {
  directRequestResource,
  jsonApiJobSpecsV2,
  fluxMonitorJobResource,
  ocrJobResource,
  keeperJobResource,
  cronJobResource,
  webJobResource,
  vrfJobResource,
} from 'support/factories/jsonApiJobs'
import globPath from 'test-helpers/globPath'
import JobsIndex, { simpleJobFilter, JobResource } from './JobsIndex'
import { ENDPOINT as OCR_ENDPOINT } from 'api/v2/jobs'
import { renderWithRouter, screen } from 'support/test-utils'

const { findAllByRole, queryByText } = screen

describe('pages/JobsIndex/JobsIndex', () => {
  it('renders the list of jobs', async () => {
    global.fetch.getOnce(
      globPath(OCR_ENDPOINT),
      jsonApiJobSpecsV2([
        ocrJobResource({
          id: '1000000',
          createdAt: new Date().toISOString(),
        }),
        fluxMonitorJobResource({
          id: '2000000',
          createdAt: new Date().toISOString(),
        }),
        directRequestResource({
          id: '3000000',
          createdAt: new Date().toISOString(),
        }),
        keeperJobResource({
          id: '4000000',
          createdAt: new Date().toISOString(),
        }),
        cronJobResource({
          id: '5000000',
          createdAt: new Date().toISOString(),
        }),
        webJobResource({
          id: '6000000',
          createdAt: new Date().toISOString(),
        }),
        vrfJobResource({
          id: '7000000',
          createdAt: new Date().toISOString(),
        }),
      ]),
    )

    renderWithRouter(
      <Route>
        <JobsIndex />
      </Route>,
    )

    expect(await findAllByRole('row')).toHaveLength(8) // Includes header

    // OCR V2 Job
    expect(queryByText('1000000')).toBeInTheDocument()

    // Flux Monitor V2 Job
    expect(queryByText('2000000')).toBeInTheDocument()

    // Direct Request V2 Job
    expect(queryByText('3000000')).toBeInTheDocument()

    // Keeper Job
    expect(queryByText('4000000')).toBeInTheDocument()

    // Cron Job
    expect(queryByText('5000000')).toBeInTheDocument()

    // Web Job
    expect(queryByText('6000000')).toBeInTheDocument()

    // VRF Job
    expect(queryByText('7000000')).toBeInTheDocument()
  })

  describe('simpleJobFilter', () => {
    it('filters by name', () => {
      const jobs: JobResource[] = [
        fluxMonitorJobResource({ name: 'FM Job' }),
        ocrJobResource({ name: 'OCR Job' }),
      ]

      expect(jobs.filter(simpleJobFilter('FM Job')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('OCR Job')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('Job')).length).toEqual(2)
    })

    it('filters by type', () => {
      const jobs: JobResource[] = [
        fluxMonitorJobResource({}),
        fluxMonitorJobResource({}),
        ocrJobResource({}),
      ]

      expect(jobs.filter(simpleJobFilter('fluxmonitor')).length).toEqual(2)
    })

    it('filters by ID', () => {
      const jobs: JobResource[] = [
        fluxMonitorJobResource({ id: 'id-3a' }),
        ocrJobResource({ id: 'id-3b' }),
      ]

      expect(jobs.filter(simpleJobFilter('i')).length).toEqual(2)
      expect(jobs.filter(simpleJobFilter('id')).length).toEqual(2)
      expect(jobs.filter(simpleJobFilter('id-')).length).toEqual(2)
      expect(jobs.filter(simpleJobFilter('id-1')).length).toEqual(0)
      expect(jobs.filter(simpleJobFilter('id-3a')).length).toEqual(1)
    })

    it('filters by job type', () => {
      const jobs: JobResource[] = [
        ocrJobResource({ id: '1' }),
        fluxMonitorJobResource({ id: '2' }),
      ]

      expect(jobs.filter(simpleJobFilter('off')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('chain')).length).toEqual(1)
    })
  })
})
