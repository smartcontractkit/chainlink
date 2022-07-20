declare type JSONValue = JSONPrimitive | JSONObject | JSONArray
declare type JSONPrimitive = string | number | boolean | null
declare type JSONObject = { [member: string]: JSONValue }
declare interface JSONArray extends Array<JSONValue> {}
declare type Pointer<T> = T | null

declare module 'math/big' {
  export type Int = string
}

declare module 'net/url' {
  export type URL = string
}

declare module 'github.com/smartcontractkit/chainlink/core/null' {
  export type Uint32 = string | null
}

declare module 'gopkg.in/guregu/null.v3' {
  export type String = string | null
  export type Time = string | null
}

declare module 'github.com/jinzhu/gorm' {
  import * as time from 'time'

  // FIXME -- needs camelCase
  export interface Model {
    ID: number
    CreatedAt: time.Time
    UpdatedAt: time.Time
    DeletedAt: Pointer<time.Time>
  }
}

declare module 'core/web/sessions_controller' {
  export interface Session {
    authenticated: boolean
  }
}

declare module 'github.com/ethereum/go-ethereum/common/hexutil' {
  /**
   * Bytes marshals/unmarshals as a JSON string with 0x prefix.
   * The empty slice marshals as "0x".
   */
  export type Bytes = string
}

declare module 'go.uber.org/zap/zapcore' {
  /**
   * A Level is a logging priority. Higher levels are more important.
   */
  export type Level = number
}

declare module 'core/store/orm' {
  import * as zapcore from 'go.uber.org/zap/zapcore'

  /**
   * LogLevel determines the verbosity of the events to be logged.
   */
  export interface LogLevel {
    Level: zapcore.Level
  }
}

declare module 'time' {
  export type Time = string
  /**
   *  FIXME -- We should define a custom marshaler for this
   */
  export type Duration = number
}

declare module 'github.com/ethereum/go-ethereum/common' {
  /**
   * Hash represents the 32 byte Keccak256 hash of arbitrary data.
   */
  export type Hash = string

  /**
   * Address represents the 20 byte address of an Ethereum account.
   */
  export type Address = string
}
