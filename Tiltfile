# Define a Kubernetes resource for Service 1
k8s_yaml(['kirzop/k8s/deployment.yaml', 'kirzop/k8s/service.yaml', 'kirzop/k8s/pvc.yaml'])

docker_build(
    'kirzop',
    './kirzop/',
)

k8s_resource(
    'kirzop',
    port_forwards=5000,
)

# PHP frontend
k8s_yaml(kustomize('kirplop/kubernetes'))

docker_build(
    'kirplop',
    './kirplop/',
)

k8s_resource(
    'kirplop',
    port_forwards=5000,
)

# The Kubernetes go operator
docker_build(
    'controller',
    './kirop',
)

k8s_yaml(kustomize('kirop/config/manager'))
k8s_yaml(kustomize('kirop/config/default'))
