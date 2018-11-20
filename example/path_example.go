package example

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func examplePaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "example/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathExampleList,
			},
		},
		&framework.Path{
			Pattern: "user/" + framework.GenericNameRegex("user"),
			Fields: map[string]*framework.FieldSchema{
				"user": &framework.FieldSchema{
					Type: framework.TypeString},
				"comment": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "example user description",
					Default:     "empty",
					Required:    false},
				"password": &framework.FieldSchema{Type: framework.TypeString},
				"version":  &framework.FieldSchema{Type: framework.TypeInt},
				"generate": &framework.FieldSchema{Type: framework.TypeBool},
			},
			ExistenceCheck: b.pathExampleExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathExampleRead,
				logical.CreateOperation: b.pathExampleCreateUpdate,
				logical.UpdateOperation: b.pathExampleCreateUpdate,
				logical.DeleteOperation: b.pathExampleDelete,
			},
		},
	}
}

func (b *backend) pathExampleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

func (b *backend) pathExampleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	version := data.Get("version").(int)

	entry, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	password := string(entry.Value)

	b.Logger().Info("reading password", "user", req.Path, "password", password)
	// Return the secret
	resp := &logical.Response{
		Data: map[string]interface{}{
			"password": password,
			"version":  version,
		},
	}
	if version != 0 {
		resp.Data["version"] = version
	}
	return resp, nil
}

func (b *backend) pathExampleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	password := data.Get("password").(string)

	b.Logger().Info("storing password", "user", req.Path, "password", password)
	entry := &logical.StorageEntry{
		Key:   req.Path,
		Value: []byte(password),
	}

	s := req.Storage
	err := s.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"password": password,
		},
	}, nil
}

func (b *backend) pathExampleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathExampleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List(ctx, "users/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}
