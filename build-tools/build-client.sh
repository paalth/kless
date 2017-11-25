#!/bin/bash

go install cmd/kless/kless.go

if [[ ! -z "$KLESS_UPLOAD_CLI" ]]; then
  gsutil cp $GOBIN/kless gs://klesscli/$BUILD_ID/bin/linux/kless
  gsutil acl ch -u AllUsers:R gs://klesscli/$BUILD_ID/bin/linux/kless
fi