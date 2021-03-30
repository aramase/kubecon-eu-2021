ARG BASEIMAGE="gcr.io/distroless/static:nonroot-amd64"
FROM $BASEIMAGE

COPY ./_output/csi /bin/

ENTRYPOINT [ "/bin/csi" ]
