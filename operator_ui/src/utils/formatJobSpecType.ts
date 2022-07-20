// formatJobSpecType formats the typename of a spec to a readable format.
export const formatJobSpecType = (typename: string) => {
  switch (typename) {
    case 'DirectRequestSpec':
      return 'Direct Request'
    case 'FluxMonitorSpec':
      return 'Flux Monitor'
    default:
      return typename.replace(/Spec$/, '')
  }
}
