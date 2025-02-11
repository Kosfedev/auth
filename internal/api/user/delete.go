package user

import (
	"context"
	"log"
)

// Delete is...
func (i *Implementation) Delete(ctx context.Context, id int64) error {
	err := i.userService.Delete(ctx, id)
	if err != nil {
		return err
	}

	log.Printf("user %d has been deleted\n", id)

	return nil
}
