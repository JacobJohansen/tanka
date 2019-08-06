package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func applyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [directory]",
		Short: "apply the configuration to the cluster",
		Args:  cobra.ExactArgs(1),
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		if err := kube.Apply(desired); err != nil {
			log.Fatalln("Applying:", err)
		}
	}
	return cmd
}

func diffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [directory]",
		Short: "differences between the configuration and the cluster",
		Args:  cobra.ExactArgs(1),
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		desired, err := kube.Reconcile(raw)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		changes, err := kube.Diff(desired)
		if err != nil {
			log.Fatalln("Diffing:", err)
		}

		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			if err := quick.Highlight(os.Stdout, changes, "diff", "terminal", "vim"); err != nil {
				log.Fatalln("Highlighting:", err)
			}
		} else {
			fmt.Println(changes)
		}
	}
	return cmd
}

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [directory]",
		Short: "jsonnet as yaml",
		Args:  cobra.ExactArgs(1),
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		raw, err := evalDict(args[0])
		if err != nil {
			log.Fatalln("Evaluating jsonnet:", err)
		}

		state, err := kube.Reconcile(raw)
		if err != nil {
			log.Fatalln("Reconciling:", err)
		}

		pretty, err := kube.Fmt(state)
		if err != nil {
			log.Fatalln("Pretty printing state:", err)
		}
		fmt.Println(pretty)
	}
	return cmd
}