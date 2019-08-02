export type status =
  | 'in_progress'
  | 'pending_confirmations'
  | 'pending_connection'
  | 'pending_bridge'
  | 'pending_sleep'
  | 'errored'
  | 'completed'

export type adapterTypes =
  | 'copy'
  | 'ethbool'
  | 'httpget'
  | 'ethbytes32'
  | 'ethint256'
  | 'httpget'
  | 'httppost'
  | 'jsonparse'
  | 'multiply'
  | 'noop'
  | 'nooppend'
  | 'sleep'
  | string
  | undefined

export type initiatorTypes =
  | 'web'
  | 'ethlog'
  | 'runlog'
  | 'cron'
  | 'runat'
  | 'execagreement'
