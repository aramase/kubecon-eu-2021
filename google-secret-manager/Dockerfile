ARG BASEIMAGE="gcr.io/distroless/static:nonroot-amd64"
FROM $BASEIMAGE

COPY ./_output/google-secret-manager /bin/

ENTRYPOINT [ "/bin/google-secret-manager" ]
