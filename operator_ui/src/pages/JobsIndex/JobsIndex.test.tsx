/* eslint-env jest */
import React from 'react'
import { Route } from 'react-router-dom'
import { jsonApiJobSpecs } from 'factories/jsonApiJobSpecs'
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
import JobsIndex, { simpleJobFilter, CombinedJobs } from './JobsIndex'
import { partialAsFull } from 'test-helpers/partialAsFull'
import { Initiator, InitiatorType } from 'core/store/models'
import { INDEX_ENDPOINT as JSON_ENDPOINT } from 'api/v2/specs'
import { ENDPOINT as OCR_ENDPOINT } from 'api/v2/jobs'

describe('pages/JobsIndex/JobsIndex', () => {
  it('renders the list of jobs', async () => {
    global.fetch.getOnce(
      globPath(JSON_ENDPOINT),
      jsonApiJobSpecs([
        {
          id: 'JsonId',
          createdAt: new Date().toISOString(),
        },
      ]),
    )

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

    // V1 Job
    expect(wrapper.text()).toContain('JsonId')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')

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
      globPath(JSON_ENDPOINT),
      jsonApiJobSpecs([
        {
          initiators: [
            partialAsFull<Initiator>({ type: 'web' as InitiatorType.WEB }),
          ],
        },
        {
          initiators: [
            partialAsFull<Initiator>({ type: 'cron' as InitiatorType.CRON }),
          ],
        },
        {
          initiators: [
            partialAsFull<Initiator>({
              type: 'run_log' as InitiatorType.RUN_LOG,
            }),
          ],
        },
      ]),
    )
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

    // Expect to have 5 jobs (3 V1 Jobs, 2 V2 Jobs)
    expect(wrapper.find('tbody').children().length).toEqual(5)

    wrapper
      .find('input[name="search"]')
      .simulate('change', { target: { value: 'web   ' } }) // Tests trimming whitespace as well

    expect(wrapper.find('tbody').children().length).toEqual(1)
  })

  describe('simpleJobFilter', () => {
    it('filters by name', () => {
      let jobs: CombinedJobs[] = jsonApiJobSpecs([
        {
          name: 'Coinmarketcap job',
        },
        {
          name: 'Kaiko job',
        },
        {
          name: 'Coinapi job',
        },
      ]).data

      jobs = jobs.concat([
        fluxMonitorJobResource({ name: 'FM Job' }),
        ocrJobResource({ name: 'OCR Job' }),
      ])

      expect(jobs.filter(simpleJobFilter('Kaiko')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('FM Job')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('OCR Job')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('Job')).length).toEqual(5)
    })

    it('filters by initiator', () => {
      let jobs: CombinedJobs[] = jsonApiJobSpecs([
        {
          initiators: [
            partialAsFull<Initiator>({ type: 'web' as InitiatorType.WEB }),
          ],
        },
        {
          initiators: [
            partialAsFull<Initiator>({ type: 'cron' as InitiatorType.CRON }),
          ],
        },
      ]).data

      jobs = jobs.concat([
        fluxMonitorJobResource({}),
        fluxMonitorJobResource({}),
        ocrJobResource({}),
      ])

      expect(jobs.filter(simpleJobFilter('cr')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('fluxmonitor')).length).toEqual(2)
    })

    it('filters by ID', () => {
      let jobs: CombinedJobs[] = jsonApiJobSpecs([
        {
          id: 'id-1a',
        },
        {
          id: 'id-1b',
        },
        {
          id: 'id-1c',
        },
        {
          id: 'id-2c',
        },
      ]).data

      jobs = jobs.concat([
        fluxMonitorJobResource({ id: 'id-3a' }),
        ocrJobResource({ id: 'id-3b' }),
      ])

      expect(jobs.filter(simpleJobFilter('i')).length).toEqual(6)
      expect(jobs.filter(simpleJobFilter('id')).length).toEqual(6)
      expect(jobs.filter(simpleJobFilter('id-')).length).toEqual(6)
      expect(jobs.filter(simpleJobFilter('id-1')).length).toEqual(3)
      expect(jobs.filter(simpleJobFilter('id-1c')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('id-3a')).length).toEqual(1)
    })

    it('filters by job type', () => {
      let jobs: CombinedJobs[] = jsonApiJobSpecs([
        {
          id: 'id-1a',
        },
        {
          id: 'id-1b',
        },
      ]).data

      jobs = jobs.concat(
        jsonApiJobSpecsV2([
          ocrJobResource({
            id: '1',
          }),
          fluxMonitorJobResource({
            id: '2',
          }),
        ]).data,
      )

      expect(jobs.filter(simpleJobFilter('direct')).length).toEqual(3)
      expect(jobs.filter(simpleJobFilter('off')).length).toEqual(1)
      expect(jobs.filter(simpleJobFilter('chain')).length).toEqual(1)
    })
  })
})
