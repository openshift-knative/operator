FROM openshift/origin-release:golang-1.14 AS builder
WORKDIR ${GOPATH}/src/knative.dev/operator
COPY . .
ENV GOFLAGS="-mod=vendor"
RUN go build -o /tmp/operator ./cmd/operator
RUN cp -Lr ${GOPATH}/src/knative.dev/operator/cmd/operator/kodata /tmp

FROM openshift/origin-base
COPY --from=builder /tmp/operator /ko-app/operator
COPY --from=builder /tmp/kodata/ /var/run/ko
ENV KO_DATA_PATH="/var/run/ko"
LABEL \
    com.redhat.component="openshift-serverless-1-tech-preview-knative-rhel8-operator-container" \
    name="openshift-serverless-1-tech-preview/knative-rhel8-operator" \
    version="v0.15.0" \
    summary="Red Hat OpenShift Serverless 1 Knative Operator" \
    maintainer="serverless-support@redhat.com" \
    description="Red Hat OpenShift Serverless 1 Knative Operator" \
    io.k8s.display-name="Red Hat OpenShift Serverless 1 Knative Operator"

ENTRYPOINT ["/ko-app/operator"] 