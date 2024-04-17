# Root chart
Root chart dedicated for integration testing

## Development
### Updating dependencies
`helm dependency update`

### Testing 
#### Render template locally
```
helm template app -f values.yaml -n crib-example . --output-dir .rendered \
--set=ingress.baseDomain="$DEVSPACE_INGRESS_BASE_DOMAIN" \
--set=ccipScripts.deployContractsAndJobs=true
```

