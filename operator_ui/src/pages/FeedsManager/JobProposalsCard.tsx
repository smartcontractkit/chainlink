import React from 'react'

import { v2 } from 'api'
import Link from 'components/Link'
import * as models from 'core/store/models'
import { tableStyles } from 'components/Table'

import Card from '@material-ui/core/Card'
import CardHeader from '@material-ui/core/CardHeader'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import Tab from '@material-ui/core/Tab'
import Tabs from '@material-ui/core/Tabs'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import { TimeAgo } from 'src/components/TimeAgo'

const tabToStatus: { [key: number]: string } = {
  0: 'pending',
  1: 'approved',
  2: 'rejected',
}

const styles = () => {
  return createStyles({
    tabsRoot: {
      borderBottom: '1px solid #e8e8e8',
    },
  })
}

interface JobProposal {
  attributes: models.JobProposal
  id: string
  type: string
}

interface JobProposalRowProps extends WithStyles<typeof tableStyles> {
  proposal: JobProposal
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

        <TableCell>{proposal.attributes.external_job_id || 'N/A'}</TableCell>
        <TableCell>
          <TimeAgo tooltip>{proposal.attributes.createdAt}</TimeAgo>
        </TableCell>
      </TableRow>
    )
  },
)

interface Props extends WithStyles<typeof styles> {}

export const JobProposalsCard = withStyles(styles)(({ classes }: Props) => {
  const [tabValue, setTabValue] = React.useState(0)
  const [proposals, setProposals] = React.useState<JobProposal[]>()

  React.useEffect(() => {
    v2.jobProposals.getJobProposals().then((res) => {
      setProposals(res.data)
    })
  }, [])

  const filteredProposals: JobProposal[] = React.useMemo(() => {
    if (!proposals) {
      return []
    }

    return proposals.filter(
      (p) => p.attributes.status === tabToStatus[tabValue],
    )
  }, [tabValue, proposals])

  return (
    <Card>
      <CardHeader title="Job Proposals" />

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
  )
})
