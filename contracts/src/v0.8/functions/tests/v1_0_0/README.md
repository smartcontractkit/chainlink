## Usage

First set the foundry profile to Functions:
```
export FOUNDRY_PROFILE=functions
```

To run all test files use:
```
forge test -vv
```

To run a specific file use:
```
forge test -vv --mp src/v0.8/functions/tests/v1_0_0/[File Name].t.sol 
```

To see coverage:
First ensure that the correct files are being evaluated. For example, if only v1 contracts are, then temporarily change the Functions profile in `./foundry.toml`.
```
forge coverage
```