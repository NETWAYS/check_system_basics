package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	// "github.com/NETWAYS/check_system_basics/internal/common/status"

	sbRegex "github.com/NETWAYS/check_system_basics/internal/common/regexp"
	fileContent "github.com/NETWAYS/check_system_basics/internal/files"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

var FileContentConfig fileContent.FileContentconfig

type EvalConfig struct {
	OKPatterns       []regexp.Regexp
	WarningPatterns  []regexp.Regexp
	CriticalPatterns []regexp.Regexp
}

var fileContentCmd = &cobra.Command{
	Use:     "fileContent",
	Short:   "Submodule to test for different properties on file content ",
	Example: ``,
	Run: func(_ *cobra.Command, _ []string) {
		if len(FileContentConfig.Paths) == 0 {
			check.Exit(check.Unknown, "At least one path (--paths) must be selected")
		}

		// Input sanity check
		for _, inputPath := range FileContentConfig.Paths {
			if !path.IsAbs(inputPath) {
				check.Exit(check.Unknown, fmt.Sprintf("Path %s is not an absolute path, but must be one", inputPath))
			}
		}

		overall := result.Overall{}

		// find files
		for _, inputPath := range FileContentConfig.Paths {
			sc, err := PathEvaluation(inputPath, FileContentConfig)
			if err != nil {
				check.ExitError(err)
			}

			overall.AddSubcheck(sc)
		}

		check.Exit(overall.GetStatus(), overall.GetOutput())
	},
}

func PathEvaluation(path string, config fileContent.FileContentconfig) (*result.PartialResult, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		// Input path is a directory
		// apply conditions for all files inside
		partialDir := result.NewPartialResult()
		partialDir.SetDefaultState(check.OK)
		partialDir.SetOutput(path)

		dirs, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, dirEntry := range dirs {
			if dirEntry.IsDir() {
				if config.Recursive {
					// TODO head down
					ssc, err := PathEvaluation(filepath.Join(path, dirEntry.Name()), config)
					if err != nil {
						return nil, err
					}

					partialDir.AddSubcheck(ssc)
				}

				// non recursive, ignore the dir
				continue
			}

			// it's a file!
			fileSC, err := EvaluateFile(filepath.Join(path, dirEntry.Name()), config)
			if err != nil {
				return nil, err
			}

			partialDir.AddSubcheck(fileSC)
		}

		return partialDir, nil
	}

	// Input path is a file
	// apply conditions directly
	return EvaluateFile(path, config)
}

func EvaluateFile(path string, config fileContent.FileContentconfig) (*result.PartialResult, error) {
	evaluationResult := result.NewPartialResult()
	evaluationResult.SetDefaultState(check.OK)

	baseName := filepath.Base(path)
	evaluationResult.SetOutput(baseName)

	// Evaluation
	// -- Pattern matching in file
	if (len(config.OKPatterns) != 0) || (len(config.WarningPatterns) != 0) || (len(config.CriticalPatterns) != 0) {
		scPattern := result.NewPartialResult()
		scPattern.SetDefaultState(config.NotFoundStatus.Status)

		// Pattern matching priority:
		// if Critical > Warning > OK
		// start with critical and first match wins
		foundPattern := false
		foundWarningPattern := false
		foundCriticalPattern := false

		var patternFound sbRegex.SBRegex

		fileDesc, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fileDesc)

		for scanner.Scan() {
			for _, critPattern := range config.CriticalPatterns {
				if critPattern.Regex.MatchString(scanner.Text()) {
					foundPattern = true
					foundCriticalPattern = true
					patternFound = critPattern

					break
				}
			}

			if foundCriticalPattern {
				// abort here, if we already have CRITICAL
				break
			}

			for _, warnPattern := range config.WarningPatterns {
				if warnPattern.Regex.MatchString(scanner.Text()) {
					foundPattern = true
					foundWarningPattern = true
					patternFound = warnPattern

					break
				}
			}

			if foundWarningPattern {
				// abort here, if we already have Warning
				break
			}

			for _, okPattern := range config.OKPatterns {
				if okPattern.Regex.MatchString(scanner.Text()) {
					foundPattern = true
					patternFound = okPattern

					break
				}
			}
		}

		// Scanner failure?
		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		if !foundPattern {
			scPattern.SetState(config.NotFoundStatus.Status)
			scPattern.SetOutput("Regex pattern did not match in file")
		} else {
			// ok we found something
			scPattern.SetOutput("Found pattern \"" + patternFound.String() + "\"")

			if foundCriticalPattern {
				scPattern.SetState(check.Critical)
			} else if foundWarningPattern {
				scPattern.SetState(check.Warning)
			} else {
				scPattern.SetState(check.OK)
			}
		}

		evaluationResult.AddSubcheck(scPattern)
	}

	return evaluationResult, nil
}

func init() {
	rootCmd.AddCommand(fileContentCmd)
	fileContentCmd.DisableFlagsInUseLine = true

	fileContentFS := fileContentCmd.Flags()
	fileContentFS.StringArrayVar(&FileContentConfig.Paths, "paths", []string{}, "File paths to evaluate")
	fileContentFS.Var(&FileContentConfig.OKPatterns, "ok-pattern", "Regex pattern in file which are OK")
	fileContentFS.Var(&FileContentConfig.WarningPatterns, "warning-pattern", "Regex pattern in file which cause a WARNING")
	fileContentFS.Var(&FileContentConfig.CriticalPatterns, "critical-pattern", "Regex pattern in file which cause a CRITICAL")
	fileContentFS.Var(&FileContentConfig.NotFoundStatus, "not-found-status", "Exit status if none of the patterns apply (OK (0), warning (1), critical (2), Uknown (3)) (default: OK)")
	fileContentFS.BoolVar(&FileContentConfig.Recursive, "recursive", false, "Recursively test all files in \"paths\"")
}
