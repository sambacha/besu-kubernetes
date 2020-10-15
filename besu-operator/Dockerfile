FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/besu-operator \
    USER_UID=1001 \
    USER_NAME=besu-operator

LABEL name="besu" \
      maintainer="support@freighttrust.com" \
      vendor="Hyperledger" \
      version="0.1.0" \
      release="1" \
      summary="Simplify the Deployment & Management of Blockchain Nodes " \
      description="The besu Operator for Kubernetes (SOK) makes it easy for Network Administrators to deploy and operate Enterprise deployments in a Kubernetes infrastructure. Packaged as a container, it uses the operator pattern to manage besu-specific custom resources, following best practices to manage all the underlying Kubernetes objects for you."

# install operator binary
COPY build/_output/bin/besu-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN mkdir /licenses && /usr/local/bin/user_setup


ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}