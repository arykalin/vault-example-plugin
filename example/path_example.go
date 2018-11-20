package example

import (
	"context"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// This returns the list of queued for import to TPP certificates
func pathExample(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "example/",

		Fields: map[string]*framework.FieldSchema{
			"user": &framework.FieldSchema{
				Type:        framework.TypeString,
				Required:    true,
				Description: "Enables or disables CORS headers on requests.",
			},
			"comment": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "A comma-separated string or array of strings indicating origins that may make cross-origin requests.",
				Default:     "Comment",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathExampleRead,
			logical.UpdateOperation: b.pathExampleUpdate,
			logical.DeleteOperation: b.pathExampleDelete,
		},

		HelpSynopsis:    "Example plugin path synopsis",
		HelpDescription: "Example plugin path description",
	}
	ret.Fields = map[string]*framework.FieldSchema{}
	return ret
}

func pathExampleList(b *backend) *framework.Path {
	ret := &framework.Path{
		Pattern: "example-list/",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathExampleList,
		},

		HelpSynopsis:    "Example plugin path synopsis",
		HelpDescription: "Example plugin path description",
	}
	return ret
}

func (b *backend) pathExampleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error
	return nil, err
}

func (b *backend) pathExampleUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error
	return nil, err
}

func (b *backend) pathExampleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error
	return nil, err
}

func (b *backend) pathExampleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error
	return nil, err
}
