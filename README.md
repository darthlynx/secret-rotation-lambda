# secret-rotation-lambda
AWS Lambda which rotates SecretManager secrets

## Build
```bash
# for your local host machine
make build

# for linux AMD64
make build-linux
```

## Package to Zip
```bash
make package
```

## Run tests

```bash
make test
```

## Unresolved issues (TODO):
1. Fix build for all platforms (Makefile)
2. Fix local testing (make local-test)
3. Implement deployment part
4. Add documentation for deployment part
5. Document lambda input parameters
