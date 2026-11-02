For trying manifests, helm charts and Kubernetes behaviour, agents may create, use and delete a
local minikube cluster (`minikube start --driver=docker`, kube context `minikube`) without
asking, as long as everything stays local. Never point such tests at a shared cluster.

Always pass `--context minikube` to kubectl and helm, so a command can never hit another cluster.
`minikube start` switches the current context to minikube: capture `kubectl config current-context`
before it and restore it afterwards with `kubectl config use-context <old>`. Switching the user's
current context is not allowed. Run `minikube delete` when done unless the user wants to keep the
cluster. Anything touching shared clusters still follows research-permissions.md.

Why: the user's current kube context decides where bare kubectl commands go, and a silent switch
to a local cluster, or a forgotten switch back, sends later commands to the wrong place.
