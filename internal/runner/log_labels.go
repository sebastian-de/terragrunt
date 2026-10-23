package runner

import (
	"path/filepath"
	"slices"

	"github.com/gruntwork-io/terragrunt/internal/component"
)

// setLogLabels gives each unit the label the `minimal` log path format shows
// for it, unique among the units the run logs for. Excluded units don't log, so
// they take no part in the uniqueness.
func setLogLabels(workingDir string, units []*component.Unit) {
	logged := make([]*component.Unit, 0, len(units))
	paths := make([]string, 0, len(units))

	for _, unit := range units {
		if unit.Excluded() {
			continue
		}

		logged = append(logged, unit)
		paths = append(paths, unit.Path())
	}

	for i, label := range pathLabels(workingDir, paths) {
		logged[i].SetLogLabel(label)
	}
}

// pathLabels returns the label the `minimal` log path format shows for each of
// the given absolute unit paths, in the same order: the last element of the
// path, with as many parent directories as it takes to tell the unit apart from
// the others.
func pathLabels(baseDir string, paths []string) []string {
	// suffixes holds each path's trailing parts, shortest first.
	suffixes := make([][]string, len(paths))

	// counts holds how many paths end in each trailing part.
	counts := make(map[string]int)

	for i, path := range paths {
		suffixes[i] = trailingParts(relOrSelf(baseDir, path))

		for _, part := range suffixes[i] {
			counts[part]++
		}
	}

	labels := make([]string, len(paths))

	for i, parts := range suffixes {
		j := slices.IndexFunc(parts, func(part string) bool { return counts[part] == 1 })
		if j < 0 {
			// The unit's path ends another unit's path, so no part of it is unique.
			j = len(parts) - 1
		}

		labels[i] = parts[j]
	}

	return labels
}

// trailingParts returns the trailing parts of path, shortest first.
func trailingParts(path string) []string {
	var parts []string

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == filepath.Separator && i+1 < len(path) {
			parts = append(parts, path[i+1:])
		}
	}

	return append(parts, path)
}

// relOrSelf returns path relative to baseDir, or path itself when it's outside baseDir.
func relOrSelf(baseDir, path string) string {
	if rel, err := filepath.Rel(baseDir, path); err == nil && filepath.IsLocal(rel) {
		return rel
	}

	return path
}
