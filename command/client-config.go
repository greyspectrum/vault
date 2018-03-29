package command

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/command/client-config"
)

const (
	// DefaultClientConfigPath is the default path to the client configuration file
	DefaultClientConfigPath = "~/.vaultrc"

	//ClientConfigPathEnv is the environment variable that can be used to
	//override where the Vault configuration is.
	ClientConfigPathEnv = "VAULT_CLIENT_CONFIG_PATH"
)

// ClientConfig is the CLI configuration for Vault that can be specified via
// a `$HOME/.vaultrc` file which is HCL-formatted (therefore HCL or JSON).
type DefaultClientConfig struct {
	// TokenHelper is the executable/command that is executed for storing
	// and retrieving the authentication token for the Vault CLI. If this
	// is not specified, then vault's internal token store will be used, which
	// stores the token on disk unencrypted.
	TokenHelper string `hcl:"token_helper"`
}

// Config loads the configuration and returns it. If the configuration
// is already loaded, it is returned.
//
// Config just calls into config.Config for backwards compatibility purposes.
// Use config.Config instead.
func ClientConfig() (*DefaultClientConfig, error) {
	conf, err := config.ClientConfig()
	return (*DefaultClientConfig)(conf), err
}

// LoadClientConfig reads the configuration from the given path. If path is
// empty, then the default path will be used, or the environment variable
// if set.
//
// LoadClientConfig just calls into config.LoadClientConfig for backwards compatibility
// purposes. Use config.LoadClientConfig instead.
func LoadClientConfig(path string) (*DefaultConfig, error) {
	conf, err := config.LoadClientConfig(path)
	return (*DefaultClientConfig)(conf), err
}

// ParseClientConfig parses the given client configuration as a string.
//
// ParseClientConfig just calls into config.ParseClientConfig for backwards compatibility
// purposes. Use config.ParseClientConfig instead.
func ParseClientConfig(contents string) (*DefaultClientConfig, error) {
	conf, err := config.ParseClientConfig(contents)
	return (*DefaultClientConfig)(conf), err
}

func checkHCLKeys(node ast.Node, valid []string) error {
	var list *ast.ObjectList
	switch n := node.(type) {
	case *ast.ObjectList:
		list = n
	case *ast.ObjectType:
		list = n.List
	default:
		return fmt.Errorf("cannot check HCL keys of type %T", n)
	}

	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := validMap[key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invalid key '%s' on line %d", key, item.Assign.Line))
		}
	}

	return result
}
