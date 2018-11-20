package example

import (
	"context"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"math/rand"
	"time"
)

func examplePaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "users/?",
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
		},
	}
	return resp, nil
}

func (b *backend) pathExampleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var password string
	if data.Get("generate").(bool) {
		password = randSeq(6)
	} else {
		password = data.Get("password").(string)
		if len(password) == 0 {
			return nil, fmt.Errorf("Must provide password or generate\n")
		}
	}
	comment := data.Get("comment").(string)
	user := data.Get("user").(string)

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
			"user":     user,
			"password": password,
			"comment":  comment,
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
	vals, err := req.Storage.List(ctx, "user/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
