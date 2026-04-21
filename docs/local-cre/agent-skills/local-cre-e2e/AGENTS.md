# Local CRE E2E Agent Notes

This directory holds a reusable agent skill for running Local CRE and CRE e2e tests.

## Scope

These instructions apply only to files under `docs/local-cre/agent-skills/local-cre-e2e/`.

## Intent

- Keep the skill practical and workflow-oriented.
- Favor commands that work from a clean checkout.
- Prefer the default topology unless the user explicitly needs a different topology or a custom override.
- Keep topology customization guidance minimal and repo-specific.

## When Updating The Skill

- Keep the skill aligned with:
  - `docs/local-cre/`
  - `core/scripts/cre/environment/configs/`
- If test commands or topology selection rules change, update this skill in the same PR when possible.
- Do not add generic AI-agent boilerplate; keep it specific to Local CRE usage in this repo.
