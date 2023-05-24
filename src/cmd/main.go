package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mcmod-update/src/repo"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ApiKeyEnv   string
	McVersion   string
	CheckUpdate string
	Download    bool
	Fabric      bool
)

func init() {
	flag.StringVar(&ApiKeyEnv, "k",
		"CurseForgeAPIKey",
		"Environment variable of CurseForge api key")
	flag.StringVar(&McVersion, "v", "1.19.2", "minecraft version")
	flag.StringVar(&CheckUpdate, "c", "", "Check for mod update")
	flag.BoolVar(&Download, "d", false, "Download updated mod and update file")
	flag.BoolVar(&Fabric, "fabric", false,
		"Use fabric instead of forge as mod loader")
}

func main() {
	flag.Parse()
	if len(ApiKeyEnv) == 0 {
		flag.Usage()
		return
	}

	ApiKey := os.Getenv(ApiKeyEnv)
	if len(ApiKey) == 0 {
		fmt.Printf("Env %q not found\n", ApiKeyEnv)
		return
	}

	modLoader := repo.ModLoaderForge
	if Fabric {
		modLoader = repo.ModLoaderFabric
	}
	cf := repo.NewCurseforgeRepo(ApiKey)
	if len(CheckUpdate) > 0 {
		checkUpdate(cf, CheckUpdate, Download)
		return
	}

	ver := repo.NewVersion(McVersion, modLoader)
	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	mods := make([]*repo.ModFile, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) != 2 {
			continue
		}
		id, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			mod, err := cf.LatestModFile(id, ver)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			if mod == nil {
				log.Printf("No mod file found, name: %q, id: %d, ver: %s %s\n",
					name, id, ver.McVersion, ver.ModLoader)
				mod = &repo.ModFile{
					ModId:     id,
					Date:      time.Now(),
					McVersion: ver.McVersion,
					ModLoader: modLoader,
				}
			}

			mod.ModName = name
			mu.Lock()
			mods = append(mods, mod)
			mu.Unlock()

			if Download && !fileExist(mod.FileName) {
				fmt.Printf("Downloading %s\n", mod.FileName)
				get(mod.Url)
			}
		}(fields[1])
	}
	wg.Wait()
	sort.Sort(repo.ModFileSlice(mods))
	js, err := json.Marshal(mods)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(js))
}

func checkUpdate(cf *repo.CurseforgeRepo, path string, download bool) {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	mods := []*repo.ModFile{}
	if err := json.Unmarshal(f, &mods); err != nil {
		panic(err)
	}

	updateCnt := atomic.Int64{}
	notFoundCnt := atomic.Int64{}
	wg := new(sync.WaitGroup)
	for i, mod := range mods {
		wg.Add(1)
		go func(index int, mod *repo.ModFile) {
			defer wg.Done()
			m, err := cf.LatestModFile(mod.ModId, repo.NewVersion(mod.McVersion, mod.ModLoader))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			if m == nil {
				log.Printf("No mod file found, name: %q, id: %d, ver: %s %s\n",
					mod.ModName, mod.ModId, mod.McVersion, mod.ModLoader)
				m = mod
				notFoundCnt.Add(1)
			}

			if !m.Date.After(mod.Date) && fileExist(m.FileName) {
				return
			}
			updateCnt.Add(1)
			fmt.Printf("Find update: %s at %s\n", m.DispName, m.Date)
			mods[index] = m
			if download && !fileExist(m.FileName) {
				fmt.Printf("Downloading %s\n", mod.FileName)
				os.Remove(mod.FileName)
				get(m.Url)
			}
		}(i, mod)
	}

	wg.Wait()

	if notFoundCnt.Load() > 0 {
		fmt.Printf("%d/%d mod(s) not found\n", notFoundCnt.Load(), len(mods))
	}

	if updateCnt.Load() == 0 {
		fmt.Println("All mods are up-to-date")
	} else {
		fmt.Printf("%d/%d mod(s) have updated\n", updateCnt.Load(), len(mods))
	}

	if download {
		js, err := json.Marshal(mods)
		if err != nil {
			panic(err)
		}
		nf, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		nf.Write(js)
		nf.Close()
	}
}

func get(url string) {
	if len(url) == 0 {
		return
	}
	content, err := repo.NewRequest("GET", url).Do()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	paths := strings.Split(url, "/")
	f, err := os.Create(paths[len(paths)-1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	f.Write(content)
}

// Check if a file exist
func fileExist(path string) bool {
	if len(path) == 0 {
		return true
	}
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
