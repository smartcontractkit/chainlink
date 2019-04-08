import { titleCase } from 'change-case'

export default (status: string) => {
  switch (status) {
    case 'in_progress':
      return 'Pending'
    case 'error':
      return 'Failed'
    case 'completed':
      return 'Complete'
    default:
      return titleCase(status)
  }
}
