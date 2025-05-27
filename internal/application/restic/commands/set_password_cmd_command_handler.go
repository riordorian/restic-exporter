package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"restic-exporter/internal/application/cqrs"
	logger "restic-exporter/internal/application/log"
	"restic-exporter/internal/application/storage"
)

type SetPasswordCmdCommandHandler struct {
	FileStorage storage.FilesystemInterface
	Logger      logger.LoggerInterface
}

func (c SetPasswordCmdCommandHandler) Handle(ctx context.Context, command cqrs.CommandInterface) error {
	q, ok := command.(SetPasswordCmdCommand)
	if !ok {
		return errors.New("invalid command type. Expected SetPasswordCmd")
	}

	accessFiles, err := c.FileStorage.FindAccessFiles(ctx, q.RootDir)
	if err != nil {
		return fmt.Errorf("failed to find access files: %v", err)
	}

	// TODO: read password from ENV
	passContent, err := os.ReadFile("./pass.txt")
	if err != nil {
		return fmt.Errorf("failed to read pass.txt: %v", err)
	}

	for _, accessFile := range accessFiles {
		if accessFile.Path == "" {
			c.Logger.Warn(fmt.Sprintf("empty path for access file %s", accessFile.Path))
		}

		cmd := exec.Command("restic", "-r", accessFile.Path, "key", "add")
		cmd.Stdin = bytes.NewReader(passContent)
		cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+accessFile.Password)

		output, err := cmd.CombinedOutput()
		if err != nil {
			c.Logger.Error(fmt.Sprintf("restic command failed: %v\nOutput: %s", err, string(output)))
			continue
		}

		c.Logger.Info(fmt.Sprintf("Success for %s. Output: %s\n", accessFile.Path, string(output)))
	}

	return nil
}
