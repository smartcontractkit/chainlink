export enum JobSpecFormats {
  JSON = 'JSON',
  TOML = 'TOML',
}

export type JobSpecFormat = keyof typeof JobSpecFormats
