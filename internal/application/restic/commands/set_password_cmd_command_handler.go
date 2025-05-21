package commands

import (
	"context"
	"errors"
	"fmt"
	"restic-exporter/internal/application/cqrs"
	"restic-exporter/internal/application/storage"
)

type SetPasswordCmdCommandHandler struct {
	FileStorage storage.FilesystemInterface
}

func (c SetPasswordCmdCommandHandler) Handle(ctx context.Context, command cqrs.CommandInterface) error {
	if q, ok := command.(SetPasswordCmdCommand); ok {
		accessFiles, err := c.FileStorage.FindAccessFiles(ctx, q.RootDir)
		if err != nil {
			return err
		}

		fmt.Println(accessFiles)

		return nil
	}

	return errors.New("invalid command type. Expected SetPasswordCmd")
}
