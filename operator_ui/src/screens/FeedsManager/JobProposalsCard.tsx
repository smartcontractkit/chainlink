import React from 'react'

import Badge from '@material-ui/core/Badge'
import Card from '@material-ui/core/Card'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import Tab from '@material-ui/core/Tab'
import Tabs from '@material-ui/core/Tabs'

import { ApprovedTable } from './ApprovedTable'
import { InactiveTable } from './InactiveTable'
import { PendingTable } from './PendingTable'
import { UpdatesTable } from './UpdatesTable'
import { SearchTextField } from 'src/components/Search/SearchTextField'

const tabToStatus: { [key: number]: string } = {
  0: 'PENDING',
  1: 'UPDATES',
  2: 'APPROVED',
  3: 'REJECTED',
  4: 'CANCELLED',
}

const styles = (theme: Theme) => {
  return createStyles({
    tabsRoot: {
      borderBottom: '1px solid #e8e8e8',
    },
    badge: {
      padding: `0 ${theme.spacing.unit * 2}px`,
    },
  })
}

const searchIncludes = (searchParam: string) => {
  const lowerCaseSearchParam = searchParam.toLowerCase()

  return (stringToSearch: string) => {
    return stringToSearch.toLowerCase().includes(lowerCaseSearchParam)
  }
}

export const search =
  (term: string) => (proposal: FeedsManager_JobProposalsFields) => {
    if (term === '') {
      return true
    }

    return matchSimple(proposal, term)
  }

// matchSimple does a simple match on the id
function matchSimple(proposal: FeedsManager_JobProposalsFields, term: string) {
  const match = searchIncludes(term)

  const dataset: string[] = [proposal.id, proposal.remoteUUID]
  if (proposal.externalJobID) {
    dataset.push(proposal.externalJobID)
  }

  return dataset.some(match)
}

interface Props extends WithStyles<typeof styles> {
  proposals: readonly FeedsManager_JobProposalsFields[]
}

export const JobProposalsCard = withStyles(styles)(
  ({ classes, proposals }: Props) => {
    const [tabValue, setTabValue] = React.useState(0)
    const [searchTerm, setSearchTerm] = React.useState('')

    const tabBadgeCounts: {
      PENDING: number
      UPDATES: number
      APPROVED: number
      REJECTED: number
      CANCELLED: number
    } = React.useMemo(() => {
      const tabBadgeCounts = {
        PENDING: 0,
        UPDATES: 0,
        APPROVED: 0,
        REJECTED: 0,
        CANCELLED: 0,
      }

      proposals.forEach((p) => {
        if (p.pendingUpdate) {
          // Always display any pending update in the updates tab counter
          tabBadgeCounts['UPDATES']++

          // Support other tabs with pending updates
          switch (p.status) {
            case 'APPROVED':
              tabBadgeCounts['APPROVED']++

              break
            case 'CANCELLED':
              tabBadgeCounts['CANCELLED']++

              break
            case 'REJECTED':
              tabBadgeCounts['REJECTED']++

              break
            default:
              break
          }
        }

        if (p.status === 'PENDING') {
          tabBadgeCounts['PENDING']++
        }
      })

      return tabBadgeCounts
    }, [proposals])

    const filteredProposals: FeedsManager_JobProposalsFields[] =
      React.useMemo(() => {
        if (!proposals) {
          return []
        }

        const activeTab = tabToStatus[tabValue]

        if (activeTab === 'UPDATES') {
          return proposals
            .filter((p) => p.pendingUpdate && search(searchTerm)(p))
            .sort((a, b) => b.latestSpec.createdAt - a.latestSpec.createdAt)
        } else {
          return proposals
            .filter(
              (p) =>
                p.status === tabToStatus[tabValue] && search(searchTerm)(p),
            )
            .sort((a, b) => b.latestSpec.createdAt - a.latestSpec.createdAt)
        }
      }, [tabValue, proposals, searchTerm])

    const renderTable = (proposals: FeedsManager_JobProposalsFields[]) => {
      switch (tabToStatus[tabValue]) {
        case 'PENDING':
          return <PendingTable proposals={proposals} />
        case 'UPDATES':
          return <UpdatesTable proposals={proposals} />
        case 'REJECTED':
        case 'CANCELLED':
          return <InactiveTable proposals={proposals} />
        case 'APPROVED':
          return <ApprovedTable proposals={proposals} />
        default:
          return null
      }
    }

    return (
      <>
        <SearchTextField
          value={searchTerm}
          onChange={setSearchTerm}
          placeholder="Search job proposals"
        />
        <Card>
          <Tabs
            value={tabValue}
            className={classes.tabsRoot}
            indicatorColor="primary"
            onChange={(_, value) => {
              setTabValue(value)
            }}
          >
            <Tab
              label={
                <Badge
                  color="primary"
                  badgeContent={tabBadgeCounts.PENDING}
                  className={classes.badge}
                  data-testid="pending-badge"
                >
                  New
                </Badge>
              }
            />
            <Tab
              label={
                <Badge
                  color="primary"
                  badgeContent={tabBadgeCounts.UPDATES}
                  className={classes.badge}
                  data-testid="updates-badge"
                >
                  Updates
                </Badge>
              }
            />
            <Tab
              label={
                <Badge
                  color="primary"
                  badgeContent={tabBadgeCounts.APPROVED}
                  className={classes.badge}
                  data-testid="approved-badge"
                >
                  Approved
                </Badge>
              }
            />
            <Tab
              label={
                <Badge
                  color="primary"
                  badgeContent={tabBadgeCounts.REJECTED}
                  className={classes.badge}
                  data-testid="rejected-badge"
                >
                  Rejected
                </Badge>
              }
            />
            <Tab
              label={
                <Badge
                  color="primary"
                  badgeContent={tabBadgeCounts.CANCELLED}
                  className={classes.badge}
                  data-testid="cancelled-badge"
                >
                  Cancelled
                </Badge>
              }
            />
          </Tabs>

          {renderTable(filteredProposals)}
        </Card>
      </>
    )
  },
)
