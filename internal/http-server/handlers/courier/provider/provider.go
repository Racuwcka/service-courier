package provider

import "context"

type Provider interface {
	Create(ctx context.Context, name string) (uint32, error)
}

type CourierProvider struct{}

func NewCourierProvider() *CourierProvider {
	return &CourierProvider{}
}

func (p *CourierProvider) Create(_ context.Context, _ string) (uint32, error) {
	return 123, nil
}
