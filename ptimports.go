// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"github.com/palantir/pkg/cobracli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/palantir/go-ptimports/ptimports"
)

var rootCmd = &cobra.Command{
	Use: "go-ptimports",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ptimports.ProcessFileFromInput("", os.Stdin, false, false, refactorFlagVal, localPrefixFlagVal, cmd.OutOrStdout())
		}
		for _, currFile := range args {
			if err := ptimports.ProcessFileFromInput(currFile, nil, listFlagVal, writeFlagVal, refactorFlagVal, localPrefixFlagVal, cmd.OutOrStdout()); err != nil {
				return errors.Wrapf(err, "failed to process file %s", currFile)
			}
		}
		return nil
	},
}

var (
	debugFlagVal    bool
	refactorFlagVal bool

	listFlagVal        bool
	writeFlagVal       bool
	localPrefixFlagVal []string
)

func init() {
	rootCmd.Flags().BoolVarP(&refactorFlagVal, "refactor", "r", false, "refactor imports to use block style imports")

	rootCmd.Flags().BoolVarP(&listFlagVal, "list", "l", false, "list files whose formatting differs from go-ptimport's")
	rootCmd.Flags().BoolVarP(&writeFlagVal, "write", "w", false, "write result to (source) file instead of stdout")
	rootCmd.Flags().StringSliceVar(&localPrefixFlagVal, "local", nil, "put imports beginning with this string after 3rd-party packages; comma-separated list")
}

func main() {
	os.Exit(cobracli.ExecuteWithDefaultParamsWithVersion(rootCmd, &debugFlagVal, ""))
}
