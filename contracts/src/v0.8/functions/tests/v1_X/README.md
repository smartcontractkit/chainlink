## Usage

First set the foundry profile to Functions:
```
export FOUNDRY_PROFILE=functions
```

**To run tests use**:

All Functions test files:
```
forge test -vvv
```

To run a specific file use:
```
forge test -vvv --mp src/v0.8/functions/tests/v1_X/[File Name].t.sol 
```

**To see coverage**:
First ensure that the correct files are being evaluated. For example, if only v0 contracts are, then temporarily change the Functions profile in `./foundry.toml`.

```
forge coverage --ir-minimum
```