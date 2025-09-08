package filesystem

import (
	"reflect"
	"testing"

	"github.com/shirou/gopsutil/v3/disk"
)

var (
	testSda = FilesystemType{
		PartStats: disk.PartitionStat{
			Device:     "/dev/sda",
			Mountpoint: "/",
			Fstype:     "ext4",
			Opts:       []string{"rw", "relatime", "errors=remount-ro"},
		},
		UsageStats: disk.UsageStat{
			Path:              "/",
			Fstype:            "ext2/ext3",
			Total:             0xe975d1000,
			Free:              0x7b4d4b000,
			Used:              0x622ced000,
			UsedPercent:       44.32754032726998,
			InodesTotal:       0x3ba000,
			InodesUsed:        0xa397b,
			InodesFree:        0x316685,
			InodesUsedPercent: 17.14806386759696,
		},
		Error: nil,
	}
	testFsHalfFull = FilesystemType{
		PartStats: disk.PartitionStat{
			Device:     "/dev/testfilesystem",
			Mountpoint: "/testMountpoint",
			Fstype:     "test",
			Opts:       []string{},
		},
		UsageStats: disk.UsageStat{
			Path:              "/testMountpoint",
			Fstype:            "test",
			Total:             4194304, // 1024 pages aka  4096 * 1024
			Free:              2097152, // 512 pages, -> 50%
			Used:              2097152,
			UsedPercent:       50,
			InodesTotal:       1024,
			InodesUsed:        512,
			InodesFree:        512,
			InodesUsedPercent: 50,
		},
	}
	fileSystemSet = []FilesystemType{
		testSda,
		testFsHalfFull,
	}
)

func TestFilterSDA(t *testing.T) {
	filters := Filters{
		IncludeFsType: []string{"ext4"},
	}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{testSda}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{testSda})
	}
}

func TestFilterHalfFull(t *testing.T) {
	filters := Filters{
		ExcludeFsType: []string{"ext4"},
	}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{testFsHalfFull}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{testFsHalfFull})
	}
}

func TestFilterNoFilter(t *testing.T) {
	filters := Filters{}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{testSda, testFsHalfFull}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{testSda, testFsHalfFull})
	}
}

func TestFilterIncludeDevicePaths(t *testing.T) {
	filters := Filters{
		IncludeDevicePaths: []string{"/dev/"},
	}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{testSda, testFsHalfFull}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{testSda, testFsHalfFull})
	}
}

func TestFilterExcludeDevicePaths(t *testing.T) {
	filters := Filters{
		ExcludeDevicePaths: []string{"/dev/s"},
	}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{testFsHalfFull}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{testFsHalfFull})
	}
}

func TestFilterRegexFilter(t *testing.T) {
	filters := Filters{
		IncludeMountPaths: []string{"^/$"},
	}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{testSda}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{testSda})
	}
}

func TestFilterRegexFilterAll(t *testing.T) {
	filters := Filters{
		ExcludeMountPaths: []string{".*"},
	}

	result, err := FilterFileSystem(fileSystemSet, &filters)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, []FilesystemType{}) {
		t.Fatalf("expected %v, got %v", result, []FilesystemType{})
	}
}
