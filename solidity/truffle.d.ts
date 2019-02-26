declare type _contractTest = (accounts: string[]) => void
declare function contract(name: string, test: _contractTest): void

/* tslint:disable-next-line:interface-name */
declare interface TransactionMeta {
  from: string,
}

declare interface IContract<T> {
  'new'(): Promise<T>,
  deployed(): Promise<T>,
  at(address: string): T,
}

/* tslint:disable-next-line:interface-name no-empty-interface */
declare interface ChainlinkedInstance {}

/* tslint:disable-next-line:interface-name */
interface Artifacts {
  require(name: string): Contract<ChainlinkedInstance>,
}

declare var artifacts: Artifacts
