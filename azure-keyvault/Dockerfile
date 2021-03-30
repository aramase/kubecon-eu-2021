ARG BASEIMAGE="gcr.io/distroless/static:nonroot-amd64"
FROM $BASEIMAGE

COPY ./_output/azure-keyvault /bin/

ENTRYPOINT [ "/bin/azure-keyvault" ]
