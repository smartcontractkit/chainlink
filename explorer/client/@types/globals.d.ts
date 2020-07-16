declare const __EXPLORER_SERVER_VERSION__: string
declare const __EXPLORER_CLIENT_VERSION__: string
declare const __GIT_SHA__: string
declare const __GIT_BRANCH__: string

declare global {
  namespace NodeJS {
    export interface ProcessEnv {
      NODE_ENV: 'development' | 'production' | 'test'
      REACT_APP_EXPLORER_BASEURL?: string
      REACT_APP_EXPLORER_GA_ID?: string
    }
  }
}
