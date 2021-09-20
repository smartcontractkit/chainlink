import { ContractFactory, Signer } from "ethers";
import { Provider } from "ethers/providers";
import { UnsignedTransaction } from "ethers/utils/transaction";
import { TransactionOverrides } from "..";
import { LinkToken } from "../LinkToken";
export declare class LinkToken__factory extends ContractFactory {
    constructor(signer?: Signer);
    deploy(overrides?: TransactionOverrides): Promise<LinkToken>;
    getDeployTransaction(overrides?: TransactionOverrides): UnsignedTransaction;
    attach(address: string): LinkToken;
    connect(signer: Signer): LinkToken__factory;
    static connect(address: string, signerOrProvider: Signer | Provider): LinkToken;
}
//# sourceMappingURL=LinkToken__factory.d.ts.map