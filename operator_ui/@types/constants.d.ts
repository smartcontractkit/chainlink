export type status =
  | 'In Progress'
  | 'Pending Confirmations'
  | 'Pending Connection'
  | 'Pending Bridge'
  | 'Pending Sleep'
  | 'Errored'
  | 'Completed'
  | null

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
