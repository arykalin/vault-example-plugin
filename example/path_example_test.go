package example

import (
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"log"
	"testing"
)

func TestExample_TestCluster(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"example": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	var err error
	err = client.Sys().Mount("example", &api.MountInput{
		Type: "example",
	})
	if err != nil {
		t.Fatal(err)
	}

	userPass := "MySecret12"
	resp, err := client.Logical().Write("example/user/user1", map[string]interface{}{
		"password": userPass,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected user info")
	}

	// Read user data
	user1, err := client.Logical().Read("example/user/user1")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("user1 data password", user1.Data["password"])
	if user1.Data["password"] != userPass {
		t.Fatalf("expected password %s has %s", userPass, user1.Data["password"])
	}

	resp, err = client.Logical().Write("example/user/user2", map[string]interface{}{
		"generate": true,
	})
	user2, err := client.Logical().Read("example/user/user2")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("user2 password", user2.Data["password"])

	resp, err = client.Logical().List("example/users")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("users list:", resp.Data)
}
