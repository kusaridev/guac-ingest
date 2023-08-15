FROM ghcr.io/guacsec/guac:v0.1.1

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]