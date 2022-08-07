FROM gcr.io/distroless/static:nonroot AS final
COPY ./bin/zts-upgrade-handler /app
USER 65532:65532
#EXPOSE 8080:8080
ENTRYPOINT ["/app"]