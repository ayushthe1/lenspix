package context

import (
	"context"

	"github.com/ayushthe1/lenspix/models"
)

// unexported keys and key types makes it more safer to use context values as other packages can't edit the values in our context
type key string

const (
	userKey key = "user"
)

// function to store a user inside of a context
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// function to retrieve a user from the context
func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	// assert the val as user type
	user, ok := val.(*models.User)
	if !ok {
		// The most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key.
		return nil
		// this nil value will end up getting converted into type *models.User. So it will be still be a pointer to a user but have a value of nil
	}

	return user

}
