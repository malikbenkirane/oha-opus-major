# OHA - Player Data Service

A tiny Go micro‑service that serves player state for a multiplayer‑game demo.
It demonstrates a hexagonal architecture, containerisation, and Argo CD‑driven
GitOps deployment.

**Project Overview**

For this homework assignment we spent roughly four hours creating a simple Go
service and the accompanying deployment artifacts. The repository includes:

- **Go service** – exposes a single endpoint: `GET /update-player-data`.
- **Dockerfile** – builds a container image for the service (included in the repo).
- **Argo CD application** – deploys the service using Kustomize.
- **Kubernetes service** – makes the application reachable locally via `kubectl port-forward`.

In this document we provide a brief overview of what could have been added with
additional time, highlighting the most disappointing shortfall: the absence of
telemetry and logs. You can inspect the `feat/otel` branch, which contains a
prototype of the tools we might have used to add observability to the service,
though the service is deliberately lacking these capabilities.


### Table of Contents

<!-- toc -->
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
- [Architecture](#architecture)
- [Trade-offs due to time limitations](#trade-offs-due-to-time-limitations)
- [Additional documents](#additional-documents)
- [Troubleshooting](#troubleshooting)
<!-- /toc -->

### Prerequisites

- Kubernetes (any cluster with RBAC)
- Argo CD (installed in the cluster) – see [argo-cd.readthedocs.io](https://argo-cd.readthedocs.io/)
- `kubectl` ([k8s.io/docs/tasks/tools](https://kubernetes.io/docs/tasks/tools/))
- `kustomize` ([kustomize.io](https://kustomize.io/))
- `argocd` CLI ([argo-cd/.../cli_installation](https://argo-cd.readthedocs.io/en/stable/cli_installation/))

## Getting started

You’ll need a Kubernetes cluster that already has Argo CD installed, together
with `kubectl`, `kustomize`, and the Argo CD CLI.  Make sure the `KUBECONFIG`
environment variable points at your kubeconfig file and that you are logged in
to Argo CD.  We performed the setup using the core release rather than the
Argo CD UI.

To create the Argo CD application, its project, and the destination namespace,
apply these manifests:

    kubectl apply -f https://raw.githubusercontent.com/malikbenkirane/oha-opus-major/a456690e/argocd/project.yaml
    kubectl apply -f https://raw.githubusercontent.com/malikbenkirane/oha-opus-major/a456690e/argocd/application.yaml

To synchronize the application resources:

    argocd app sync omoha-demo

**Testing the service**

1. **Forward a local port to the service’s port 80**  
   (In this example we use local port 8080.)

   ```
   kubectl port-forward -n omoha-demo svc/demo-player-data 8080:80
   ```

2. **Send a request to the forwarded port**

   ```
   curl localhost:8080/update-player-data
   ```

3. **Expected response**

   ```
   [{"PositionX":0,"PositionY":0,"Orientation":0,"CanFire":false,"Weapons":[{"Name":"butterfly","Damage":0}],"Health":0}]
   ```

## Architecture

    .
    └── argocd

The Argo CD manifests that define the project, application, and its namespace
are located in the repository’s root under the **argocd** folder:

    .
    ├── cmd
    │  └── serve
    ├── internal
    │  ├── adapter
    │  │  ├── player_data_repository
    │  │  │  └── mock
    │  │  └── player_data_server
    │  │     └── http
    │  ├── domain
    │  │  ├── player
    │  │  └── weapon
    │  └── port
    └── service

The Go service is built on a hexagonal architecture.

- Ports are defined as a collection of interfaces that abstract the underlying
  operations.

- Two concrete port implementations are provided:
  - A mock repository for player data.
  - An HTTP server that implements the player‑data server interface.

- The domain layer contains the "business" entities used by the application.

- The service composes the HTTP server instance with the mock player‑repository
  and exposes a Service object.

The latter provides a `Run` method, which the `cmd/serve` package invokes to
perform the final bundling.

    type Service interface {
            Run(ctx context.Context) error
    }

Review the [feat/otel](
https://github.com/malikbenkirane/oha-opus-major/blob/feat/otel/internal/otel/otel.go)
branch to see an additional entity I intended to include in this project’s
scope; however, time constraints prevented its integration. The current setup
also lacks liveness, readiness, and startup probes. There was an idea to add a
sidecar container that could interact with the horizontally scalable deployment
when the replicas become unresponsive, along with several other concepts
related to this approach that we can discuss and potentially adopt as an
exercise.

    .
    └── kustomize
       ├── base
       └── overlays
           └── demo

The `kustomize` directories contain the kustomization resources, so a
`kustomization.yaml` file should be present in the `base` directory and in each
subdirectory of the `overlays` folder.

We only provide a demo overlay here, but to configure base‑resource patches for
specific environments (e.g., dev, prod, staging), we would place them in
directories within the overlays folder.

## Trade-offs due to time limitations

We may have devoted excessive time to building the service, leaving us
insufficient time to cover the essential SRE/DevOps elements of the assignment.

It’s been some time since I wanted to explore Argo CD, and I’m grateful for the
chance to do so!

This assignment prompted me to explore several technologies I hadn’t previously
used, giving me a valuable opportunity to broaden my toolkit and discover new
ways to tackle common software‑engineering challenges.  

I focused on making the project reproducible and thoroughly documented for
future engineers—including my own future self—by tackling essential aspects
such as deployment, observability, tracing, metrics, logging, and
configuration, so anyone reviewing the work can quickly grasp and extend it.

Further details on the time‑related trade‑offs are available in the [trade‑offs.md](
docs/trade-offs.md) document.

## Additional documents

With the aim of making this effort memorable and substantive—not just a quick,
easy task—I’ve written documentation for anyone new to Argo CD. You can find it
in the [argocd.md](docs/argocd.md) file.

Another interesting document covers game design. Its TL;DR highlights topics we
can consider when building our own multiplayer game. You can find it
at [game.md](docs/game.md), where it also describes the model used by the Go
service.


## Troubleshooting

- **Argo CD sync fails** – Ensure the cluster context matches `KUBECONFIG` and that the service account used by Argo CD has permission to create resources in the `omaha-demo` namespace.
- **Port‑forward hangs** – Verify that the `demo-player-data` Service exists (`kubectl get svc -n omoha-demo`). 
