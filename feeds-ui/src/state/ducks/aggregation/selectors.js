import { createSelector } from 'reselect'

const NODE_NAMES = [
  {
    address: '0x049bd8c3adc3fe7d3fc2a44541d955a537c2a484',
    name: 'Fiews',
  },
  {
    address: '0x240bae5a27233fd3ac5440b5a598467725f7d1cd',
    name: 'LinkPool',
  },
  {
    address: '0x4565300c576431e5228e8aa32642d5739cf9247d',
    name: 'Certus One',
  },
  {
    address: '0x58c69aff4df980357034ea98aad35bbf78cbd849',
    name: 'Wetez',
  },
  {
    address: '0x79c6e11be1c1ed4d91fbe05d458195a2677f14a5',
    name: 'Validation Capital',
  },
  {
    address: '0x7a9d706b2a3b54f7cf3b5f2fcf94c5e2b3d7b24b',
    name: 'LinkForest',
  },
  {
    address: '0x7e94a8a23687d8c7058ba5625db2ce358bcbd244',
    name: 'SNZPool',
  },
  {
    address: '0x89f70fa9f439dbd0a1bc22a09befc56ada04d9b4',
    name: 'Chainlink',
  },
  {
    address: '0x8c85a06eb3854df0d502b2b00169dbfb8b603bf3',
    name: 'Kaiko',
  },
  {
    address: '0x8cfb1d4269f0daa003cdea567ac8f76c0647764a',
    name: 'Simply VC',
  },
  {
    address: '0xb92ec7d213a28e21b426d79ede3c9bbcf6917c09',
    name: 'stake.fish',
  },
  {
    address: '0xf3b450002c7bc300ea03c9463d8e8ba7f821b7c6',
    name: 'Newroad',
  },
  {
    address: '0xf5a3d443fccd7ee567000e43b23b0e98d96445ce',
    name: 'Chainlayer',
  },
  {
    address: '0x992Ef8145ab8B3DbFC75523281DaD6A0981891bb',
    name: 'Figment Networks',
  },
  {
    address: '0x83dA1beEb89Ffaf56d0B7C50aFB0A66Fb4DF8cB1',
    name: 'Omniscience',
  },
  {
    address: '0x0Ce0224ba488ffC0F46bE32b333a874Eb775c613',
    name: 'Cosmostation',
  },
  {
    address: '0x64FE692be4b42F4Ac9d4617aB824E088350C11C2',
    name: 'Ztake.org',
  },
  {
    address: '0x260A96cEC05328f678754D1ACF143C8ac1DF079A',
    name: 'HashQuark',
  },
  {
    address: '0x38b6ab6B9294CCe1Ccb59c3e7D390690B4c18B1A',
    name: 'Prophet',
  },
  {
    address: '0x2Ed7E9fCd3c0568dC6167F0b8aEe06A02CD9ebd8',
    name: 'Secure Data Links',
  },
  {
    address: '0x78E76126719715Eddf107cD70f3A31dddF31f85A',
    name: 'Honeycomb.market',
  },
  {
    address: '0x24A718307Ce9B2420962fd5043fb876e17430934',
    name: 'Infinity Stones',
  },
]

const oracles = state => state.aggregation.oracles
const oracleResponse = state => state.aggregation.oracleResponse
const currentAnswer = state => state.aggregation.currentAnswer
const contractAddress = state => state.aggregation.contractAddress

const oraclesList = createSelector([oracles], list => {
  if (!list) return []

  const names = {}

  NODE_NAMES.forEach(n => {
    names[n.address.toUpperCase()] = n.name
  })

  const result = list.map(a => {
    return {
      address: a,
      name: names[a.toUpperCase()] || 'Unknown',
      type: 'oracle',
    }
  })

  return result
})

const networkGraphNodes = createSelector(
  [oraclesList, contractAddress],
  (list, address) => {
    if (!list) return []

    let result = [
      {
        type: 'contract',
        name: 'Aggregation Contract',
        address,
      },
      ...list,
    ]

    result = result.map((a, i) => {
      return { ...a, id: i }
    })

    return result
  },
)

const networkGraphState = createSelector(
  [oracleResponse, currentAnswer],
  (list, answer) => {
    if (!list) return []

    const contractData = {
      currentAnswer: answer,
      type: 'contract',
    }

    return [...list, contractData]
  },
)

export { oraclesList, networkGraphNodes, networkGraphState }
