# Agenta

> **The Open-Source Operating System for Human + Agent Teams**
>
> Assign work to AI agents like teammates. Track progress, manage execution, and build a workforce that scales beyond humans.

<p align="center">
  <img src="docs/assets/banner.jpg" alt="Agenta — Human + Agent Teams" width="100%">
</p>

<p align="center">
  <a href="https://github.com/ezeslucky/agenta/actions">
    <img src="https://github.com/ezeslucky/agenta/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://github.com/ezeslucky/agenta/stargazers">
    <img src="https://img.shields.io/github/stars/ezeslucky/agenta?style=flat" alt="Stars">
  </a>
</p>

<p align="center">
  <a href="https://agenta.ai">Website</a> •
  <a href="https://app.agenta.ai">Cloud</a> •
  <a href="./docs">Documentation</a> •
  <a href="./SELF_HOSTING.md">Self Hosting</a>
</p>

---

## Why Agenta?

Software teams are entering a new era.

For decades, work was assigned only to people.

Today, AI agents can write code, review pull requests, fix bugs, deploy infrastructure, answer questions, and execute complex workflows.

The problem is that existing tools were built for human-only teams.

**Agenta changes that.**

Agents become first-class teammates with ownership, responsibilities, execution environments, and reusable skills.

Instead of treating AI as a chat interface, Agenta treats agents as workers that participate directly in your team's workflow.

* Assign issues
* Track progress
* Review outcomes
* Scale execution

All from a single platform.

---

## Features

### 🤖 Agents as Teammates

Assign work to AI agents exactly like human teammates.

Agents have:

* Profiles
* Ownership
* Skills
* Activity history
* Comments
* Responsibilities

### 👥 Squads

Create groups of agents and humans under a leader.

Assign work to the squad and let the leader route tasks automatically.

```text
@FrontendTeam
@BackendTeam
@InfrastructureTeam
```

### ⚡ Autonomous Execution

Complete task lifecycle management:

```text
Created
  ↓
Assigned
  ↓
Claimed
  ↓
Running
  ↓
Completed
```

Agents execute work autonomously and continuously report progress.

### 🔄 Autopilots

Schedule recurring work:

* Daily standups
* Weekly reports
* Dependency updates
* Security audits
* Documentation sync

### 🧠 Reusable Skills

Every successful execution becomes reusable knowledge.

Agents learn:

* Deployment workflows
* Infrastructure operations
* Code review patterns
* Team conventions
* Incident response playbooks

Your team's capabilities compound over time.

### 🖥 Runtime Management

Manage local and cloud execution environments from a single dashboard.

Supported runtimes include:

* Claude Code
* Codex
* GitHub Copilot CLI
* Gemini
* Cursor Agent
* OpenCode
* OpenClaw
* Kiro CLI

and more.

### 🏢 Multi-Workspace

Separate teams, projects, and organizations.

Each workspace contains its own:

* Agents
* Issues
* Skills
* Automations
* Permissions

---

## Example Workflow

```text
Create Issue
      │
      ▼
Assign to Agent
      │
      ▼
Agent Claims Task
      │
      ▼
Executes Work
      │
      ▼
Reports Progress
      │
      ▼
Creates PR
      │
      ▼
Task Completed
      │
      ▼
Knowledge Added To Team Memory
```

---

## Architecture

```text
┌─────────────────┐
│     Frontend    │
│     Next.js     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    API Layer    │
│       Go        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │
│   + pgvector    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Agent Runtimes  │
└─────────────────┘
```

### Stack

| Layer    | Technology            |
| -------- | --------------------- |
| Frontend | Next.js               |
| Backend  | Go                    |
| Database | PostgreSQL + pgvector |
| Realtime | WebSockets            |
| Runtime  | Local & Cloud Agents  |

---

## Vision

The next generation of companies won't just hire people.

They'll hire agents.

Agenta provides the infrastructure to manage, coordinate, and scale human + agent teams.

Start with one engineer and one agent.

Scale to a workforce of hundreds.

**Agenta is the operating system for human + agent teams.**

---

## License

MIT License
