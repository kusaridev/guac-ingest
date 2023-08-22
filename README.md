# Overview

Github Action for ingesting SBOMs and Attestations into
[GUAC](https://github.com/guacsec/guac) as part of your github
workflow. Authentication is provided by GitHub's OIDC provider and a specified
auth token provider. This will enable quick and easy integration to your GUAC
instance with very minimal input.

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

## Example Usage

```yaml
    - name: GUAC ingestion
      uses: kusaridev/guac-github-action@v1
      with:
        files: './spdx.json'
        gql-addr: 'https://guac.kusari.dev/query'
        token-url: ${{ secrets.TOKENURL }}
        client-id: ${{ secrets.CLIENTID }}
```
