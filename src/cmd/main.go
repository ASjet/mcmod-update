package main

import (
	"flag"
	"fmt"
	"mcmod-update/src"
	"os"
)

var (
	ApiKeyEnv string
	McVersion string
	ModLoader string
	Init      string
	Record    string
	OptDep    bool
)

func init() {
	flag.StringVar(&ApiKeyEnv, "k",
		"CURSE_FORGE_APIKEY",
		"Environment variable of CurseForge api key")
	flag.StringVar(&McVersion, "v", "1.19.2", "minecraft version")
	flag.StringVar(&ModLoader, "l", "forge", "mod loader")
	flag.StringVar(&Init, "i", "", "Initialize mod records")
	flag.StringVar(&Record, "r", "", "mod record file")
	flag.BoolVar(&OptDep, "o", false, "download optional dependencies")
}

func main() {
	flag.Parse()
	if len(ApiKeyEnv) == 0 {
		flag.Usage()
		return
	}

	apiKey := os.Getenv(ApiKeyEnv)
	if len(apiKey) == 0 {
		fmt.Printf("Env %q not found\n", ApiKeyEnv)
		return
	}

	if len(Init) == 0 {
		src.CheckUpdate(apiKey, Record, OptDep)
		return
	}

	src.Init(apiKey, McVersion, ModLoader, Init, Record, OptDep)
}
