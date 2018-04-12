# Requirements

## v1.0.0

1. Have a CLI.
2. Accept a Go-package URL and a comma-separated list of `<OLD_TYPE>=<NEW_TYPE>`
   pairs.
3. Fail if any `<OLD_TYPE>` doesn't exist as a type-alias to `interface{}` in
   the package.
4. Fail if any `<OLD_TYPE>` is a duplicate.
5. Download the Go-files in the package into a directory relative to where the
   CLI was ran named 'internal/\<NEW_TYPES\>/\<PACKAGE_DIRECTORY\>/' where
   `<NEW_TYPES>` are all `<NEW_TYPE>` values made lower-case and concatenated
   and `<PACKAGE_DIRECTORY>` is the package's leaf-directory.
6. Replace all occurrences of `<OLD_TYPE>` with `<NEW_TYPE>` in the package, for
   each pair.
7. Remove all `<OLD_TYPE>` type-aliases.
8. Rewrite all comments referencing `<OLD_TYPE>` to reference the corresponding
   `<NEW_TYPE>`.
9. Format the package.
10. Include any extra imports.
11. Fail if the the package doesn't build.
