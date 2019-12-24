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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/palantir/pkg/cobracli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/x0rzkov/go-ptimports/v2/ptimports"
)

var rootCmd = &cobra.Command{
	Use: "go-ptimports",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &ptimports.Options{
			Refactor:      refactorFlagVal,
			Simplify:      simplifyFlagVal,
			FormatOnly:    formatOnlyFlagVal,
			LocalPrefixes: localPrefixFlagVal,
			// RecursiveDir:  recursiveDirFlagVal,
		}

		if recursiveDirFlagVal {
			var dirs []string
			if len(args) == 0 {
				dirs = append(dirs, ".")
			} else {
				for _, arg := range args {
					dirs = append(dirs, arg)
				}
			}
			for _, dirname := range dirs {
				err := godirwalk.Walk(dirname, &godirwalk.Options{
					Callback: func(currFile string, de *godirwalk.Dirent) error {
						if skipVendorFlagVal && strings.HasPrefix(currFile, "vendor") {
							if verboseFlagVal {
								fmt.Printf("skipping %s \n", currFile)
							}
							return nil
						}
						if strings.HasSuffix(currFile, ".go") {
							if verboseFlagVal {
								fmt.Printf("processing %s %s\n", de.ModeType(), currFile)
							}
							if err := ptimports.ProcessFileFromInput(currFile, nil, listFlagVal, writeFlagVal, opts, cmd.OutOrStdout()); err != nil {
								return errors.Wrapf(err, "failed to process file %s", currFile)
							}
						}
						return nil
					},
					Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
				})
				if err != nil {
					return errors.Wrapf(err, "failed to process directory %s", dirname)
				}
			}
		} else {
			if len(args) == 0 {
				return ptimports.ProcessFileFromInput("", os.Stdin, false, false, opts, cmd.OutOrStdout())
			}
			for _, currFile := range args {
				if err := ptimports.ProcessFileFromInput(currFile, nil, listFlagVal, writeFlagVal, opts, cmd.OutOrStdout()); err != nil {
					return errors.Wrapf(err, "failed to process file %s", currFile)
				}
			}
		}
		return nil
	},
}

var (
	simplifyFlagVal        bool
	refactorFlagVal        bool
	formatOnlyFlagVal      bool
	recursiveDirFlagVal    bool
	skipVendorFlagVal      bool
	localAutoDetectFlagVal bool
	localPrefixFlagVal     []string
	listFlagVal            bool
	writeFlagVal           bool
	verboseFlagVal         bool
)

func Execute() int {
	return cobracli.ExecuteWithDefaultParams(rootCmd)
}

func init() {
	rootCmd.Flags().BoolVarP(&simplifyFlagVal, "simplify", "s", false, "simplify code in the manner that gofmt does")
	rootCmd.Flags().BoolVarP(&refactorFlagVal, "refactor", "r", false, "refactor imports to use block style imports")
	rootCmd.Flags().BoolVar(&formatOnlyFlagVal, "format-only", false, "do not add or remove imports")
	rootCmd.Flags().StringSliceVar(&localPrefixFlagVal, "local", nil, "put imports beginning with this string after 3rd-party packages; comma-separated list")

	rootCmd.Flags().BoolVarP(&verboseFlagVal, "verbose", "v", false, "verbose mode")
	rootCmd.Flags().BoolVarP(&localAutoDetectFlagVal, "local-detect", "o", false, "auto-detect local")
	rootCmd.Flags().BoolVarP(&skipVendorFlagVal, "skip-vendor", "k", false, "skip vendor dir.")
	rootCmd.Flags().BoolVarP(&recursiveDirFlagVal, "recursive-dir", "d", false, "walk through directory.")
	rootCmd.Flags().BoolVarP(&listFlagVal, "list", "l", false, "list files whose formatting differs from go-ptimport's")
	rootCmd.Flags().BoolVarP(&writeFlagVal, "write", "w", false, "write result to (source) file instead of stdout")
}
