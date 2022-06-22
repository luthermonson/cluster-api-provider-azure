#!/bin/bash

rm generated/*
mkdir -p generated
kustomize build --load-restrictor LoadRestrictionsNone global > ./generated/core-global.yaml
kustomize build --load-restrictor LoadRestrictionsNone base > ./generated/core-base.yaml