/**
 * @packageDocumentation
 *
 * This file deals with contract helpers to deal with ethers.js contract abstractions
 */
import { ethers, Signer } from 'ethers'
import { Provider } from 'ethers/providers'
export * from './generated/LinkTokenFactory'

/**
 * The type of any function that is deployable
 */
type Deployable = {
  deploy: (...deployArgs: any[]) => Promise<any>
}

/**
 * Get the return type of a function, and unbox any promises
 */
export type Instance<T extends Deployable> = T extends {
  deploy: (...deployArgs: any[]) => Promise<infer U>
}
  ? U
  : never

type Override<T, Method extends string> = {
  [K in keyof T]: K extends Method ? any : T[K]
}

export type CallableOverrideInstance<
  T extends Deployable,
  Callables extends string
> = T extends {
  deploy: (...deployArgs: any[]) => Promise<infer ContractInterface>
}
  ? Omit<Override<ContractInterface, Callables>, 'connect'> & {
      connect(
        signer: string | Signer | Provider,
      ): Override<ContractInterface, Callables>
    }
  : never
export function callable(oldContract: ethers.Contract, methods: string[]): any {
  const oldAbi = oldContract.interface.abi
  const newAbi = oldAbi.map(fragment => {
    if (!methods.includes(fragment.name ?? '')) {
      return fragment
    }
    return {
      ...fragment,
      stateMutability: 'view',
    }
  })
  const contract = new ethers.Contract(
    oldContract.address,
    newAbi,
    oldContract.signer,
  )

  return contract
}
