import { ethers } from 'ethers'

export const name = 'TEST_CONTRACT'

const events = [1, 2, 3].map((n) => ({
  topic: `topic${n}`,
  name: `event${n}`,
  decode: jest.fn().mockReturnValue({ [`input${n}`]: `value${n}` }),
  inputs: [{ name: `input${n}` }],
}))

export const eventNames = events.map((e) => e.name)

export const contract = ({
  interface: {
    events: {
      event1: events[0],
      event2: events[1],
      event3: events[2],
    },
  },
  on: jest.fn(),
} as unknown) as ethers.Contract

export const eventEmission1 = {
  topics: ['topic1'],
  data: 'eventEmissionData',
  blockNumber: 1,
}

export const eventEmission2 = {
  topics: ['topic2'],
  data: 'eventEmissionData',
  blockNumber: 1,
}
