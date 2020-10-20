/**
 * @packageDocumentation
 *
 * This file deals with contract helpers to deal with ethers.js contract abstractions
 */
import { ethers, Signer, ContractTransaction } from 'ethers'
import { Provider } from 'ethers/providers'
import { FunctionFragment } from 'ethers/utils'
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

type Override<T> = {
  [K in keyof T]: T[K] extends (...args: any[]) => Promise<ContractTransaction>
    ? (...args: any[]) => Promise<any>
    : T[K]
}

export type CallableOverrideInstance<T extends Deployable> = T extends {
  deploy: (...deployArgs: any[]) => Promise<infer ContractInterface>
}
  ? Omit<Override<ContractInterface>, 'connect'> & {
      connect(signer: string | Signer | Provider): CallableOverrideInstance<T>
    }
  : never

export function callable(oldContract: ethers.Contract, methods: string[]): any {
  const oldAbi = oldContract.interface.abi
  const newAbi = oldAbi.map((fragment) => {
    if (!methods.includes(fragment.name ?? '')) {
      return fragment
    }

    if ((fragment as FunctionFragment)?.constant === false) {
      return {
        ...fragment,
        stateMutability: 'view',
        constant: true,
      }
    }
    return {
      ...fragment,
      stateMutability: 'view',
    }
  })
  const contract = new ethers.Contract(
    oldContract.address,
    newAbi,
    oldContract.signer ?? oldContract.provider,
  )

  return contract
}
