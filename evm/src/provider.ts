import { IEthereumProvider } from '@nomiclabs/buidler/types'
import { JsonRpcProvider } from 'ethers/providers'

// https://github.com/nomiclabs/buidler/blob/master/packages/buidler-ethers/src/ethers-provider-wrapper.ts
export class EthersProviderWrapper extends JsonRpcProvider {
  private readonly _buidlerProvider: IEthereumProvider
  constructor(buidlerProvider: IEthereumProvider) {
    super()
    this._buidlerProvider = buidlerProvider
  }
  public async send(method: string, params: any): Promise<any> {
    const result = await this._buidlerProvider.send(method, params)
    // We replicate ethers' behavior.
    this.emit('debug', {
      action: 'send',
      request: {
        id: 42,
        jsonrpc: '2.0',
        method,
        params,
      },
      response: result,
      provider: this,
    })
    return result
  }
}
