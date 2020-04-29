export default {
  initiators: [
    {
      type: 'fluxmonitor',
      params: {
        address: '0x0000000000000000000000000000000000000000', // set before use
        requestData: {
          data: {
            coin: 'ETH',
            market: 'USD',
          },
        },
        feeds: [''], // set before use
        precision: 2,
        threshold: 5,
        idleTimer: {
          disabled: true,
        },
        pollTimer: {
          frequency: '5s'
        }


      },
    },
  ],
  tasks: [
    {
      type: 'multiply',
      confirmations: null,
      params: {
        times: 100,
      },
    },
    {
      type: 'ethint256',
      confirmations: null,
      params: {},
    },
    {
      type: 'ethtx',
      confirmations: null,
      params: {},
    },
  ],
  startAt: null,
  endAt: null,
}
