# Docker file for AWS Lambda Go runtime
FROM public.ecr.aws/lambda/go:1

COPY bin/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap

CMD [ "bootstrap" ]
