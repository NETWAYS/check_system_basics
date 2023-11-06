# check_system_basics

`check_system_basics` is a monitoring plugin, which is capable to check various Linux metrics such
as memory or filesystem usage via subcommands.

In the current version check_system_basics supports the `memory`, `filesystem`, `psi`, `sensors`, `netdev` and `load` sub command.

## Usage

### memory

A sub command to measure and evaluate memory and swap usage. This is not a trivial topic, but there is [detailed information](https://www.thegeekdiary.com/understanding-proc-meminfo-file-analyzing-memory-utilization-in-linux/) available for those who search for it.

For the memory usage thresholds can be applied to either available, free or used memory. The recommended way is to set thresholds for available memory,
since this is probably the metric most administrators are interested in.


### filesystem

A sub command to check the usage of the currently mounted filesystems.

Thresholds can be applied on absolute values (in bytes) or percentage values of the free space on the filesystem
or the free inodes.

 * With `--exclude-fs-type` and `--include-fs-type` specific filesystem types can be excluded or explicitly included.
 * With `--exclude-device-path` and `--include-device-path` specific device paths can be excluded or explicitly included.
 * With `--exclude-mount-path` and `--include-mount-path` specific mount paths can be excluded or explicitly included.


### load

A sub command to retrieve the current system load values and alerts if they are not within the defined thresholds.

By default no thresholds are applied.


### psi

Note: The Pressure stall information interface is not available on all current Linux distributions (specifically it is not
activated in their kernel configuration).
Therefore this command is not available and will exit with an error.

A sub command to retrieve the current ["pressure stall information"](https://lwn.net/Articles/759781/) values of the system. These are useful metrics to determine a shortage
of resources on the system, the resources being cpu, memory and io.
For each of these resources a 10 second, 60 second and 300 second aggregate percentage value is available for which
amount of time a process did not immediately receive the resource it was asking for and was therefore stalled.

The PSI interface might not be available on your systems, since not all distributions build it into their kernel (RHEL for example).
In this case, the Plugin will return UNKNOWN.

At least on RHEL systems the PSI interface can be enabled via appending "psi=1" to the kernel commandline (`/etc/default/grub`).

The checks includes the three components CPU, IO and Memory by default, but individual components can be selected with the following flagS:

```
--include-cpu
--include-memory
--include-io
```

Default thresholds are applied to all of the measurements.




### sensors

A sub command to read the sensors exposed by the linux kernel and display whether they
are within their respective thresholds (if any are set on the system side).
Additionally it will export the respective values as performance data to be rendered
by a graphing system.

There are no parameters available at the point of writing.


## Building

### Necessary tools

 * the [`golang` toolchain](https://go.dev/)

### Compiling

```
go build
```
executed in the main folder of this repository will generate an executable. If you are not on a linux system,
it should probably look like this:

```
GOOS=linux go build
```

### Creating the Icinga2 CheckCommand config (if necessary)

```
./check_system_basics --dump-icinga2-config > myConfigFile.conf
```
