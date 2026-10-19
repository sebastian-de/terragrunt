package runner

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/gruntwork-io/terragrunt/internal/hclparse"
)

// UnitLabels returns the label the `basename` log path format shows for each
// unit at the given paths, keyed by the unit's absolute path: the last element
// of the path, with as many parent directories as it takes to tell the unit
// apart from the others. Relative paths are resolved against baseDir, and the
// labels never reach above it. Stack directories are skipped, because every
// generated unit has them.
func UnitLabels(baseDir string, paths []string) map[string]string {
	// suffixes holds each unit's trailing parts, shortest first.
	suffixes := make(map[string][]string, len(paths))

	for _, path := range paths {
		if !filepath.IsAbs(path) {
			path = filepath.Join(baseDir, path)
		}

		path = filepath.Clean(path)
		suffixes[path] = trailingParts(relOrSelf(baseDir, path))
	}

	// counts holds how many units end in each trailing part.
	counts := make(map[string]int)

	for _, parts := range suffixes {
		for _, part := range parts {
			counts[part]++
		}
	}

	labels := make(map[string]string, len(suffixes))

	for path, parts := range suffixes {
		if i := slices.IndexFunc(parts, func(part string) bool { return counts[part] == 1 }); i >= 0 {
			labels[path] = parts[i]

			continue
		}

		// The unit's path without stack directories matches another unit's, or
		// ends it. Its path in full, stack directories and all, is unique.
		labels[path] = relOrSelf(baseDir, path)
	}

	return labels
}

// trailingParts returns the trailing parts of path, shortest first, without
// stack directories.
func trailingParts(path string) []string {
	var elems []string

	for elem := range strings.SplitSeq(path, string(filepath.Separator)) {
		if elem != "" && elem != "." && elem != hclparse.StackDir {
			elems = append(elems, elem)
		}
	}

	parts := make([]string, len(elems))

	for k := range elems {
		parts[k] = strings.Join(elems[len(elems)-1-k:], string(filepath.Separator))
	}

	return parts
}

// relOrSelf returns path relative to baseDir, or path itself when it's outside
// baseDir.
func relOrSelf(baseDir, path string) string {
	if rel, err := filepath.Rel(baseDir, path); err == nil && filepath.IsLocal(rel) {
		return rel
	}

	return path
}
