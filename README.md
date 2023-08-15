# Overview
Github Action for ingesting SBOMs and Attestations into [GUAC](https://github.com/guacsec/guac) as part of your github workflow.
This will enable quick and easy integration to your GUAC instance with very minimal input.

Currently this wraps the `guacone collect files` command.

## Inputs
### `files` 
**Required** Path to directory or specific file to ingest

### `gql-addr`  
**Required** GUAC GraphQL API Endpoint - example:  http://localhost:8080/query

## Outputs

### `console_out`
Raw output of the guacone command

## Example Usage

```- name: GUAC ingestion
    uses: actions/guacone@v1
    id: guacone
    with:
        files: './spdx.json'
        gql-addr: 'http://localhost:1234/query'
```