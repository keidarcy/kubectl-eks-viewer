package main

import (
	"os"

	"github.com/keidarcy/kubectl-eks-viewer/pkg/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	command := cmd.NewCmd(streams)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
