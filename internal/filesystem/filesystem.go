package filesystem

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// nolint: revive, golint
type FilesystemType struct {
	PartStats  disk.PartitionStat
	UsageStats disk.UsageStat
	Error      error
}

type tmpFileSystemWrapper struct {
	usage disk.UsageStat
	err   error
}

func GetDiskUsageSingle(ctx context.Context, timeout time.Duration, fs *FilesystemType) {
	myCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resChan := make(chan tmpFileSystemWrapper, 1)

	go func() {
		tmp := tmpFileSystemWrapper{}
		usageStats, err := disk.Usage(fs.PartStats.Mountpoint)

		if err == nil {
			tmp.usage = *usageStats
		}

		tmp.err = err

		resChan <- tmp
	}()

	select {
	case tmp := <-resChan:
		if tmp.err != nil {
			fs.Error = tmp.err
			return
		}

		fs.UsageStats = tmp.usage
	case <-myCtx.Done():
		err := errors.New("Timeout exceeded for fs " + fs.PartStats.Mountpoint + ". Maybe hanging network filesystem?")
		fs.Error = err
	}
}

func GetDiskUsage(ctx context.Context, timeout time.Duration, fsList []FilesystemType, _ *CheckConfig) error {
	for index := range fsList {
		GetDiskUsageSingle(ctx, timeout/time.Duration(len(fsList)), &fsList[index])
	}

	return nil
}

func FilterFileSystem(filesystems []FilesystemType, filters *Filters) ([]FilesystemType, error) {
	// Filesystem Type
	if len(filters.ExcludeFsType) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isExcluded := false

			for _, exclude := range filters.ExcludeFsType {
				match, err := regexp.MatchString(exclude, filesystems[fs].PartStats.Fstype)
				if err != nil {
					return []FilesystemType{}, err
				}

				if match {
					isExcluded = true
					break
				}
			}

			if !isExcluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	if len(filters.IncludeFsType) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isIncluded := false

			for _, include := range filters.IncludeFsType {
				match, err := regexp.MatchString(include, filesystems[fs].PartStats.Fstype)
				if err != nil {
					return []FilesystemType{}, err
				}

				if match {
					isIncluded = true
					break
				}
			}

			if isIncluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	// Device paths
	if len(filters.ExcludeDevicePaths) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isExcluded := false

			for _, exclude := range filters.ExcludeDevicePaths {
				match, err := regexp.MatchString(exclude, filesystems[fs].PartStats.Device)
				if err != nil {
					return []FilesystemType{}, err
				}

				if match {
					isExcluded = true
					break
				}
			}

			if !isExcluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	if len(filters.IncludeDevicePaths) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isIncluded := false

			for _, include := range filters.IncludeDevicePaths {
				match, err := regexp.MatchString(include, filesystems[fs].PartStats.Device)
				if err != nil {
					return []FilesystemType{}, err
				}

				if match {
					isIncluded = true
					break
				}
			}

			if isIncluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	// Mount paths
	if len(filters.ExcludeMountPaths) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isExcluded := false

			for _, exclude := range filters.ExcludeMountPaths {
				match, err := regexp.MatchString(exclude, filesystems[fs].PartStats.Mountpoint)
				if err != nil {
					return []FilesystemType{}, err
				}

				if match {
					isExcluded = true
					break
				}
			}

			if !isExcluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	if len(filters.IncludeMountPaths) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isIncluded := false

			for _, include := range filters.IncludeMountPaths {
				match, err := regexp.MatchString(include, filesystems[fs].PartStats.Mountpoint)
				if err != nil {
					return []FilesystemType{}, err
				}

				if match {
					isIncluded = true
					break
				}
			}

			if isIncluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	// Mount options
	if len(filters.ExcludeOptions) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isExcluded := false

			for _, exclude := range filters.ExcludeOptions {
				for _, mountOption := range filesystems[fs].PartStats.Opts {
					match, err := regexp.MatchString(exclude, mountOption)
					if err != nil {
						return []FilesystemType{}, err
					}

					if match {
						isExcluded = true
						break
					}
				}

				if isExcluded {
					break
				}
			}

			if !isExcluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	if len(filters.IncludeOptions) != 0 {
		newList := make([]FilesystemType, 0)

		for fs := range filesystems {
			isIncluded := false

			for _, include := range filters.IncludeOptions {
				for _, mountOption := range filesystems[fs].PartStats.Opts {
					match, err := regexp.MatchString(include, mountOption)
					if err != nil {
						return []FilesystemType{}, err
					}

					if match {
						isIncluded = true
						break
					}
				}

				if isIncluded {
					break
				}
			}

			if isIncluded {
				newList = append(newList, filesystems[fs])
			}
		}

		filesystems = newList
	}

	return filesystems, nil
}
