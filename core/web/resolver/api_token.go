package resolver

import "github.com/smartcontractkit/chainlink/core/auth"

type APITokenResolver struct {
	token auth.Token
}

func NewAPIToken(token auth.Token) *APITokenResolver {
	return &APITokenResolver{token}
}

func (r *APITokenResolver) AccessKey() string {
	return r.token.AccessKey
}

func (r *APITokenResolver) Secret() string {
	return r.token.Secret
}

// -- CreateAPIToken Mutation --

type CreateAPITokenPayloadResolver struct {
	token     *auth.Token
	inputErrs map[string]string
}

func NewCreateAPITokenPayload(token *auth.Token, inputErrs map[string]string) *CreateAPITokenPayloadResolver {
	return &CreateAPITokenPayloadResolver{token, inputErrs}
}

func (r *CreateAPITokenPayloadResolver) ToCreateAPITokenSuccess() (*CreateAPITokenSuccessResolver, bool) {
	if r.inputErrs != nil {
		return nil, false
	}

	return NewCreateAPITokenSuccess(r.token), true
}

func (r *CreateAPITokenPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type CreateAPITokenSuccessResolver struct {
	token *auth.Token
}

func NewCreateAPITokenSuccess(token *auth.Token) *CreateAPITokenSuccessResolver {
	return &CreateAPITokenSuccessResolver{token}
}

func (r *CreateAPITokenSuccessResolver) Token() *APITokenResolver {
	return NewAPIToken(*r.token)
}

// -- DeleteAPIToken Mutation --

type DeleteAPITokenPayloadResolver struct {
	token     *auth.Token
	inputErrs map[string]string
}

func NewDeleteAPITokenPayload(token *auth.Token, inputErrs map[string]string) *DeleteAPITokenPayloadResolver {
	return &DeleteAPITokenPayloadResolver{token, inputErrs}
}

func (r *DeleteAPITokenPayloadResolver) ToDeleteAPITokenSuccess() (*DeleteAPITokenSuccessResolver, bool) {
	if r.inputErrs != nil {
		return nil, false
	}

	return NewDeleteAPITokenSuccess(r.token), true
}

func (r *DeleteAPITokenPayloadResolver) ToInputErrors() (*InputErrorsResolver, bool) {
	if r.inputErrs != nil {
		var errs []*InputErrorResolver

		for path, message := range r.inputErrs {
			errs = append(errs, NewInputError(path, message))
		}

		return NewInputErrors(errs), true
	}

	return nil, false
}

type DeleteAPITokenSuccessResolver struct {
	token *auth.Token
}

func NewDeleteAPITokenSuccess(token *auth.Token) *DeleteAPITokenSuccessResolver {
	return &DeleteAPITokenSuccessResolver{token}
}

func (r *DeleteAPITokenSuccessResolver) Token() *APITokenResolver {
	return NewAPIToken(*r.token)
}
