# GUAC action

This Action ingests SBOMs and Attestations into
[GUAC](https://github.com/guacsec/guac) as part of your github
workflow. Authentication is provided by GitHub's OIDC provider and a specified
auth token provider. This will enable quick and easy integration to your GUAC
instance with very minimal input.

For details on how to query and utilize the data upon ingestion, please see documentataion for the [GUAC use cases](https://docs.guac.sh/guac-use-cases/). 

This action will only work with an OAuth2 protected GUAC GraphQL API endpoint - e.g. [The Kusari hosted GUAC platform](https://www.kusari.dev/).  

## Usage

See [action.yaml](action.yaml)

```yaml
steps:
  - uses: actions/checkout@v3

  - uses: [Your build and SBOM/Provenance generation steps]

  - uses: kusaridev/guac-ingest@v0
    name: GUAC ingestion
    with:
      files: './spdx.json'
      gql-addr: 'https://[tenant-id].api.kusari.cloud/query'
      token-url: ${{ secrets.TOKENURL }}
      client-id: ${{ secrets.CLIENTID }}
```

## Inputs

### `files`

**Required** Path to directory or specific file to ingest

### `gql-addr`

**Required** GUAC GraphQL API Endpoint - example:  https://guac.instance/query

### `token-url`

**Required** URL of auth token provider - example:  https://token.provider/oauth2/token

### `gql-addr`

**Required** Client ID for auth token provider - example:  abcd-efgh-1234...

## Outputs

### `console_out`

Raw output of the guacone command

# License

The scripts and documentation in this project are released under the [Apache License](LICENSE)