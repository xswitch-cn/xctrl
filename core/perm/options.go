package perm

import (
	"time"
)

// Options token 选项
type Options struct {
	// Token is an auth token
	Token string
	// Public key base64 encoded
	PublicKey string
	// Private key base64 encoded
	PrivateKey string
	// Endpoints to exclude
	Exclude []string
}

// Option token options
type Option func(o *Options)

// Exclude ecludes a set of endpoints from authorization
func Exclude(e ...string) Option {
	return func(o *Options) {
		o.Exclude = e
	}
}

// PublicKey is the JWT public key
func PublicKey(key string) Option {
	return func(o *Options) {
		o.PublicKey = key
	}
}

// PrivateKey is the JWT private key
func PrivateKey(key string) Option {
	return func(o *Options) {
		o.PrivateKey = key
	}
}

// Token sets an auth token
func Token(t string) Option {
	return func(o *Options) {
		o.Token = t
	}
}

// GenerateOptions token 生成参数
type GenerateOptions struct {
	// Metadata associated with the account
	Metadata map[string]string
	// Roles/scopes associated with the account
	Roles []string
	//Expiry of the token
	Expiry time.Time
}

// GenerateOption token 生成参数
type GenerateOption func(o *GenerateOptions)

// Metadata for the generated account
func Metadata(md map[string]string) func(o *GenerateOptions) {
	return func(o *GenerateOptions) {
		o.Metadata = md
	}
}

// Expiry for the generated account's token expires
func Expiry(ex time.Time) func(o *GenerateOptions) {
	return func(o *GenerateOptions) {
		o.Expiry = ex
	}
}

// NewGenerateOptions from a slice of options
func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	var options GenerateOptions
	for _, o := range opts {
		o(&options)
	}
	//set defualt expiry of token
	if options.Expiry.IsZero() {
		options.Expiry = time.Now().Add(time.Hour * 24)
	}
	return options
}
