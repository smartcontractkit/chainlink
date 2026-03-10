# System Level Smoke Tests

Read and setup [devenv](../README.md).

Enter the shell: `cl sh`, select the product environment you need then run the corresponding test suite:
```bash
obs up -f # spin up full observability stack in case you need to run soak/load/chaos tests
# spin up DF1 (OCR2) environment with product orchestration
up env.toml,products/ocr2/basic.toml
# run smoke tests
test ocr2 TestSmoke/rounds
```