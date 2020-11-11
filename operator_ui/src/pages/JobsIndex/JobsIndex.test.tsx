/* eslint-env jest */
import React from 'react'
import { Route } from 'react-router-dom'
import { jsonApiJobSpecs } from 'factories/jsonApiJobSpecs'
import { jsonApiOcrJobSpecs } from 'factories/jsonApiOcrJobSpecs'
import { syncFetch } from 'test-helpers/syncFetch'
import globPath from 'test-helpers/globPath'
import { mountWithProviders } from 'test-helpers/mountWithTheme'
import JobsIndex, { simpleJobFilter } from './JobsIndex'
import { partialAsFull } from '@chainlink/ts-helpers'
import { Initiator, InitiatorType } from 'core/store/models'
import { INDEX_ENDPOINT as JSON_ENDPOINT } from 'api/v2/specs'
import { ENDPOINT as OCR_ENDPOINT } from 'api/v2/ocrSpecs'

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
      jsonApiOcrJobSpecs([
        {
          id: 'OcrId',
          createdAt: new Date().toISOString(),
        },
      ]),
    )

    const wrapper = mountWithProviders(<Route component={JobsIndex} />)

    await syncFetch(wrapper)

    expect(wrapper.text()).toContain('JsonId')
    expect(wrapper.text()).toContain('OcrId')
    expect(wrapper.text()).toContain('web')
    expect(wrapper.text()).toContain('just now')
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
      jsonApiOcrJobSpecs([
        {
          id: 'OcrId',
          createdAt: new Date().toISOString(),
        },
      ]),
    )

    const wrapper = mountWithProviders(<Route component={JobsIndex} />)

    await syncFetch(wrapper)

    // Expect to have 3 jobs initially
    expect(wrapper.find('tbody').children().length).toEqual(4)

    wrapper
      .find('input[name="search"]')
      .simulate('change', { target: { value: 'web' } })

    expect(wrapper.find('tbody').children().length).toEqual(1)
  })

  describe('simpleJobFilter', () => {
    it('filters by name', () => {
      const jobs = jsonApiJobSpecs([
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

      const search = 'Kaiko'
      expect(jobs.filter(simpleJobFilter(search)).length).toEqual(1)
    })

    it('filters by initiator', () => {
      const jobs = jsonApiJobSpecs([
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

      const search = 'cr'
      expect(jobs.filter(simpleJobFilter(search)).length).toEqual(1)
    })

    it('filters by ID', () => {
      const jobs = jsonApiJobSpecs([
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

      expect(jobs.filter(simpleJobFilter('i')).length).toEqual(4)
      expect(jobs.filter(simpleJobFilter('id')).length).toEqual(4)
      expect(jobs.filter(simpleJobFilter('id-')).length).toEqual(4)
      expect(jobs.filter(simpleJobFilter('id-1')).length).toEqual(3)
      expect(jobs.filter(simpleJobFilter('id-1c')).length).toEqual(1)
    })
  })
})
