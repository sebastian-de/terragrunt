package options

// UnitLabelKeyName is the log field that holds a unit's [UnitLabel].
const UnitLabelKeyName = "unit-label"

// UnitLabel is the label the `minimal` path format shows for a unit, together
// with the two paths the unit can be rendered under: Prefix is the unit's path
// relative to the discovery working dir, Path is its absolute path, as used by
// `--log-show-abs-paths`. Code that logs with another prefix, e.g. while reading
// a dependency's config, keeps the field, so the label only applies while the
// rendered prefix is still the unit's own.
type UnitLabel struct {
	Prefix string
	Path   string
	Label  string
}
