# Concourse gitlab resource

## Usage

In a pipeline:

### Resource type definition

```yaml
- name: "gitlab"
  type: registry-image
  source:
    repository: <resource/repository>
    tag: latest
```

### Merge request status

This allows you to get a trigger when a merge request change to a specific state.

#### Source

```yaml
resources:
- name: merge_request_status
  type: gitlab
  source:
    # base url to the gitlab server
    server_url: "https://server_url.com"
    # gitlab project id
    project_id: "<project_id>"
    # gitlab token, pat, group or project level with at least
    # api, read api, write api permissions
    token: "<token>"
    # feature to use
    feature: "merge_request_status"
    merge_request_status:
      # <required> merge_request_iid to watch
      merge_request_iid: 1
      # <required> filter the merge request status that will trigger
      # Any of opened, closed, locked or merged.
      state: ["closed", "merged"]
```


#### Get

```yaml
get: merge_request_status
trigger: true # || false
```

#### Put

No action implemented.

### Deployments

#### Source

```yaml
resources:
- name: deployments
  type: gitlab
  source:
    # base url to the gitlab server
    server_url: "https://server_url.com"
    # gitlab project id
    project_id: "<project_id>"
    # gitlab token, pat, group or project level with at least
    # api, read api, write api permissions
    token: "<token>"
    # feature to use, in this case deployments
    feature: "deployments"
    # <optional> filter the deployment status 
    # One of created, running, success, failed, canceled, or blocked.
    status: "running"
    # <optional> filter by the environment name
    environment: "($ $env $)"
    # <optional> configure the way the resource fetch deployments
    # either value are:
    # - default: the resource will behave like a classic concourse resource
    #     it will pull the latest version and list all newer version at once
    # - yield: this mode is more adapted if you want to ensure to process all 
    #     deployments in the order they were sent. When set to yield, the resource
    #     will always return the older deployment.
    #     Mainly ment to be used to track running deployments.
    mode: "yield"
```

#### Get

```yaml
get: deployments
trigger: true # || false
```

#### Put

##### Create a deployment

> [!Important]
> the `environment` attribute in the `source` must
> be specified to create a deployment.

```yaml
put: deployments
params:
  # the resource need the path to its metadata directory, that matches the resource name
  metadata_dir: deployments
  # Specify the action
  action: create
  # git sha
  sha: <git_sha>
  # branch or tag
  ref: <git_ref>
  # indicate if the ref is a git tag
  tag: true | false
  # Specify the status of the deployment
  # One of running, success, failed, or canceled
  status: running
```

#### Update the deployment

```yaml
put: deployments
params:
  # the resource need the path to its metadata directory, that matches the resource name
  metadata_dir: deployments
  # Specify the new status of the deployment
  # One of running, success, failed, or canceled
  status: running
```

### Delete a deployment

```yaml
put: deployments
params:
  # the resource need the path to its metadata directory, that matches the resource name
  metadata_dir: deployments
```

### Environments

Source configuration.

```yaml
resources:
- name: environment
  type: gitlab
  source:
    project_id: "<gitlab_project_id>"
    token: "<gitlab_token>"
    feature: "environments"
    environments:
      # name of the environment for check steps
      name: "<environment>"
```

#### Create or update the environment

```yaml
put: "environments"
params:
  action: "upsert"
  # the resource need the path to its metadata directory, that matches the resource name
  metadata_dir: environments
  # environment name
  name: "<environment_name>"
  # <optional> description displayed on the environment
  description: |
    <nice_descriptive_text>
  # <optional> attach an external url to the environment, could be a link to this component
  external_url: "($ .console_url $)/organizations/($ .organization $)/projects/($ .project $)/environments/($ .environment $)/components/($ .component $)/overview"
  # <optional> gitlab environment tier
  # Allowed values are production, staging, testing, development, and other.
  tier: <environment_tier>
```

#### Stop an environment

```yaml
put: "environment"
params:
  action: "stop"
```

#### Delete an environment

```yaml
put: "environment"
params:
  action: "delete"
```

#### Delete all stopped environments

```yaml
put: "environment"
params:
  action: "delete_stopped"
  # <optional> The date before which environments can be deleted.
  # Defaults to 30 days ago. Expected in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ).
  before: "<timestamp>"
  # <optional> Maximum number of environments to delete. Defaults to 100.
  limit: "<timestamp>"
  # Defaults to true for safety reasons.
  # It performs a dry run where no actual deletion is performed. Set to false to
  # actually delete the environment.
  dry_run: false
```

#### Stop stale environments

```yaml
put: "environment"
params:
  action: "stop_stale"
  # <optional> The date before which environments can be stopped.
  # Defaults to 30 days ago. Expected in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ).
  before: "<timestamp>"
```

## Development

This a [concourse resource](https://concourse-ci.org/resources.html) meant to interact
with Gitlab CI.

Before reading this readme, please read the links below for context:
- https://concourse-ci.org/resources.html
- https://concourse-ci.org/implementing-resource-types.html

### Quickstart

The resource is made in go, all actions like build/test/etc have been added to a
[justfile](https://github.com/casey/just).

Below is a list of the main commands

#### Build

```bash
go build
```

#### Build docker image

```bash
docker build . --tag <tag>
```

#### Tests

```bash
go test ./...
```

#### Lint

```bash
just lint
```

### TODO

- [ ] move all deployments source to a `deployment` object like environments
- [ ] add e2e

### Implementation

#### Entrypoint

A concourse resource must contain three scripts:
```
/opt/resource/check
/opt/resource/in
/opt/resource/out
```

To integrate in the docker image, the binary will be copied in `/opt/resource/gitlab-resource`,
to integrate with concourse, we created three scripts [here](./scripts/) named
`check`, `in`, `out` that contains a shebang like `#! /opt/resource/gitlab-resource`.

This will call gitlab resource like that:
```
gitlab-resource /opt/resource/<check,in,out> <args>
```

The [main](./main.go) contains the logic to handle that an trigger the action.

This trick allow us to make a very small docker image with only the binary and
some utility files, making the attack surface minimal.

#### Features

Each feature will have its own sets of params and settings, they must be added
in the [models folder](./models/).

Then you will have to create a handler that implements the [Feature interface](./features/interface.go)
and add its creation in the NewFeatureHandler function.

For the implementation details, look up how the [deployments](./features/deployments/)
or [environments](./features/environments/) are made.

### Tests

There are no end to end tests for now (as it would require way more time to implement
right now).

Some units tests have been itegrated using the defaults go test system.
