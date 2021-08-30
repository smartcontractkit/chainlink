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
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import JobsIndex, { simpleJobFilter, JobSpecV2 } from './JobsIndex'
import { ENDPOINT as OCR_ENDPOINT } from 'api/v2/jobs'

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

    const wrapper = mountWithProviders(<Route component={JobsIndex} />)

    await syncFetch(wrapper)

    // OCR V2 Job
    expect(wrapper.text()).toContain('1000000')

    // Flux Monitor V2 Job
    expect(wrapper.text()).toContain('2000000')

    // Direct Request V2 Job
    expect(wrapper.text()).toContain('3000000')

    // Keeper V2 Job
    expect(wrapper.text()).toContain('4000000')
  })

  it('allows searching', async () => {
    global.fetch.getOnce(
      globPath(OCR_ENDPOINT),
      jsonApiJobSpecsV2([
        ocrJobResource({
          id: 'OcrId',
          createdAt: new Date().toISOString(),
        }),
        fluxMonitorJobResource({
          id: 'FluxMonitorId',
          createdAt: new Date().toISOString(),
        }),
      ]),
    )

    const wrapper = mountWithProviders(<Route component={JobsIndex} />)

    await syncFetch(wrapper)

    // Expect to have 2 jobs
    expect(wrapper.find('tbody').children().length).toEqual(2)
  })

  describe('simpleJobFilter', () => {
    it('filters by name', () => {
      const jobs: JobSpecV2[] = [
        fluxMonitorJobResource({ name: 'FM Job' }),
        ocrJobResource({ name: 'OCR Job' }),
      ]

      expect(jobs.filter(simpleJobFilter('FM Job')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('OCR Job')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('Job')).length).toEqual(2)
    })

    it('filters by type', () => {
      const jobs: JobSpecV2[] = [
        fluxMonitorJobResource({}),
        fluxMonitorJobResource({}),
        ocrJobResource({}),
      ]

      expect(jobs.filter(simpleJobFilter('fluxmonitor')).length).toEqual(2)
    })

    it('filters by ID', () => {
      const jobs: JobSpecV2[] = [
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
      const jobs: JobSpecV2[] = [
        ocrJobResource({ id: '1' }),
        fluxMonitorJobResource({ id: '2' }),
      ]

      expect(jobs.filter(simpleJobFilter('off')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('chain')).length).toEqual(1)
    })
  })
})
