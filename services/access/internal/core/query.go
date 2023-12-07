package core

type AuthorizeRequest struct {
	Token      string
	Permission string
	Resource   string
}

type AuthorizeResponse struct {
	HavePermission bool
}
