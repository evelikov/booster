package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var (
	outputFile         = flag.String("output", "booster.img", "Output initrd file")
	forceOverwriteFile = flag.Bool("force", false, "Overwrite existing initrd file")
	initBinary         = flag.String("initBinary", "/usr/lib/booster/init", "Booster 'init' binary location")
	compression        = flag.String("compression", "", `Output file compression ("zstd", "gzip", "none")`)
	kernelVersion      = flag.String("kernelVersion", "", "Linux kernel version to generate initramfs for")
	configFile         = flag.String("config", "", "Configuration file path")
	debugEnabled       = flag.Bool("debug", false, "Enable debug output")
	universal          = flag.Bool("universal", false, "Add wide range of modules/tools to allow this image boot at different machines")
	strip              = flag.Bool("strip", false, "Strip ELF binaries before adding it to the image")
	pprofcpu           = flag.String("pprof.cpu", "", "Write cpu profile to file")
)

func debug(format string, v ...interface{}) {
	if *debugEnabled {
		fmt.Printf(format, v...)
	}
}

func runGenerator() error {
	if *pprofcpu != "" {
		f, err := os.Create(*pprofcpu)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}
		defer pprof.StopCPUProfile()
	}

	file := *configFile
	if file == "" {
		_, err := os.Stat(defaultConfigPath)
		if err == nil {
			file = defaultConfigPath
		} else if !os.IsNotExist(err) {
			// It is OK if the default config is missing. In this case we consider if the default config is empty.
			return err
		}
	}

	conf, err := readGeneratorConfig(file)
	if err != nil {
		return err
	}

	return generateInitRamfs(conf)
}

func main() {
	flag.Parse()

	if err := runGenerator(); err != nil {
		log.Fatal(err)
	}
}
