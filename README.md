# sentry-krm-extensions

Kubernetes kinds for the Sentry configuration that the `jianyuan/sentry`
Terraform provider manages and no maintained Kubernetes operator does. Sentry
publishes no operator; the only two community ones
([jace-ys/sentry-operator](https://github.com/jace-ys/sentry-operator) and
[sr/kube-sentry-controller](https://github.com/sr/kube-sentry-controller)) were
last updated in 2021 and 2023 and cover teams, projects and client keys only.

By the numbers: `jianyuan/sentry` (officially sponsored by Sentry; community tier on
the registry) is at version 0.15.8 with 24 resources and
29,510,829 downloads. Operators cover 0. This repository adds
**21 kinds** standing in for **24 of 24 (100%)** of them (10 generated from Sentry's
OpenAPI spec covering 11 resources, plus 11 hand-written descriptors covering 13 resources).

**Notes.** Every Sentry path ends in a slash; requests keep it.
`Detector` stands in for both `sentry_uptime_monitor` and `sentry_metric_monitor`, and a
field that is invalid for the detector's `type` fails at the API. `MetricAlert` stands in
for both `sentry_metric_alert` and `sentry_alert`. `ProjectSpikeProtection` stands in for
both `sentry_project_spike_protection` and `sentry_all_projects_spike_protection`
(set `projects: ["__ALL_PROJECTS__"]` for org-wide spike protection).

**Readiness:** **9 kinds** are proven `live` against `sentry.io` (`TestEveryKindLive`);
the remaining **12 kinds** are `fake` (proven against the fake API built from `types.yaml`
via `TestEveryKindOnFake`, skipped in `hack/live/live.yaml` where the test organization
lacks integrations/member scopes, uses newer Detectors/Workflows instead of legacy endpoints,
or where read-back uses a different shape).

## How it works

The kinds are generated, not hand-written. `hack/gen/vendor.yaml` pins
Sentry's OpenAPI spec ([getsentry/sentry-api-schema](https://github.com/getsentry/sentry-api-schema),
FSL-1.1-Apache-2.0; fetched at build time, never committed here);
[krm-rest-runtime](https://github.com/highfivelabs/krm-rest-runtime)'s generator finds
every resource the API exposes in a create/read/update/delete shape, keeps those that
stand in for a Terraform resource in the gap list (`hack/gen/terraform-gaps.json`), and writes:

- `types.yaml`: per kind, its API calls, path parameters and Terraform equivalents
- `config/crds/`: one CRD per kind, its spec the API's request body
- `zz_generated.report.md`: the full list

One generic manager (`cmd/manager`) runs every kind on krm-rest-runtime's
reconciler: it creates the resource, reads it back every 10 minutes and
puts back anything changed outside Kubernetes, and deletes it with the
object. Spec fields keep Sentry's own names, so the API reference
applies field for field. Path parameters are spec fields too.

Hand corrections to the automatic matching are in `vendor.yaml`'s `overrides`:
`creates` says which POST makes a team or project (Sentry creates them on a
different path than they live at), and `map` points the monitors at their endpoints.

## Install

```sh
kubectl kustomize config/default | KO_DOCKER_REPO=<registry> ko resolve -f - | kubectl apply --server-side -f -
```

Then, in each namespace that manages Sentry, a Secret with `token` (an internal
integration or organization token) and a Connection named `default`
(`config/samples/connection.yaml`). Sentry hosts per region: set `baseURL` to
`https://us.sentry.io/api/0` (default), `https://de.sentry.io/api/0`, or your
self-hosted `https://<host>/api/0`. Objects use the Connection unless
`spec.connectionRef.name` names another. Secret values go in a Kubernetes Secret,
not the spec: `spec.secretRefs: {<field>: {name, key}}` fills that body field.

To run only some API groups, install those groups' CRDs from `config/crds`
and pass `--groups` to the manager.

## Kinds

| Kind | Terraform equivalent | Shape | Read |
| --- | --- | --- | --- |
| `cron.sentry.hfl.dev/v1alpha1` `CronMonitor` | `sentry_cron_monitor` | crud | `GET /organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/` |
| `dashboard.sentry.hfl.dev/v1alpha1` `Dashboard` | `sentry_dashboard` | crud | `GET /organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/` |
| `key.sentry.hfl.dev/v1alpha1` `Key` | `sentry_key` | crud | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/keys/{key_id}/` |
| `detector.sentry.hfl.dev/v1alpha1` `Detector` | `sentry_metric_monitor`, `sentry_uptime_monitor` | crud | `GET /organizations/{organization_id_or_slug}/detectors/{detector_id}/` |
| `notification.sentry.hfl.dev/v1alpha1` `NotificationAction` | `sentry_notification_action` | crud | `GET /organizations/{organization_id_or_slug}/notifications/actions/{action_id}/` |
| `organization.sentry.hfl.dev/v1alpha1` `Organization` | `sentry_organization` | settings | `GET /organizations/{organization_id_or_slug}/` |
| `organization.sentry.hfl.dev/v1alpha1` `OrganizationMember` | `sentry_organization_member` | crud | `GET /organizations/{organization_id_or_slug}/members/{member_id}/` |
| `project.sentry.hfl.dev/v1alpha1` `Project` | `sentry_project` | crud | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/` |
| `project.sentry.hfl.dev/v1alpha1` `ProjectOwnership` | `sentry_project_ownership` | settings | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/ownership/` |
| `team.sentry.hfl.dev/v1alpha1` `Team` | `sentry_team` | crud | `GET /teams/{organization_id_or_slug}/{team_id_or_slug}/` |
| `alert.sentry.hfl.dev/v1alpha1` `IssueAlert` | `sentry_issue_alert` | crud | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/rules/{rule_id}/` |
| `alert.sentry.hfl.dev/v1alpha1` `MetricAlert` | `sentry_metric_alert`, `sentry_alert` | crud | `GET /organizations/{organization_id_or_slug}/alert-rules/{alert_rule_id}/` |
| `plugin.sentry.hfl.dev/v1alpha1` `Plugin` | `sentry_plugin` | upsert | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/plugins/{plugin_id}/` |
| `organization.sentry.hfl.dev/v1alpha1` `OrganizationRepository` | `sentry_organization_repository` | crud | `GET /organizations/{organization_id_or_slug}/repos/{repository_id}/` |
| `organization.sentry.hfl.dev/v1alpha1` `OrganizationCodeMapping` | `sentry_organization_code_mapping` | crud | `GET /organizations/{organization_id_or_slug}/code-mappings/{code_mapping_id}/` |
| `project.sentry.hfl.dev/v1alpha1` `ProjectSymbolSource` | `sentry_project_symbol_source` | crud | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/symbol-sources/{symbol_source_id}/` |
| `team.sentry.hfl.dev/v1alpha1` `TeamMember` | `sentry_team_member` | upsert | `GET /organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/` |
| `integration.sentry.hfl.dev/v1alpha1` `IntegrationOpsgenie` | `sentry_integration_opsgenie` | settings | `GET /organizations/{organization_id_or_slug}/integrations/{integration_id}/` |
| `integration.sentry.hfl.dev/v1alpha1` `IntegrationPagerduty` | `sentry_integration_pagerduty` | settings | `GET /organizations/{organization_id_or_slug}/integrations/{integration_id}/` |
| `project.sentry.hfl.dev/v1alpha1` `ProjectInboundDataFilter` | `sentry_project_inbound_data_filter` | settings | `GET /projects/{organization_id_or_slug}/{project_id_or_slug}/filters/` |
| `project.sentry.hfl.dev/v1alpha1` `ProjectSpikeProtection` | `sentry_project_spike_protection`, `sentry_all_projects_spike_protection` | assoc | `GET /organizations/{organization_id_or_slug}/projects/` |

## Not covered

<!-- BEGIN generated: not covered (krmgen.py, from vendor.yaml hardBlocks and draftNotes) -->
Every one of the 24 Terraform gap types has a kind.
<!-- END generated: not covered -->

## For agents (MCP)

`cmd/mcp` serves this repository's kinds over the Model Context Protocol
(stdio), from `types.yaml`, `config/crds` and the live tests' records.
`.mcp.json` registers it, so an agent session opened in this repository
has it without setup. Elsewhere, run it from the repository root:

```sh
go run ./cmd/mcp
```

Its tools: `list_kinds`, `find_kinds_for_terraform`, `describe_kind` (every
field, the REST calls, the live record), `example` (the live-tested spec
where there is one), `validate_manifest` (the CRD schema, as the API server
checks it) and `plan_on_fake` (the controller run against the fake API,
with every REST call it makes). Nothing it does reaches the real API.
krm-rest-runtime's README ("For agents") has the details.

## Regenerate

```sh
hack/gen/generate.sh   # needs python3 + pyyaml, and krm-rest-runtime beside this repo
go test ./...
```
