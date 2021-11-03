export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Time: any;
};

export type Bridge = {
  readonly __typename?: 'Bridge';
  readonly confirmations: Scalars['Int'];
  readonly createdAt: Scalars['Time'];
  readonly minimumContractPayment: Scalars['String'];
  readonly name: Scalars['String'];
  readonly outgoingToken: Scalars['String'];
  readonly url: Scalars['String'];
};

/** BridgePayload defines the response to fetch a single bridge by name */
export type BridgePayload = Bridge | NotFoundError;

/** BridgesPayload defines the response when fetching a page of bridges */
export type BridgesPayload = PaginatedPayload & {
  readonly __typename?: 'BridgesPayload';
  readonly metadata: PaginationMetadata;
  readonly results: ReadonlyArray<Bridge>;
};

export type Chain = {
  readonly __typename?: 'Chain';
  readonly createdAt: Scalars['Time'];
  readonly enabled: Scalars['Boolean'];
  readonly id: Scalars['ID'];
  readonly nodes: ReadonlyArray<Node>;
  readonly updatedAt: Scalars['Time'];
};

/** CreateBridgeInput defines the input to create a bridge */
export type CreateBridgeInput = {
  readonly confirmations: Scalars['Int'];
  readonly minimumContractPayment: Scalars['String'];
  readonly name: Scalars['String'];
  readonly url: Scalars['String'];
};

/** CreateBridgeInput defines the response when creating a bridge */
export type CreateBridgePayload = CreateBridgeSuccess;

/** CreateBridgeSuccess defines the success response when creating a bridge */
export type CreateBridgeSuccess = {
  readonly __typename?: 'CreateBridgeSuccess';
  readonly bridge: Bridge;
  readonly incomingToken: Scalars['String'];
};

export type CreateFeedsManagerInput = {
  readonly bootstrapPeerMultiaddr?: Maybe<Scalars['String']>;
  readonly isBootstrapPeer: Scalars['Boolean'];
  readonly jobTypes: ReadonlyArray<JobType>;
  readonly name: Scalars['String'];
  readonly publicKey: Scalars['String'];
  readonly uri: Scalars['String'];
};

/** CreateFeedsManagerPayload defines the response when creating a feeds manager */
export type CreateFeedsManagerPayload = CreateFeedsManagerSuccess | InputErrors | NotFoundError | SingleFeedsManagerError;

/**
 * CreateFeedsManagerSuccess defines the success response when creating a feeds
 * manager
 */
export type CreateFeedsManagerSuccess = {
  readonly __typename?: 'CreateFeedsManagerSuccess';
  readonly feedsManager: FeedsManager;
};

export type Error = {
  readonly code: ErrorCode;
  readonly message: Scalars['String'];
};

/**
 * TODO - Add Chain config into the response.
 * Config    types.ChainCfg `json:"config"`
 */
export type ErrorCode =
  | 'INVALID_INPUT'
  | 'NOT_FOUND'
  | 'UNPROCESSABLE';

export type FeedsManager = {
  readonly __typename?: 'FeedsManager';
  readonly bootstrapPeerMultiaddr?: Maybe<Scalars['String']>;
  readonly createdAt: Scalars['Time'];
  readonly id: Scalars['ID'];
  readonly isBootstrapPeer: Scalars['Boolean'];
  readonly isConnectionActive: Scalars['Boolean'];
  readonly jobTypes: ReadonlyArray<JobType>;
  readonly name: Scalars['String'];
  readonly publicKey: Scalars['String'];
  readonly uri: Scalars['String'];
};

/** FeedsManagerPayload defines the response to fetch a single feeds manager by id */
export type FeedsManagerPayload = FeedsManager | NotFoundError;

/** FeedsManagersPayload defines the response when fetching feeds managers */
export type FeedsManagersPayload = {
  readonly __typename?: 'FeedsManagersPayload';
  readonly results: ReadonlyArray<FeedsManager>;
};

export type InputError = Error & {
  readonly __typename?: 'InputError';
  readonly code: ErrorCode;
  readonly message: Scalars['String'];
  readonly path: Scalars['String'];
};

export type InputErrors = {
  readonly __typename?: 'InputErrors';
  readonly errors: ReadonlyArray<InputError>;
};

export type JobType =
  | 'FLUX_MONITOR'
  | 'OCR';

export type Mutation = {
  readonly __typename?: 'Mutation';
  readonly createBridge: CreateBridgePayload;
  readonly createFeedsManager: CreateFeedsManagerPayload;
  readonly updateBridge: UpdateBridgePayload;
  readonly updateFeedsManager: UpdateFeedsManagerPayload;
};


export type MutationCreateBridgeArgs = {
  input: CreateBridgeInput;
};


export type MutationCreateFeedsManagerArgs = {
  input: CreateFeedsManagerInput;
};


export type MutationUpdateBridgeArgs = {
  input: UpdateBridgeInput;
  name: Scalars['String'];
};


export type MutationUpdateFeedsManagerArgs = {
  id: Scalars['ID'];
  input: UpdateFeedsManagerInput;
};

export type Node = {
  readonly __typename?: 'Node';
  readonly chain: Chain;
  readonly createdAt: Scalars['Time'];
  readonly httpURL: Scalars['String'];
  readonly id: Scalars['ID'];
  readonly name: Scalars['String'];
  readonly updatedAt: Scalars['Time'];
  readonly wsURL: Scalars['String'];
};

export type NotFoundError = Error & {
  readonly __typename?: 'NotFoundError';
  readonly code: ErrorCode;
  readonly message: Scalars['String'];
};

