# Bundled Skills

This directory contains OpenCode skills that are bundled into the Chetter runner image.

## Current Skills

| Skill | Purpose |
|---|---|
| `golang-pro` | Go implementation guidance: concurrency, services, generics, performance, tests. |
| `mysql` | MySQL/InnoDB schema, indexes, query tuning, transactions, and operations. |
| `tidb-sql` | TiDB-specific SQL behavior, compatibility notes, vector/full-text features, and transactions. |

## How Skills Reach Agents

Skills are baked into `ghcr.io/flatout-works/chetter-runner` by `runner/Dockerfile.chetter`:

```dockerfile
COPY tools/skills/ /opt/opencode/.agents/skills/
```

The runner starts OpenCode with `HOME=/opt/opencode` in Docker and Kata modes, so OpenCode discovers bundled skills from:

```text
/opt/opencode/.agents/skills/**/SKILL.md
```

There is no current runtime copy from host `~/.agents/skills` into task workspaces. Changing a bundled skill requires rebuilding and redeploying the runner image.

## Task Skill Hints

Task requests and schedules can include a `skills` list. Chetter stores that list and prepends it to the prompt as a hint:

```text
Requested OpenCode skills: golang-pro, mysql. Use these skills when applicable.
```

The hint does not install a skill. The skill must already exist in the runner image.

## Adding Or Updating Skills

1. Add or edit a directory under `tools/skills/<skill-name>/`.
2. Ensure it contains `SKILL.md` with valid skill frontmatter.
3. Keep large supporting material in `references/` under the skill directory.
4. Rebuild the runner image:

   ```bash
   make docker-build-runner
   ```

5. For production, rebuild and push images:

   ```bash
   REGISTRY=ghcr.io/flatout-works TAG=main ./deploy/build-and-push.sh
   ```

6. Redeploy Chetter runners so new tasks use the updated image.

## Verification

To inspect a built runner image:

```bash
docker run --rm --entrypoint sh ghcr.io/flatout-works/chetter-runner:main \
  -lc 'find /opt/opencode/.agents/skills -maxdepth 2 -name SKILL.md -print'
```

Expected output includes the bundled skill files, for example:

```text
/opt/opencode/.agents/skills/golang-pro/SKILL.md
/opt/opencode/.agents/skills/mysql/SKILL.md
/opt/opencode/.agents/skills/tidb-sql/SKILL.md
```
