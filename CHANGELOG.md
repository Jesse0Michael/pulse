## [1.0.6](https://github.com/Jesse0Michael/pulse/compare/v1.0.5...v1.0.6) (2023-10-18)

### Bug Fixes

- add summary endpoint implementation ([dbc5101](https://github.com/Jesse0Michael/pulse/commit/dbc510185147a520dbcc17dc01d0a384c73b2b64))

### Chores

- address linter issues ([c221064](https://github.com/Jesse0Michael/pulse/commit/c2210641814b9fbf9ec88b052ff38ca021919fcb))

### Continuous Integration

- add linter to configuration ([ddcbf23](https://github.com/Jesse0Michael/pulse/commit/ddcbf23b0fdd1a5c16a600fc1c76cb68d9eab58c))
- lint with GO ([21fa359](https://github.com/Jesse0Michael/pulse/commit/21fa359575c7568772df499764c286427e53e111))
- separate linting and testing jobs ([885e7c8](https://github.com/Jesse0Michael/pulse/commit/885e7c892567e3667d82d0ec06a891296f01f376))
- update test CI ([e35f973](https://github.com/Jesse0Michael/pulse/commit/e35f9737b53413b5f5e388fa5ba15386df19ffcb))
- validate all codebase ([681f60c](https://github.com/Jesse0Michael/pulse/commit/681f60ca7992825d016f68df4975cc0f47b2c10c))

## [1.0.5](https://github.com/Jesse0Michael/pulse/compare/v1.0.4...v1.0.5) (2023-10-18)

### Bug Fixes

- update text fixtures ([6bdda6d](https://github.com/Jesse0Michael/pulse/commit/6bdda6de58042c541f937474902457cd75b68d20))

## [1.0.4](https://github.com/Jesse0Michael/pulse/compare/v1.0.3...v1.0.4) (2023-10-17)

### Bug Fixes

- update default server timeout ([e7d8bf3](https://github.com/Jesse0Michael/pulse/commit/e7d8bf30a90b1f8200bed3c012a72cba0083bde6))

### Documentation

- change README.md syntax ([9f7ea7c](https://github.com/Jesse0Michael/pulse/commit/9f7ea7c2dd6d5b3f3556cd5e929b17976801be40))

### Tests

- add behavioral tests for the pulse cli ([8f6d16d](https://github.com/Jesse0Michael/pulse/commit/8f6d16da842b2361c05bb3aec77764c764270ae1))

## [1.0.3](https://github.com/Jesse0Michael/pulse/compare/v1.0.2...v1.0.3) (2023-10-16)

### Documentation

- update README with application purpose, usage, and examples ([d62e9e9](https://github.com/Jesse0Michael/pulse/commit/d62e9e9a0d10e052b0bc8a22b2d40712622d910d))

### Other

- Merge branch 'main' of ssh://github.com/Jesse0Michael/pulse ([e3d9e20](https://github.com/Jesse0Michael/pulse/commit/e3d9e205bf86c74669608d8cc88354322072b32d))

## [1.0.2](https://github.com/Jesse0Michael/pulse/compare/v1.0.1...v1.0.2) (2023-10-16)

### Bug Fixes

- use correct directory to build cli ([fb211cf](https://github.com/Jesse0Michael/pulse/commit/fb211cf72ab8724a8ebce7cf1f935e5cedd8694c))

### Chores

- remove debug line ([23a7c54](https://github.com/Jesse0Michael/pulse/commit/23a7c54fce8a6c2394561345139ccd103d4d692f))

### Other

- Merge branch 'main' of ssh://github.com/Jesse0Michael/pulse ([9c90c8f](https://github.com/Jesse0Michael/pulse/commit/9c90c8f2718763672fcb988a84a2ae7145439593))

## [1.0.1](https://github.com/Jesse0Michael/pulse/compare/v1.0.0...v1.0.1) (2023-10-16)

### Chores

- add Environment configuration to CLI usage ([d4fa0e9](https://github.com/Jesse0Michael/pulse/commit/d4fa0e9cf50290fad89ede172c687b20b1e696bb))
- rename pulse cmd directory ([cb09a8b](https://github.com/Jesse0Michael/pulse/commit/cb09a8bcfee80771882ec7e3bf451c9a8b12f921))

### Other

- Merge branch 'main' of ssh://github.com/Jesse0Michael/pulse ([0b69300](https://github.com/Jesse0Michael/pulse/commit/0b6930019384a267bf68ff16695228dcfcfd9fc8))

# 1.0.0 (2023-10-16)

### Bug Fixes

- add initial openAI client ([a090eee](https://github.com/Jesse0Michael/pulse/commit/a090eeeea66da444708dc4fcf7fe6f2eb6747471))
- clear default github service url ([ab0bc62](https://github.com/Jesse0Michael/pulse/commit/ab0bc62c7bbbd2e51ecd5079fb5e8486cd7f5a99))
- make github service testable ([f0f17d7](https://github.com/Jesse0Michael/pulse/commit/f0f17d7f27bb550c32651872df4c06107d74304d))
- openai service summary output formatting ([b9fbbf1](https://github.com/Jesse0Michael/pulse/commit/b9fbbf145aae7b8d6b8f8b137aff51a65b5d61a2))
- parse and format github activity for AI prompt ([2bbb191](https://github.com/Jesse0Michael/pulse/commit/2bbb191d4eecafc157b19cff4548cb6a887de636))
- refactor code into Pulser service struct to use in both the cli and api ([099f690](https://github.com/Jesse0Michael/pulse/commit/099f690d1d2f341d9fef3857560b2db71b815fea))
- reorganize code to service package ([baca0e2](https://github.com/Jesse0Michael/pulse/commit/baca0e209ce9203efdd1384f4e1805111743811b))

### Continuous Integration

- tag repository merges ([8d73d99](https://github.com/Jesse0Michael/pulse/commit/8d73d999173c438336f4827d3fe564bb23fb4178))
- use standard GITHUB_TOKEN ([5744ab1](https://github.com/Jesse0Michael/pulse/commit/5744ab166b1bfff9230698f016b7324937340b11))

### Documentation

- add Makefile ([8fd6b41](https://github.com/Jesse0Michael/pulse/commit/8fd6b41e669f622426bcda34c048db51b528254a))

### Features

- add pulse cli skeleton ([96ea854](https://github.com/Jesse0Michael/pulse/commit/96ea8542dcc9dae86580b474d55d60a0465fc637))
- init pulse repository with api skeleton ([a225778](https://github.com/Jesse0Michael/pulse/commit/a2257788fcc00cefbfd6f839befa1382676d9691))
- update api route to use github and username ([ef3b48a](https://github.com/Jesse0Michael/pulse/commit/ef3b48a7e5e3a7f33a3df3a455d7b4251fd0a97b))

### Tests

- test openai service ([d686c8c](https://github.com/Jesse0Michael/pulse/commit/d686c8c63b574e653d2ab398869d7ad9150add3c))