export type PaginatedPayload = {
  readonly metadata: PaginationMetadata;
};

export type PaginationMetadata = {
  readonly __typename?: 'PaginationMetadata';
  readonly total: Scalars['Int'];
};

export type Query = {
  readonly __typename?: 'Query';
  readonly bridge: BridgePayload;
  readonly bridges: BridgesPayload;
  readonly chain: Chain;
  readonly chains: ReadonlyArray<Chain>;
  readonly feedsManager: FeedsManagerPayload;
  readonly feedsManagers: FeedsManagersPayload;
};


export type QueryBridgeArgs = {
  name: Scalars['String'];
};


export type QueryBridgesArgs = {
  limit?: Maybe<Scalars['Int']>;
  offset?: Maybe<Scalars['Int']>;
};


export type QueryChainArgs = {
  id: Scalars['ID'];
};


export type QueryChainsArgs = {
  limit?: Maybe<Scalars['Int']>;
  offset?: Maybe<Scalars['Int']>;
};


export type QueryFeedsManagerArgs = {
  id: Scalars['ID'];
};

export type SingleFeedsManagerError = Error & {
  readonly __typename?: 'SingleFeedsManagerError';
  readonly code: ErrorCode;
  readonly message: Scalars['String'];
};

/** UpdateBridgeInput defines the input to update a bridge */
export type UpdateBridgeInput = {
  readonly confirmations: Scalars['Int'];
  readonly minimumContractPayment: Scalars['String'];
  readonly name: Scalars['String'];
  readonly url: Scalars['String'];
};

/** CreateBridgeInput defines the response when updating a bridge */
export type UpdateBridgePayload = NotFoundError | UpdateBridgeSuccess;

/** UpdateBridgeSuccess defines the success response when updating a bridge */
export type UpdateBridgeSuccess = {
  readonly __typename?: 'UpdateBridgeSuccess';
  readonly bridge: Bridge;
};

export type UpdateFeedsManagerInput = {
  readonly bootstrapPeerMultiaddr?: Maybe<Scalars['String']>;
  readonly isBootstrapPeer: Scalars['Boolean'];
  readonly jobTypes: ReadonlyArray<JobType>;
  readonly name: Scalars['String'];
  readonly publicKey: Scalars['String'];
  readonly uri: Scalars['String'];
};

/** UpdateFeedsManagerPayload defines the response when updating a feeds manager */
export type UpdateFeedsManagerPayload = InputErrors | NotFoundError | UpdateFeedsManagerSuccess;

/**
 * UpdateFeedsManagerSuccess defines the success response when updating a feeds
 * manager
 */
export type UpdateFeedsManagerSuccess = {
  readonly __typename?: 'UpdateFeedsManagerSuccess';
  readonly feedsManager: FeedsManager;
};


    declare global {
      export type FetchFeedsManagersVariables = Exact<{ [key: string]: never; }>;


export type FetchFeedsManagers = { readonly __typename?: 'Query', readonly feedsManagers: { readonly __typename?: 'FeedsManagersPayload', readonly results: ReadonlyArray<{ readonly __typename: 'FeedsManager', readonly id: string, readonly name: string, readonly uri: string, readonly publicKey: string, readonly jobTypes: ReadonlyArray<JobType>, readonly isBootstrapPeer: boolean, readonly isConnectionActive: boolean, readonly bootstrapPeerMultiaddr?: string | null | undefined, readonly createdAt: any }> } };

export type UpdateFeedsManagerVariables = Exact<{
  id: Scalars['ID'];
  input: UpdateFeedsManagerInput;
}>;


export type UpdateFeedsManager = { readonly __typename?: 'Mutation', readonly updateFeedsManager: { readonly __typename?: 'InputErrors', readonly errors: ReadonlyArray<{ readonly __typename?: 'InputError', readonly path: string, readonly message: string, readonly code: ErrorCode }> } | { readonly __typename?: 'NotFoundError', readonly message: string, readonly code: ErrorCode } | { readonly __typename?: 'UpdateFeedsManagerSuccess', readonly feedsManager: { readonly __typename?: 'FeedsManager', readonly id: string, readonly name: string, readonly uri: string, readonly publicKey: string, readonly jobTypes: ReadonlyArray<JobType>, readonly isBootstrapPeer: boolean, readonly isConnectionActive: boolean, readonly bootstrapPeerMultiaddr?: string | null | undefined, readonly createdAt: any } } };

export type CreateFeedsManagerVariables = Exact<{
  input: CreateFeedsManagerInput;
}>;


export type CreateFeedsManager = { readonly __typename?: 'Mutation', readonly createFeedsManager: { readonly __typename?: 'CreateFeedsManagerSuccess', readonly feedsManager: { readonly __typename?: 'FeedsManager', readonly id: string, readonly name: string, readonly uri: string, readonly publicKey: string, readonly jobTypes: ReadonlyArray<JobType>, readonly isBootstrapPeer: boolean, readonly isConnectionActive: boolean, readonly bootstrapPeerMultiaddr?: string | null | undefined, readonly createdAt: any } } | { readonly __typename?: 'InputErrors', readonly errors: ReadonlyArray<{ readonly __typename?: 'InputError', readonly path: string, readonly message: string, readonly code: ErrorCode }> } | { readonly __typename?: 'NotFoundError', readonly message: string, readonly code: ErrorCode } | { readonly __typename?: 'SingleFeedsManagerError', readonly message: string, readonly code: ErrorCode } };

    }