//
// Copyright 2023 Kusari, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/guacsec/guac/pkg/assembler/clients/helpers"
	"github.com/guacsec/guac/pkg/handler/collector"
	"github.com/guacsec/guac/pkg/handler/collector/file"
	"github.com/guacsec/guac/pkg/handler/processor"
	"github.com/guacsec/guac/pkg/handler/processor/process"
	"github.com/guacsec/guac/pkg/ingestor/parser"
	"github.com/guacsec/guac/pkg/logging"
	"github.com/sigstore/cosign/v2/pkg/providers"
	_ "github.com/sigstore/cosign/v2/pkg/providers/github"
	"golang.org/x/oauth2"
)

// This program grabs an OIDC ID token from GitHub and then exhcanges it for an
// Auth token from the passed in URL. In case a client secret is provided, it'd
// use that and perform client_credentials flow instead.
// Then it uses GUAC collectors to process files and push them to a GUAC gql
// server. It is only meant to be run by a GitHub Action (see action.yaml)
// and therefore has limited options and argument processing.
func main() {
	ctx := logging.WithLogger(context.Background())
	logger := logging.FromContext(ctx)

	if len(os.Args) != 6 {
		logger.Fatalf("Invalid args")
	}

	files := os.Args[1]
	gqlEP := os.Args[2]
	tokenURL := os.Args[3]
	clientID := os.Args[4]
	clientSecret := os.Args[5]

	token, err := authToken(ctx, tokenURL, clientID, clientSecret, gqlEP)
	if err != nil {
		logger.Fatalf("Unable to get auth token: %v", err)
	}

	c := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(token),
		},
	}
	gqlClient := graphql.NewClient(gqlEP, c)

	if err := collectFiles(ctx, files, gqlClient); err != nil {
		logger.Fatalf("Unable to send files to guac: %v", err)
	}
}

func authToken(ctx context.Context, tokenURL, clientID string, clientSecret string, gqlEP string) (*oauth2.Token, error) {
	logger := logging.FromContext(ctx)
	if !providers.Enabled(ctx) {
		return nil, fmt.Errorf("incorrect environment")
	}
	token, err := providers.Provide(ctx, tokenURL)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, fmt.Errorf("empty token")
	}
	logger.Infof("ID token aquired")

	gqlUrl, err := url.Parse(gqlEP)
	if err != nil {
		return nil, fmt.Errorf("unable to parse audience from gql-addr")
	}
	audience := fmt.Sprintf("%v://%v", gqlUrl.Scheme, gqlUrl.Hostname())

	var conf oauth2.Config
	conf.Endpoint.TokenURL = tokenURL
	conf.Endpoint.AuthStyle = oauth2.AuthStyleInParams
	options := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer"),
		oauth2.SetAuthURLParam("scope", "openid"),
		oauth2.SetAuthURLParam("client_id", clientID),
		oauth2.SetAuthURLParam("client_secret", clientSecret),
		oauth2.SetAuthURLParam("audience", audience),
		oauth2.SetAuthURLParam("assertion", token),
	}
	tok, err := conf.Exchange(ctx, "", options...)
	if err != nil {
		return nil, err
	}
	if tok.AccessToken == "" {
		return nil, fmt.Errorf("empty token")
	}
	logger.Infof("Access token aquired")
	return tok, nil
}

func collectFiles(ctx context.Context, files string, gqlClient graphql.Client) error {
	logger := logging.FromContext(ctx)

	fileCollector := file.NewFileCollector(ctx, files, false, time.Second)
	if err := collector.RegisterDocumentCollector(fileCollector, file.FileCollector); err != nil {
		return fmt.Errorf("unable to register file collector: %w", err)
	}

	assemblerFunc := helpers.GetBulkAssembler(ctx, gqlClient)

	emit := func(d *processor.Document) error {
		logger.Infof(d.SourceInformation.Source)
		docTree, err := process.Process(ctx, d)
		if err != nil {
			return fmt.Errorf("unable to process doc: %w, format: %v, document: %v", err, d.Format, d.Type)
		}

		predicates, _, err := parser.ParseDocumentTree(ctx, docTree)
		if err != nil {
			return fmt.Errorf("unable to ingest doc tree: %w", err)
		}

		if err := assemblerFunc(predicates); err != nil {
			return fmt.Errorf("unable to assemble graphs: %w", err)
		}
		logger.Infof("completed doc %+v", d.SourceInformation.Source)
		return nil
	}

	errHandler := func(err error) bool {
		if err == nil {
			logger.Info("collector ended gracefully")
			return true
		}
		logger.Errorf("collector ended with error: %v", err)
		return false
	}

	return collector.Collect(ctx, emit, errHandler)
}
