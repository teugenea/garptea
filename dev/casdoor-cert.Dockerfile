FROM casbin/casdoor:v1.595.0
USER root
COPY dev-root-CA.crt /usr/local/share/ca-certificates
RUN update-ca-certificates