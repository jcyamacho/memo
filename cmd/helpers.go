package cmd

import (
	"fmt"
	"io"
	"os"
)

func readContent(args []string) (string, error) {
	return readContentFrom(args, stdinReader(os.Stdin))
}

func readContentFrom(args []string, stdin io.Reader) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	if stdin == nil {
		return "", fmt.Errorf("content is required")
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	return string(data), nil
}

func stdinReader(stdin *os.File) io.Reader {
	info, err := stdin.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice != 0 {
		return nil
	}
	return stdin
}

func printOutput(output string, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
