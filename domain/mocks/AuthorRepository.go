// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/syahidfrd/go-boilerplate/domain"
)

// AuthorRepository is an autogenerated mock type for the AuthorRepository type
type AuthorRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, author
func (_m *AuthorRepository) Create(ctx context.Context, author *domain.Author) error {
	ret := _m.Called(ctx, author)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Author) error); ok {
		r0 = rf(ctx, author)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *AuthorRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx
func (_m *AuthorRepository) Fetch(ctx context.Context) ([]domain.Author, error) {
	ret := _m.Called(ctx)

	var r0 []domain.Author
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Author); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Author)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *AuthorRepository) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Author
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Author); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Author)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, author
func (_m *AuthorRepository) Update(ctx context.Context, author *domain.Author) error {
	ret := _m.Called(ctx, author)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Author) error); ok {
		r0 = rf(ctx, author)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
