package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mcmod-update/src/model"
	v1 "mcmod-update/src/repo/curseforge/v1"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func Init(apiKey, version, modLoader, modsPath, recordPath string, optDep bool) {
	modsF, err := os.Open(modsPath)
	if err != nil {
		panic(err)
	}
	defer modsF.Close()

	oldIdFileMap := make(map[int32]*model.File)
	scanner := bufio.NewScanner(modsF)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) < 2 {
			continue
		}

		id, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		oldIdFileMap[int32(id)] = &model.File{
			ModID:     int32(id),
			ModName:   fields[1],
			Date:      time.Now(),
			McVersion: version,
			ModLoader: modLoader,
		}
	}

	idFileMap := update(apiKey, oldIdFileMap, optDep)

	files := make([]*model.File, 0, len(idFileMap))
	downloadFile := make(map[string]*model.File)
	for _, file := range idFileMap {
		downloadFile[file.FileName] = file
		files = append(files, file)
	}

	if !downloadFiles(downloadFile) {
		return
	}

	if err := updateRecordFile(files, recordPath); err != nil {
		panic(err)
	}
}

func CheckUpdate(apiKey, recordPath string, optDep bool) {
	oldRecord, err := os.ReadFile(recordPath)
	if err != nil {
		panic(err)
	}

	files := []*model.File{}
	if err := json.Unmarshal(oldRecord, &files); err != nil {
		panic(err)
	}

	oldIdFileMap := make(map[int32]*model.File)
	for _, file := range files {
		oldIdFileMap[file.ModID] = file
	}

	idFileMap := update(apiKey, oldIdFileMap, optDep)
	updates := make(map[string]*model.File)
	newFiles := make([]*model.File, 0, len(idFileMap))
	for _, nf := range idFileMap {
		// newFiles include files that not exist yet
		newFiles = append(newFiles, nf)
		if len(nf.DownloadUrl) == 0 {
			continue
		}
		if of, ok := oldIdFileMap[nf.ModID]; ok {
			// Old file exists and is up-to-date
			if !nf.Date.After(of.Date) && FileExist(of.FileName) {
				continue
			}
			updates[of.FileName] = nf
		} else {
			updates[nf.FileName] = nf
		}

		fmt.Printf("Find update: %s at %s\n", nf.DispName, nf.Date)
	}

	if len(updates) == 0 {
		fmt.Println("All mods are up-to-date")
		return
	} else {
		fmt.Printf("%d/%d mod(s) have updates\n", len(updates), len(newFiles))
	}

	if !downloadFiles(updates) {
		return
	}

	if err = updateRecordFile(newFiles, recordPath); err != nil {
		panic(err)
	}
}

func update(apiKey string, fileMap map[int32]*model.File, optDep bool) map[int32]*model.File {
	mu := sync.Mutex{}
	newFileMap := make(map[int32]*model.File)
	notFoundCnt := atomic.Int64{}
	wg := new(sync.WaitGroup)
	f := v1.NewAdaptor(apiKey)
	for modId, file := range fileMap {
		wg.Add(1)
		go func(modId int32, file *model.File) {
			defer wg.Done()
			files, err := f.GetLatestModFileWithDeps(modId, file.McVersion,
				file.ModLoader, optDep)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			if len(files) == 0 {
				log.Printf("No mod file found, name: %q, id: %d, ver: %s %s\n",
					file.ModName, modId, file.McVersion, file.ModLoader)
				notFoundCnt.Add(1)
				files = []*model.File{file}
			}

			mu.Lock()
			for _, f := range files {
				if f.ModID == modId {
					f.ModName = file.ModName
				}
				newFileMap[f.ModID] = f
			}
			mu.Unlock()
		}(modId, file)
	}
	wg.Wait()

	if notFound := notFoundCnt.Load(); notFound > 0 {
		fmt.Printf("%d/%d mod(s) not found\n", notFound, len(newFileMap))
	}

	return newFileMap
}

func updateRecordFile(files []*model.File, recordPath string) error {
	tmpFile := recordPath + ".new"
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	if err = recordUpdate(files, f); err != nil {
		return err
	}

	f.Close()
	if err = os.Rename(tmpFile, recordPath); err != nil {
		return err
	}

	return nil
}

func recordUpdate(files []*model.File, record io.Writer) error {
	sort.Sort(model.FileSliceSortByModId(files))
	b, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return err
	}

	if _, err = record.Write(b); err != nil {
		return err
	}

	return nil
}

func downloadFiles(fileMap map[string]*model.File) bool {
	fmt.Printf("Download updates?[y/N]: ")
	input := ""
	if _, err := fmt.Scanln(&input); err != nil {
		return false
	}

	if !yes(input) {
		return false
	}

	wg := new(sync.WaitGroup)
	for oldName, file := range fileMap {
		wg.Add(1)
		go func(oldName string, file *model.File) {
			defer wg.Done()

			if len(file.DownloadUrl) == 0 {
				return
			}

			if FileExist(file.FileName) {
				fmt.Printf("Skip: file %q exist\n", file.FileName)
				return
			}

			os.Remove(oldName)

			log.Printf("Downloading %s\n", file.FileName)
			if err := Get(file.DownloadUrl, file.FileName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}(oldName, file)
	}
	wg.Wait()

	return true
}

func yes(input string) bool {
	return input == "y" || input == "Y"
}
