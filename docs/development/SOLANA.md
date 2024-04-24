# Handling Chainlink Common Breaking Changes in Chainlink Core and Chainlink Solana

The order of operation for updating seems to have been:
1. Merge chainlink-common with breaking changes
2. Create chainlink-solana PR with updated common ref
3. Create a chainlink core PR pinned to chainlink-common ref and the chainlink solana PR
3. Update chainlink-solana PR with chainlink/integration-tests + chainlink/core ref
4. Merge chainlink-solana
5. Merge chainlink core

```mermaid
%%{init: { 'gitGraph': {'parallelCommits': false,  'showBranches': true, 'showCommitLabel':false,'mainBranchName': 'common'}} }%%

gitGraph TB:
    branch core order: 2
    branch solana order: 4
    # Setup dummy commits
    checkout common
    checkout core
    checkout solana
    checkout common
    commit
    checkout core
    commit
    checkout solana
    commit

    # merge chainlink-common
    checkout common
    commit

    branch common-pr order: 1
    commit tag:"breaking changes" type: HIGHLIGHT
    checkout common
    merge common-pr

    # dummy commits for spacing
    checkout core
    commit
    checkout solana
    commit

    # create chainlink-solana PR with updated common ref
    checkout solana
    commit
    branch solana-pr order: 4
    checkout solana-pr
    commit tag: "update common ref" type: HIGHLIGHT

    # create a chainlink core PR pinned to chainlink-common ref
    # and the chainlink solana PR
    checkout core
    branch core-pr order: 3
    commit tag: "update common ref" type: HIGHLIGHT
    
    commit tag: "update solana ref" type: HIGHLIGHT
    merge solana-pr
    commit
    # update chainlink-solana PR with 
    # chainlink/integration-tests + chainlink/core ref
    checkout solana-pr
    
    commit tag: "update w/ core-integ tests + core ref" type: HIGHLIGHT
    merge core-pr

    # merge chainlink-solana
    checkout solana
    merge solana-pr type: HIGHLIGHT

    # then merge chainlink core
    checkout core
    merge core-pr type: HIGHLIGHT
    merge common
    checkout solana
    merge common
    commit
    checkout core
    commit
``` 
