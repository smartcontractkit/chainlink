import { ethers } from 'ethers'

export type FF<T extends ethers.Contract> = PatchConnect<FlattenFunctions<T>>

type PatchConnect<T extends ethers.Contract> = Omit<T, 'connect'> & {
  connect(...args: Parameters<T['connect']>): PatchConnect<T>
}

type FlattenFunctions<T extends ethers.Contract> = {
  [K in keyof T['functions']]: T['functions'][K]
} &
  T
