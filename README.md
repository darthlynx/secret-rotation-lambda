# secret-rotation-lambda
AWS Lambda which rotates SecretManager secrets

## Build
```bash
# for your local host machine
make build
```

## Package to Zip
```bash
make package
```

## Run unit tests

```bash
make test
```

## Run lambda simulation locally

Prerequisites:
- Docker desktop
- AWS credentials

```bash
export AWS_ACCESS_KEY_ID=your_key_id
export AWS_SECRET_ACCESS_KEY=your_secret_key
# Optional, if using temporary credentials:
export AWS_SESSION_TOKEN=your_token

make local-build
```

## Unresolved issues (TODO):
1. Fix build for all platforms (Makefile)
2. Implement deployment part
3. Add documentation for deployment part
4. Document lambda input parameters
