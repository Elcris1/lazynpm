package app

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type InfoKey string
const (
	Name InfoKey = "name"
	Version InfoKey = "version"
	FileLocation InfoKey = "fileLocation"
	PackageManager InfoKey = "packageManager"
)
type Package struct {
	Info PacakgeEntry
	Scripts PacakgeEntry
	Dependecies PacakgeEntry
	DevDependencies PacakgeEntry
}

type PacakgeEntry struct {
	Keys []string
	Data map[string]string
}

type PackageInfo struct {
	FileLocation    string
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	PackageManager  string            `json:"packageManager"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func parseMapToEntry(m map[string]string) PacakgeEntry {
	if m == nil {
		return PacakgeEntry{
			Keys: []string{},
			Data: make(map[string]string),
		}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return PacakgeEntry{
		Keys: keys,
		Data: m,
	}
}

func ParsePackageInfo(info *PackageInfo) Package {
	if info == nil {
		return Package{}
	}

	infoData := map[string]string {
		string(FileLocation): info.FileLocation,
		string(Name): info.Name,
		string(Version): info.Version,
		string(PackageManager): info.PackageManager,
	}

	return Package {
		Info: parseMapToEntry(infoData),
		Scripts: parseMapToEntry(info.Scripts),
		Dependecies: parseMapToEntry(info.Dependencies),
		DevDependencies: parseMapToEntry(info.DevDependencies),
	}
	
}

func LoadPackage(filename string, pkg *Package)  error {
	data, err := os.ReadFile(filename)
	if err != nil {
		//%q stands for double quoted string
		//%w stand for wrpaped error format
		return fmt.Errorf("loadPackage: failed to read file: %q, %w", filename, err)
	}
	info := &PackageInfo{
		FileLocation: filename,
	}
	if err := json.Unmarshal(data, info); err != nil {
		return fmt.Errorf("loadPackage: failed to parse JSON in %q, %w", filename, err)
	}
	*pkg = ParsePackageInfo(info)
	return nil
}

func (p PackageInfo) String() string {
	content := fmt.Sprintf("FileLocation: %s\nPacageInfo:\n\t- Name: %s\n\t- Version: %s\n\t- PackageManager: %s\n", p.FileLocation, p.Name, p.Version, p.PackageManager)
	appendSection(&content, "Scripts", p.Scripts)
	appendSection(&content, "Dependencies", p.Dependencies)
	appendSection(&content, "DevDependencies", p.DevDependencies)
	return content

}

func appendSection(buffer *string, section string, content map[string]string) {
	if len(content) == 0 {
		*buffer += fmt.Sprintf("%s: no content\n", section)
	} else {
		*buffer += fmt.Sprintf("%s:\n", section)
	}
	for key, val := range content {
		*buffer += fmt.Sprintf("\t- %q: %q\n", key, val)
	}

}