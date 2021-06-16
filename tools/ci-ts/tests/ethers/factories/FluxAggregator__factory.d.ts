import { ContractFactory, Signer } from 'ethers'
import { Provider } from 'ethers/providers'
import { UnsignedTransaction } from 'ethers/utils/transaction'
import { BigNumberish } from 'ethers/utils'
import { TransactionOverrides } from '..'
import { FluxAggregator } from '../FluxAggregator'
export declare class FluxAggregator__factory extends ContractFactory {
  constructor(signer?: Signer)
  deploy(
    _link: string,
    _paymentAmount: BigNumberish,
    _timeout: BigNumberish,
    _validator: string,
    _minSubmissionValue: BigNumberish,
    _maxSubmissionValue: BigNumberish,
    _decimals: BigNumberish,
    _description: string,
    overrides?: TransactionOverrides,
  ): Promise<FluxAggregator>
  getDeployTransaction(
    _link: string,
    _paymentAmount: BigNumberish,
    _timeout: BigNumberish,
    _validator: string,
    _minSubmissionValue: BigNumberish,
    _maxSubmissionValue: BigNumberish,
    _decimals: BigNumberish,
    _description: string,
    overrides?: TransactionOverrides,
  ): UnsignedTransaction
  attach(address: string): FluxAggregator
  connect(signer: Signer): FluxAggregator__factory
  static connect(
    address: string,
    signerOrProvider: Signer | Provider,
  ): FluxAggregator
}
//# sourceMappingURL=FluxAggregator__factory.d.ts.map
