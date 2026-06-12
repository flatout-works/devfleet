---
name: devfleet
description: Use Devfleet to submit, track, and manage remote agent tasks: runner health, task status, schedules, and cancellation. Triggers on devfleet-related workflow requests, slash commands (/devfleet-*), task management, and fleet health checks.
---

# Devfleet Remote Development Runner Fleet

Devfleet is a self-hosted MCP server for running autonomous AI development agents. It gives your AI tooling a way to submit software development work to a fleet of containerized runners.

- **MCP endpoint:** `https://devfleet.flatout.works/mcp` (hosted) or your own instance
- **Source repo:** `https://github.com/flatout-works/devfleet`

## Available MCP Tools

All tools are prefixed `devfleet_` and available via the `devfleet` MCP server.

### Fleet Health
| Tool | Purpose |
|---|---|
| `devfleet_runner_health` | Runner fleet health, running/stale tasks, image versions, and latest task event age |
| `devfleet_list_tasks` | List recent tasks, optional status filter |
| `devfleet_list_schedules` | List cron task schedules |

### Task Lifecycle
| Tool | Purpose |
|---|---|
| `devfleet_submit_task` | Submit a new task to the fleet |
| `devfleet_task_status` | Get current status and result for a task |
| `devfleet_task_progress` | Get distilled progress timeline |
| `devfleet_task_events` | Get full event history |
| `devfleet_task_latest_event` | Get most recent event |
| `devfleet_cancel_task` | Cancel a pending or running task |
| `devfleet_clear_queue` | Clear queued task messages (requires confirm) |

### Schedules
| Tool | Purpose |
|---|---|
| `devfleet_schedule_task` | Create/activate a cron schedule |
| `devfleet_update_schedule` | Update a schedule by name |
| `devfleet_run_schedule` | Run a schedule immediately |
| `devfleet_delete_schedule` | Delete a schedule by name |
| `devfleet_sync_schedules` | Load schedules from a YAML directory and upsert them |

### Arcane (Vulnerability Scanning, Optional)
| Tool | Purpose |
|---|---|
| `devfleet_arcane_list_images` | List Docker images in an Arcane environment |
| `devfleet_arcane_image_summary` | Vulnerability summary for a specific image |
| `devfleet_arcane_environment_summary` | Aggregated vulnerability counts across all images |
| `devfleet_arcane_list_vulnerabilities` | Detailed vulnerability list with filtering |
| `devfleet_arcane_scanner_status` | Scanner availability and version |

## Common Workflows

### Check Fleet Status
Use `/devfleet-status` or ask:
```
Tell me about the ongoing tasks
```
The agent will call `devfleet_list_tasks` and `devfleet_runner_health` to show what's running, what's done, what's stale, and what failed.

### Submit a Task
Use `/devfleet-submit` or ask explicitly. When submitting, specify:
- `git_url`: your repository URL
- `git_ref`: usually `main`
- `agent_image`: your runner image, such as `ghcr.io/your-org/devfleet-runner:main`
- `prompt`: clear, scoped instructions
- Optional: `agent`, `provider_id`, `model_id`, `variant_id`, `skills`, `timeout_sec`

### Track a Task
```
Show progress for task task_<id>
Show the latest event for task task_<id>
```

### Diagnose Stale Tasks
A running task is stale in fleet health when `last_event_sec > 600`. Check its events and progress to understand what step it is stuck on. Consider canceling and resubmitting.

### Manage Schedules
Schedules are declarative YAML files. See `schedules/` for sample templates. The workflow:

```
Use devfleet_sync_schedules to sync schedules from schedules/
```

This calls `devfleet_sync_schedules` which reads all YAML files and upserts them.

## Working with Schedules

### Adding a New Schedule

1. Copy an existing sample from `schedules/` as a starting point.
2. Edit it with your repo details and prompt.
3. Sync to Devfleet:
   ```
   Use devfleet_sync_schedules to sync schedules from schedules/
   ```

### Customizing a Schedule

Each schedule YAML supports these fields:

| Field | Required | Description |
|---|---|---|
| `name` | yes | Unique schedule name (slug) |
| `enabled` | no | `true` to activate, `false` to pause (default false) |
| `cron_expr` | yes | Five-field cron or `@hourly`, `@daily` |
| `prompt` | yes | Task prompt run on each cron fire |
| `git_url` | yes | Repository URL to clone |
| `git_ref` | no | Branch/tag/commit (default main) |
| `agent_image` | no | Runner image override |
| `agent` | no | OpenCode agent to use |
| `provider_id` | no | LLM provider for model selection |
| `model_id` | no | Model ID |
| `variant_id` | no | Model variant |
| `skills` | no | List of skill names to load |
| `timeout_sec` | no | Task timeout in seconds |

### Tweaking an Existing Schedule

To change a schedule's cron expression, prompt, model, or other fields:

1. Edit the `schedules/*.yaml` file directly.
2. Sync again:
   ```
   Use devfleet_sync_schedules to sync schedules from schedules/
   ```
   Or update individually:
   ```
   Use devfleet_update_schedule to change nightly-changelog-update's model to opencode/minimax-m3
   ```

### Pausing a Schedule
Set `enabled: false` in the YAML and sync, or:
```
Use devfleet_update_schedule to disable nightly-issue-fixer
```

### Running a Schedule Manually
```
Use devfleet_run_schedule to run the nightly-changelog-update schedule now
```

### Deleting a Schedule
```
Use devfleet_delete_schedule to delete nightly-docs-update
```
Remove the corresponding YAML file so `devfleet_sync_schedules` doesn't recreate it.

### Keeping Schedules in Your Repo

The recommended pattern is to store schedule YAMLs in your own repo (not in devfleet's `schedules/` directory). When you set up your project:

1. Create a `schedules/` directory in your project repo.
2. Copy the samples from devfleet's `schedules/` as starting points.
3. Customize for your project (repo URL, agent image, prompt details).
4. Sync with Devfleet pointing at your project's `schedules/` directory.

This way your schedules are version-controlled alongside your code and can be reviewed in PRs.

## Safety Rules

- Never send secrets (API keys, tokens, passwords) in task prompts or env vars
- Task prompts must explicitly state whether file edits and PR creation are allowed
- Tell tasks to create branches and PRs rather than pushing to the default branch
- Use `timeout_sec` appropriate for the work (e.g., 600 for quick checks, 3600 for code changes)
- Devfleet clones from Git; tasks cannot access uncommitted local changes
- For recurring schedules, check `schedules/` YAMLs into version control

## Model Selection

Common model choices for different task types:

| Task type | Suggested model | Notes |
|---|---|---|
| Docs/changelog | synthetic/kimi-k2.6 | Fast, cheap, good at summarizing |
| PR reviews | opencode-go/minimax-m3 | Thorough, structured output |
| Code quality / fixes | opencode/deepseek-v4-pro | Good at Go code analysis |
| Bugfixes | deepseek/deepseek-chat | Budget-friendly for simple fixes |

Prefer reliable, cost-effective models for scheduled maintenance tasks. Reserve large/expensive models for complex implementation work.
