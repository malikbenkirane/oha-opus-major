# About ArgoCD

We won’t cover the deployment details of ArgoCD in your cluster; instead, we’ll
provide a concise, introductory overview of project management and ArgoCD
application management.

Comprehensive details on Argo CD operations are provided in the [operator manual](
https://github.com/argoproj/argo-cd/tree/master/docs/operator-manual
). If you prefer a declarative setup with Kubernetes manifests, refer to
[this document](
https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/declarative-setup.md
).

Target the Argo CD namespace and ensure `KUBECONFIG` points to the correct
cluster context:

    kubectl config set-context --current --namespace argocd

For this demo we use the core server; log in before issuing further Argo CD
commands:

    argocd login --core

## Project management

You can see all projects that are defined by running:

    argocd proj list

And you can retrieve the full `AppProject` resource definitions with:

    kubectl get appproject -o yaml

Take a look at [project.yaml](
https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/project.yaml)
in argocd operator manual for a complete project‑specification example.

For this demo, however, we’ll create a minimal project using the CLI:

    arogcd proj create omoha-demo

Keeping this file under CVS version control is recommended, so we'll store it
in the repository’s top‑level `argocd` directory.

We also configure authorized destinations and source repositories.
Refer to the `argocd/project.yaml` file in the repository’s source.

## Application management

To list every application that has been defined, run:

    argocd app list

You can also use this command to verify the health and current state of each
application.

For the declartive setup checkout argocd operator manual [application.yaml](
https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/application.yaml
).

For this assignment we’ll work with a single overlay called **demo**.  
- The **main** branch will be referenced via `spec.source.revision`.
- The overlay files are located in `kustomize/overlays/demo`, and the
  `kustomization.yaml` resides in the `kustomwhize` directory at the
  repository’s root, so we’ll configure `spec.source.path` to point there.

If additional versions of the application were required, we would introduce
more manifest files and revisit the naming conventions and directory layout
used for the declarative manifests under the `argocd` root. Since we only have
one application right now, the configuration lives in a single file:
`argocd/application.yaml`.

Refer to the `argocd/application.yaml` file in the repository’s source.
