package galapagos

import "github.com/otto-de/terraform-provider-galapagos/internal/rest"

var (
	ACCOUNT_REST     = rest.RESTType{Plural: "accounts", Singular: "account"}
	APPLICATION_REST = rest.RESTType{Plural: "applications", Singular: "application"}
)
