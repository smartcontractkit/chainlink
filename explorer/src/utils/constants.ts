export enum PostgresErrorCode {
  UNIQUE_CONSTRAINT_VIOLATION = '23505',
}

export const NORMAL_CLOSE = 1000

export const ACCESS_KEY_HEADER = 'x-explore-chainlink-accesskey'
export const SECRET_HEADER = 'x-explore-chainlink-secret'
export const ADMIN_USERNAME_HEADER = 'x-explore-admin-username'
export const ADMIN_PASSWORD_HEADER = 'x-explore-admin-password'
export const CORE_VERSION_HEADER = 'x-explore-chainlink-core-version'
export const CORE_SHA_HEADER = 'x-explore-chainlink-core-sha'

export const ADMIN_USERNAME_PARAM = 'username'
export const ADMIN_PASSWORD_PARAM = 'password'
