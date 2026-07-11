#!/usr/bin/env python3
from pathlib import Path
import re


ROOT = Path(__file__).resolve().parents[1]
release = (ROOT / ".github/workflows/ghcr-build.yml").read_text()
sync = (ROOT / ".github/workflows/upstream-sync.yml").read_text()
zizmor = (ROOT / ".github/workflows/zizmor.yml").read_text()
dockerfile = (ROOT / "Dockerfile-build").read_text()
build_workflow = (ROOT / ".github/workflows/build.yaml").read_text()

assert "ghcr.io/masonjames/ntfy:latest" not in release
assert "type=raw,value=latest" not in release
assert "gh workflow run ntfy-release.yml" in release
assert "--repo masonjames/platform-infra" in release
assert "permission-actions: write" in release
assert "permission-checks: read" in release
assert "runs-on: [self-hosted, hetzner]" not in release
assert "vars.NTFY_AUTO_DEPLOY_ENABLED == 'true'" in release
assert "TZ=America/New_York" in release
assert "^7-03[0-5][0-9]$" in release
assert "deployment-rollback" in release
assert "service:ntfy|0052cc" in release
assert "issues: write" in release
assert "ghcr.io/masonjames/ntfy@$IMAGE_DIGEST" in release
assert "actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32" in release
assert "aquasecurity/trivy-action@ed142fd0673e97e23eac54620cfb913e5ce36c25" in release
assert "actions/create-github-app-token@bcd2ba49218906704ab6c1aa796996da409d3eb1" in sync
assert "permission-contents: write" in sync
assert "persist-credentials: false" in sync
assert "git merge --no-ff upstream/main" in sync
assert "gh pr create" in sync
assert "zizmorcore/zizmor-action@b1d7e1fb5de872772f31590499237e7cce841e8e" in zizmor
assert "online-audits: false" in zizmor
assert 'org.opencontainers.image.revision="$COMMIT"' in dockerfile
assert 'org.opencontainers.image.source="$SOURCE_URL"' in dockerfile
assert re.search(r"^FROM golang:[^@\n]+@sha256:[0-9a-f]{64} AS builder$", dockerfile, re.MULTILINE)
assert re.search(r"^FROM alpine:[^@\n]+@sha256:[0-9a-f]{64}$", dockerfile, re.MULTILINE)
assert "ADD ./action ./action" in dockerfile
assert "ADD ./template ./template" in dockerfile
assert build_workflow.index("Build the deployment image without publishing") < build_workflow.index("Build all the things")

for workflow in (release, sync, zizmor):
    for match in re.finditer(r"uses:\s+[^\s@]+@([^\s#]+)", workflow):
        ref = match.group(1)
        assert re.fullmatch(r"[0-9a-f]{40}", ref), f"unpinned action ref: {ref}"

print("Mason ntfy automation contracts passed")
