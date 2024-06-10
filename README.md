# guac-ingest Action

This Action ingests SBOMs and Attestations into the [Kusari hosted GUAC platform](https://www.kusari.dev/) as part of your github workflow. This will enable quick and easy integration to your GUAC instance with very minimal input.

Authentication credentials (client-id, client-secret) are provided by the Kusari team. 

For details on how to query and utilize the data upon ingestion, please see documentataion for the [GUAC use cases](https://docs.guac.sh/guac-use-cases/). 


## Usage

See [action.yaml](action.yaml)

```yaml
steps:
  - uses: actions/checkout@v4

  - uses: [Your build and SBOM/Provenance generation steps]

  - uses: kusaridev/guac-ingest@v0
    name: GUAC Ingestion
    with:
      files: './spdx.json'
      api-addr: 'https://[kusari-tenant-id].api.us.kusari.cloud'
      client-id: ${{ secrets.KUSARI_CLIENT_ID }}
      client-secret: ${{ secrets.KUSARI_CLIENT_SECRET }}
```

## Inputs

### `files`

**Required** - Path to directory or specific file to ingest

### `api-addr`

**Required** - Kusari hosted GUAC tenant api endpoint

### `client-id`

**Required** - Client id for auth token provider

### `client-secret`

**Required** - Client secret for auth token provider

### `token-url`

Url for auth token provider

## Outputs

### `console_out`

Raw output of the kusari-uploader command

# License

The scripts and documentation in this project are released under the [Apache License](LICENSE)