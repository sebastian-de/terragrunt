package runner

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terragrunt/internal/component"
	"github.com/gruntwork-io/terragrunt/test/helpers/venvtest"

	"github.com/stretchr/testify/assert"
)

func TestPathLabels(t *testing.T) {
	t.Parallel()

	base := venvtest.Root("/work/live")
	outside := venvtest.Root("/other")

	stacked := func(elems ...string) string {
		path := base
		for _, elem := range elems {
			path = filepath.Join(path, ".terragrunt-stack", elem)
		}

		return path
	}

	tests := []struct {
		expected []string
		name     string
		paths    []string
	}{
		{
			name:  "unique_names_keep_the_last_element",
			paths: []string{stacked("openebs", "mayastor-diskpools"), filepath.Join(base, "vpc")},
			expected: []string{
				"mayastor-diskpools",
				"vpc",
			},
		},
		{
			name:  "same_names_add_parents_with_stack_dirs",
			paths: []string{stacked("app-a", "db"), stacked("app-b", "db"), stacked("app-a", "vpc")},
			expected: []string{
				filepath.Join("app-a", ".terragrunt-stack", "db"),
				filepath.Join("app-b", ".terragrunt-stack", "db"),
				"vpc",
			},
		},
		{
			name:  "parents_are_added_until_unique",
			paths: []string{filepath.Join(base, "eu", "app", "db"), filepath.Join(base, "us", "app", "db")},
			expected: []string{
				filepath.Join("eu", "app", "db"),
				filepath.Join("us", "app", "db"),
			},
		},
		{
			name:  "name_ending_another_falls_back_to_full_path",
			paths: []string{filepath.Join(base, "db"), filepath.Join(base, "app", "db")},
			expected: []string{
				"db",
				filepath.Join("app", "db"),
			},
		},
		{
			name:  "stack_dir_tells_units_apart",
			paths: []string{stacked("app", "db"), filepath.Join(base, "app", "db")},
			expected: []string{
				filepath.Join(".terragrunt-stack", "db"),
				filepath.Join("app", "db"),
			},
		},
		{
			name:  "paths_outside_base_dir_use_their_own_parents",
			paths: []string{filepath.Join(outside, "db"), filepath.Join(base, "db")},
			expected: []string{
				filepath.Join(filepath.Base(outside), "db"),
				"db",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, pathLabels(base, tc.paths))
		})
	}
}

// TestSetLogLabels checks setLogLabels labels the units a run logs for, and
// that excluded units neither get a label nor lengthen the others'.
func TestSetLogLabels(t *testing.T) {
	t.Parallel()

	base := venvtest.Root("/work/live")

	excluded := component.NewUnit(filepath.Join(base, "app-a", "db"))
	excluded.SetExcluded(true)

	units := []*component.Unit{
		component.NewUnit(filepath.Join(base, "app-a", "db")),
		component.NewUnit(filepath.Join(base, "app-b", "db")),
		component.NewUnit(filepath.Join(base, "vpc")),
		excluded,
	}

	setLogLabels(base, units)

	assert.Equal(t, filepath.Join("app-a", "db"), units[0].LogLabel())
	assert.Equal(t, filepath.Join("app-b", "db"), units[1].LogLabel())
	assert.Equal(t, "vpc", units[2].LogLabel())
	assert.Empty(t, excluded.LogLabel(), "an excluded unit is not logged for")
}
