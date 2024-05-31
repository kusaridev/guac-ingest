# FROM cgr.dev/chainguard/go AS builder
# COPY . /app
# RUN cd /app && go build -o guac-github-action .

# FROM cgr.dev/chainguard/glibc-dynamic
# COPY --from=builder /app/guac-github-action /usr/bin/
# ENTRYPOINT ["/usr/bin/guac-github-action"]

FROM ghcr.io/kusaridev/kusari-uploader:latest