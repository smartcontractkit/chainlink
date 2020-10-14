const { run, sleep } = require('./utils')

module.exports = {
  chainlinkLogin,
  addJobSpec,
  setPriceFeedValue,
  successfulJobRuns,
  makeJobSpecEthlog,
  makeJobSpecFluxMonitor,
}

async function chainlinkLogin(nodeURL, tmpdir) {
  let resp = await run(
    `chainlink -j admin login -f ${__dirname}/../../chainlink/apicredentials`,
    { CLIENT_NODE_URL: nodeURL, ROOT: tmpdir },
  )
  console.log(resp)
}

async function addJobSpec(nodeURL, jobSpec, tmpdir) {
  let jobSpecJSON = JSON.stringify(jobSpec)
  let resp = await run(['chainlink', '-j', 'jobs', 'create', jobSpecJSON], {
    CLIENT_NODE_URL: nodeURL,
    ROOT: tmpdir,
  })
  try {
    console.log('RESP ~>', resp)
    return JSON.parse(resp)
  } catch (err) {
    throw new Error('error adding job spec: ' + err.toString())
  }
}

async function setPriceFeedValue(externalAdapterURL, value) {
  const url = new URL('result', externalAdapterURL).href
  const response = await fetch(url, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ result: value }),
  })
  if (!response.ok) {
    throw new Error('failed to set external adapter price feed')
  }
}

async function successfulJobRuns(nodeURL, jobSpecID, num, tmpdir) {
  while (true) {
    let runs = JSON.parse(
      await run(`chainlink -j runs list`, {
        CLIENT_NODE_URL: nodeURL,
        ROOT: tmpdir,
      }),
    )
    console.log('runs ~>', runs)
    const failedRun = (run) =>
      run.jobId === jobSpecID && run.status === 'errored'
    if (Array.isArray(runs) && runs.find(failedRun)) {
      throw new Error('job run errored')
    }
    await sleep(3000)
  }
}

function makeJobSpecEthlog(oracleContractAddress) {
  return {
    initiators: [
      { type: 'ethlog', params: { address: oracleContractAddress } },
    ],
    tasks: [
      { type: 'httpGet', params: { url: 'http://localhost:8000/data' } },
      { type: 'jsonParse', params: { path: ['bryn', 'age'] } },
      { type: 'multiply', params: { times: 10000 } },
      { type: 'ethtx' },
    ],
  }
}

const initiators = [
  {
    type: 'fluxmonitor',
    params: {
      address: aggregatorContractAddress,
      requestData: {
        data: {
          coin: 'ETH',
          market: 'USD',
        },
      },
      feeds: [feedAddr],
      precision: 2,
      threshold: 5,
      idleTimer: {
        disabled: true,
      },
      pollTimer: {
        period: '15s',
      },
    },
  },
]

const tasks = [
  {
    type: 'multiply',
    confirmations: null,
    params: {
      times: 100,
    },
  },
  {
    type: 'ethuint256',
    confirmations: null,
    params: {},
  },
  {
    type: 'ethtx',
    confirmations: null,
    params: {},
  },
]

function makeJobSpecFluxMonitor(aggregatorContractAddress, feedAddr) {
  return {
    initiators: initiators,
    tasks: tasks,
  }
}
