package lg

import (
	"io/ioutil"
	"os"
	"sort"
)

func addPlaceholderFile(path string) {
	err := ioutil.WriteFile(path+".gitignore", []byte(
		"# Ignore everything in this directory\n*\n# Except this file\n!.gitignore",
	), 0755)
	if err != nil {
		panic(err)
	}
}

func refreshDir(dirname string) {
	err := os.RemoveAll(dirname)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(dirname, 0777)
	if err != nil {
		panic(err)
	}

	addPlaceholderFile(dirname)
}

func oldestFiles(files []os.FileInfo) []string {
	if len(files) < 6 {
		return []string{}
	}

	times := make([]int, len(files))

	for index, file := range files {
		times[index] = int(file.ModTime().Unix())
	}
	sort.Sort(sort.Reverse(sort.IntSlice(times)))
	cutoff := times[5]

	toDelete := []string{}

	for _, file := range files {
		if int(file.ModTime().Unix()) <= cutoff {
			toDelete = append(toDelete, file.Name())
		}
	}

	return toDelete
}

func logFiles(allFiles []os.FileInfo) []os.FileInfo {
	logFiles := []os.FileInfo{}
	for _, file := range allFiles {
		if file.Name() != ".gitignore" {
			logFiles = append(logFiles, file)
		}
	}

	return logFiles
}

func deleteOldestFiles(files []os.FileInfo) {
	for _, name := range oldestFiles(logFiles(files)) {
		err := os.Remove(oldLogDir() + name)

		if err != nil {
			panic(err)
		}
	}
}
