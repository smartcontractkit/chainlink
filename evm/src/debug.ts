import debug from 'debug'

export function makeDebug(name: string): debug.Debugger {
  return debug(name)
}
