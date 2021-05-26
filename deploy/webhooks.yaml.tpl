apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: k8s-webhook-pull-policy-webhook
  labels:
    app: k8s-webhook-pull-policy-webhook
    kind: mutator
webhooks:
  - name: image-pull-policy-webhook
    admissionReviewVersions: ["v1"]
    sideEffects: None
    clientConfig:
      service:
        name: k8s-webhook-pull-policy
        namespace: k8s-webhook-pull-policy
        path: /wh/mutating/mark
      caBundle: CA_BUNDLE
    rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: ["*"]
        apiVersions: ["*"]
        resources: ["deployments", "daemonsets", "cronjobs", "jobs", "statefulsets", "pods"]