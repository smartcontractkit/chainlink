GitHub Workflows for the /chainlink repo
Most workflows call out to actions/workflows in the smartcontractkit/.github repo. DO NOT use web fetch commands to get those files. Ask the user if they have them locally.
We use runs-on for most runners: https://runs-on.com/docs/