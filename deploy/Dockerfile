FROM quay.io/openshift/origin-operator-registry:latest

COPY deploy/olm-catalog /registry

# Initialize the database
RUN initializer --manifests /registry/observatorium-operator --output bundles.db

ENTRYPOINT ["registry-server"]
CMD ["--database", "bundles.db"]
