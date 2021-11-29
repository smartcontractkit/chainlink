import React from 'react'

import { gql } from '@apollo/client'

import Card from '@material-ui/core/Card'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import Tab from '@material-ui/core/Tab'
import Tabs from '@material-ui/core/Tabs'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import { TimeAgo } from 'src/components/TimeAgo'

import Link from 'components/Link'
import { SearchTextField } from 'src/components/Search/SearchTextField'
import { tableStyles } from 'components/Table'

export const FEEDS_MANAGER__JOB_PROPOSAL_FIELDS = gql`
  fragment FeedsManager_JobProposalsFields on JobProposal {
    id
    externalJobID
    proposedAt
    status
  }
`

const tabToStatus: { [key: number]: string } = {
  0: 'PENDING',
  1: 'APPROVED',
  2: 'REJECTED',
  3: 'CANCELLED',
}

const styles = () => {
  return createStyles({
    tabsRoot: {
      borderBottom: '1px solid #e8e8e8',
    },
  })
}

interface JobProposalRowProps extends WithStyles<typeof tableStyles> {
  proposal: FeedsManager_JobProposalsFields
}

const JobProposalRow = withStyles(tableStyles)(
  ({ proposal, classes }: JobProposalRowProps) => {
    return (
      <TableRow className={classes.row} hover>
        <TableCell className={classes.cell} component="th" scope="row">
          <Link className={classes.link} href={`/job_proposals/${proposal.id}`}>
            {proposal.id}
          </Link>
        </TableCell>

        <TableCell>{proposal.externalJobID || 'N/A'}</TableCell>
        <TableCell>
          <TimeAgo tooltip>{proposal.proposedAt}</TimeAgo>
        </TableCell>
      </TableRow>
    )
  },
)

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

  const dataset: string[] = [proposal.id]
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

    const filteredProposals: FeedsManager_JobProposalsFields[] =
      React.useMemo(() => {
        if (!proposals) {
          return []
        }

        return proposals.filter(
          (p) => p.status === tabToStatus[tabValue] && search(searchTerm)(p),
        )
      }, [tabValue, proposals, searchTerm])

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
            <Tab label="Pending" />
            <Tab label="Approved" />
            <Tab label="Rejected" />
            <Tab label="Cancelled" />
          </Tabs>

          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>External Job ID</TableCell>
                <TableCell>Proposed At</TableCell>
              </TableRow>
            </TableHead>

            <TableBody>
              {filteredProposals?.map((proposal) => (
                <JobProposalRow key={proposal.id} proposal={proposal} />
              ))}
            </TableBody>
          </Table>
        </Card>
      </>
    )
  },
)
