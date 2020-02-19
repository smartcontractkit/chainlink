import Listing from './Listing.component'
import { connect } from 'react-redux'

import { listingOperations, listingSelectors } from 'state/ducks/listing'

const mapStateToProps = state => ({
  groups: listingSelectors.groups(state),
})

const mapDispatchToProps = {
  fetchAnswers: listingOperations.fetchAnswers,
}

export default connect(mapStateToProps, mapDispatchToProps)(Listing)
