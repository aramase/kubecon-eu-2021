FROM gcr.io/distroless/static:nonroot
ARG TARGETARCH

COPY ./_output/${TARGETARCH}/csi /bin/

ENTRYPOINT [ "/bin/csi" ]
