import { GanacheProvider, JSONRPCRequestPayload } from 'ethereum-types'
import { Subprovider, Callback, ErrorCallback } from '@0x/subproviders'
import { makeDebug } from '../debug'

// TODO: Could not find a declaration file for module '@eth-optimism/ovm-toolchain'
/* eslint-disable @typescript-eslint/no-var-requires */
const ovm = require('@eth-optimism/ovm-toolchain')

const log = makeDebug('ovm-ganache-subprovider')

/**
 * This class implements the [web3-provider-engine](https://github.com/MetaMask/provider-engine) subprovider interface.
 * It intercepts all JSON RPC requests and relays them to an in-process ganache instance.
 */
export class OVMGanacheSubprovider extends Subprovider {
  private readonly ganacheProvider: GanacheProvider
  /**
   * Instantiates an OVMGanacheSubprovider
   * @param opts The desired opts with which to instantiate the Ganache provider
   */
  constructor(opts: any) {
    super()
    this.ganacheProvider = ovm.ganache.provider(opts)
  }

  /**
   * This method conforms to the web3-provider-engine interface.
   * It is called internally by the ProviderEngine when it is this subproviders
   * turn to handle a JSON RPC request.
   * @param payload JSON RPC payload
   * @param _next Callback to call if this subprovider decides not to handle the request
   * @param end Callback to call if subprovider handled the request and wants to pass back the request.
   */
  // tslint:disable-next-line:prefer-function-over-method async-suffix
  public async handleRequest(
    payload: JSONRPCRequestPayload,
    _next: Callback,
    end: ErrorCallback,
  ): Promise<void> {
    log(`sending request payload to ovm ganache: ${JSON.stringify(payload)}`)
    this.ganacheProvider.sendAsync(
      payload,
      (err: Error | null, result: any) => {
        const resultStr = JSON.stringify(result).slice(0, 1000)
        const errStr = JSON.stringify(err)
        const message = `payload sent: result ${resultStr} [possibly truncated], err ${errStr}`
        log(message)
        end(err, result && result.result)
      },
    )
  }
}
