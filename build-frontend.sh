#!/usr/bin/env bash
set -xeuo pipefail
cd frontend || exit 1
npm ci
npm run build
