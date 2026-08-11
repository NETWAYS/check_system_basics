package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	sbRegex "github.com/NETWAYS/check_system_basics/internal/common/regexp"
	fileContent "github.com/NETWAYS/check_system_basics/internal/files"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/spf13/cobra"
)

var FileContentConfig fileContent.FileContentconfig

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
	if (len(config.OKPatterns) != 0) || (len(config.WarningPatterns) != 0) || (len(config.CriticalPatterns) != 0) || config.MetricPattern.IsSet {
		// Pattern matching priority:
		// if Critical > Warning > OK
		// start with critical and first match wins
		foundOKPattern := false
		foundWarningPattern := false
		foundCriticalPattern := false
		foundMetric := false

		var metric float64

		var patternFound sbRegex.SBRegex

		fileDesc, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fileDesc)

		for scanner.Scan() {
			if !foundCriticalPattern {
				for _, critPattern := range config.CriticalPatterns {
					if critPattern.Regex.MatchString(scanner.Text()) {
						foundCriticalPattern = true
						patternFound = critPattern

						break
					}
				}
			}

			if !foundWarningPattern {
				for _, warnPattern := range config.WarningPatterns {
					if warnPattern.Regex.MatchString(scanner.Text()) {
						foundWarningPattern = true
						patternFound = warnPattern

						break
					}
				}
			}

			if !foundOKPattern {
				for _, okPattern := range config.OKPatterns {
					if okPattern.Regex.MatchString(scanner.Text()) {
						foundOKPattern = true
						patternFound = okPattern

						break
					}
				}
			}

			if config.MetricPattern.IsSet {
				if config.MetricPattern.Regex.MatchString(scanner.Text()) {
					foundMetric = true
					matchSlice := config.MetricPattern.Regex.FindStringSubmatch(scanner.Text())

					if matchSlice == nil {
						check.ExitError(fmt.Errorf("Metrix Regex submatch failed somehow"))
					}

					metric, err = strconv.ParseFloat(matchSlice[1], 64)
					if err != nil {
						check.ExitError(err)
					}
				}
			}

			if (config.MetricPattern.IsSet || foundMetric) &&
				(len(config.CriticalPatterns) == 0 || foundCriticalPattern) &&
				(len(config.WarningPatterns) == 0 || foundWarningPattern) &&
				(len(config.OKPatterns) == 0 || foundOKPattern) {
				// Abort if we are done
				break
			}
		}

		// Scanner failure?
		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		scPattern := result.NewPartialResult()
		scPattern.SetDefaultState(config.NotFoundStatus.Status)

		if !foundCriticalPattern && !foundWarningPattern && !foundOKPattern {
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

		if config.MetricPattern.IsSet {
			// Expected a metric match
			scMetric := result.NewPartialResult()

			if !foundMetric {
				scMetric.SetState(config.MetricNotFoundStatus.Status)
				scMetric.SetOutput("Metric not found")
			} else {
				scMetric.SetState(check.OK)
				scMetric.SetOutput(fmt.Sprintf("%s: %g", config.MetricLabel, metric))

				pdMetrcis := check.Perfdata{
					Value: metric,
					Label: config.MetricLabel,
				}

				if config.MetricThresholds.Warn.IsSet {
					pdMetrcis.Warn = &config.MetricThresholds.Warn.Th
				}

				if config.MetricThresholds.Warn.Th.DoesViolate(metric) {
					scMetric.SetState(check.Warning)
				}

				if config.MetricThresholds.Crit.IsSet {
					pdMetrcis.Crit = &config.MetricThresholds.Crit.Th
				}

				if config.MetricThresholds.Crit.Th.DoesViolate(metric) {
					scMetric.SetState(check.Critical)
				}

				scMetric.AddPerfdata(&pdMetrcis)
			}

			evaluationResult.AddSubcheck(scMetric)
		}
	}

	return evaluationResult, nil
}

func init() {
	rootCmd.AddCommand(fileContentCmd)
	fileContentCmd.DisableFlagsInUseLine = true

	fileContentFS := fileContentCmd.Flags()
	fileContentFS.StringArrayVar(&FileContentConfig.Paths, "paths", []string{}, "File paths to evaluate")
	fileContentFS.BoolVar(&FileContentConfig.Recursive, "recursive", false, "Recursively test all files in \"paths\"")

	fileContentFS.Var(&FileContentConfig.OKPatterns, "ok-pattern", "Regex pattern in file which are OK")
	fileContentFS.Var(&FileContentConfig.WarningPatterns, "warning-pattern", "Regex pattern in file which cause a WARNING")
	fileContentFS.Var(&FileContentConfig.CriticalPatterns, "critical-pattern", "Regex pattern in file which cause a CRITICAL")
	fileContentFS.Var(&FileContentConfig.NotFoundStatus, "not-found-status", "Exit status if none of the patterns apply (OK (0), warning (1), critical (2), Uknown (3)) (default: OK)")

	fileContentFS.Var(&FileContentConfig.MetricPattern, "metric-pattern", "Regex pattern to find numerical values in the file")
	fileContentFS.StringVar(&FileContentConfig.MetricLabel, "metric-label", "metric", "(Perfdata) label for matched metrics")
	fileContentFS.Var(&FileContentConfig.MetricThresholds.Warn, "metric-warning", "Warning threshold for the matched metric")
	fileContentFS.Var(&FileContentConfig.MetricThresholds.Crit, "metric-critical", "Critical threshold for the matched metric")
	fileContentFS.Var(&FileContentConfig.MetricNotFoundStatus, "metric-not-found-status", "Exit status if the metric patterns were not found. (OK (0), warning (1), critical (2), Uknown (3)) (default: OK)")
}
