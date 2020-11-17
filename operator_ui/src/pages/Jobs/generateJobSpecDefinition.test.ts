import { InitiatorType } from 'core/store/models'
import { generateJSONDefinition } from './generateJobSpecDefinition'

describe('generateJSONDefinition', () => {
  it('generates valid definition', () => {
    const jobSpecAttributesInput = {
      initiators: [
        {
          type: 'web' as InitiatorType.WEB,
        },
      ],
      id: '7f4e4ca5f9ce4131a080a214947736c5',
      name: 'Bitstamp ticker',
      createdAt: '2020-11-17T10:25:44.040459Z',
      tasks: [
        {
          ID: 6,
          type: 'httpget',
          confirmations: 0,
          params: {
            get: 'https://bitstamp.net/api/ticker/',
          },
          CreatedAt: '2020-11-17T10:25:44.043094Z',
          UpdatedAt: '2020-11-17T10:25:44.043094Z',
          DeletedAt: null,
        },
        {
          ID: 7,
          type: 'jsonparse',
          confirmations: null,
          params: {
            path: ['last'],
          },
          CreatedAt: '2020-11-17T10:25:44.043948Z',
          UpdatedAt: '2020-11-17T10:25:44.043948Z',
          DeletedAt: null,
        },
        {
          ID: 8,
          type: 'multiply',
          confirmations: null,
          params: {
            times: 100,
          },
          CreatedAt: '2020-11-17T10:25:44.04456Z',
          UpdatedAt: '2020-11-17T10:25:44.04456Z',
          DeletedAt: null,
        },
        {
          ID: 9,
          type: 'ethuint256',
          confirmations: null,
          params: {},
          CreatedAt: '2020-11-17T10:25:44.045404Z',
          UpdatedAt: '2020-11-17T10:25:44.045404Z',
          DeletedAt: null,
        },
        {
          ID: 10,
          type: 'ethtx',
          confirmations: null,
          params: {},
          CreatedAt: '2020-11-17T10:25:44.046211Z',
          UpdatedAt: '2020-11-17T10:25:44.046211Z',
          DeletedAt: null,
        },
      ],
      minPayment: '1000000',
      updatedAt: '2020-02-09T15:13:03Z',
      startAt: '2020-02-09T15:13:03Z',
      endAt: null,
      errors: [],
      earnings: null,
    }

    const expectedOutput = {
      initiators: [{ type: 'web' }],
      name: 'Bitstamp ticker',
      startAt: '2020-02-09T15:13:03Z',
      tasks: [
        {
          confirmations: 0,
          params: { get: 'https://bitstamp.net/api/ticker/' },
          type: 'httpget',
        },
        { params: { path: ['last'] }, type: 'jsonparse' },
        { params: { times: 100 }, type: 'multiply' },
        { type: 'ethuint256' },
        { type: 'ethtx' },
      ],
    }

    const output = generateJSONDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })
})

describe('generateTOMLDefinition', () => {
  it('generates valid definition', () => {
    const jobSpecAttributesInput = {
      initiators: [
        {
          type: 'web' as InitiatorType.WEB,
        },
      ],
      id: '7f4e4ca5f9ce4131a080a214947736c5',
      name: 'Bitstamp ticker',
      createdAt: '2020-11-17T10:25:44.040459Z',
      tasks: [
        {
          ID: 6,
          type: 'httpget',
          confirmations: 0,
          params: {
            get: 'https://bitstamp.net/api/ticker/',
          },
          CreatedAt: '2020-11-17T10:25:44.043094Z',
          UpdatedAt: '2020-11-17T10:25:44.043094Z',
          DeletedAt: null,
        },
        {
          ID: 7,
          type: 'jsonparse',
          confirmations: null,
          params: {
            path: ['last'],
          },
          CreatedAt: '2020-11-17T10:25:44.043948Z',
          UpdatedAt: '2020-11-17T10:25:44.043948Z',
          DeletedAt: null,
        },
        {
          ID: 8,
          type: 'multiply',
          confirmations: null,
          params: {
            times: 100,
          },
          CreatedAt: '2020-11-17T10:25:44.04456Z',
          UpdatedAt: '2020-11-17T10:25:44.04456Z',
          DeletedAt: null,
        },
        {
          ID: 9,
          type: 'ethuint256',
          confirmations: null,
          params: {},
          CreatedAt: '2020-11-17T10:25:44.045404Z',
          UpdatedAt: '2020-11-17T10:25:44.045404Z',
          DeletedAt: null,
        },
        {
          ID: 10,
          type: 'ethtx',
          confirmations: null,
          params: {},
          CreatedAt: '2020-11-17T10:25:44.046211Z',
          UpdatedAt: '2020-11-17T10:25:44.046211Z',
          DeletedAt: null,
        },
      ],
      minPayment: '1000000',
      updatedAt: '2020-02-09T15:13:03Z',
      startAt: '2020-02-09T15:13:03Z',
      endAt: null,
      errors: [],
      earnings: null,
    }

    const expectedOutput = {
      initiators: [{ type: 'web' }],
      name: 'Bitstamp ticker',
      startAt: '2020-02-09T15:13:03Z',
      tasks: [
        {
          confirmations: 0,
          params: { get: 'https://bitstamp.net/api/ticker/' },
          type: 'httpget',
        },
        { params: { path: ['last'] }, type: 'jsonparse' },
        { params: { times: 100 }, type: 'multiply' },
        { type: 'ethuint256' },
        { type: 'ethtx' },
      ],
    }

    const output = generateJSONDefinition(jobSpecAttributesInput)
    expect(output).toEqual(expectedOutput)
  })
})
