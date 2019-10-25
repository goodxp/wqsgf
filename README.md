# wqSGF

wqSGF is a SGF parser for Go game(weiqi).
  - load and save .sgf files (from/to game tree structure for further coding)
  - encode/decode SGF strings
  - helper functions for convertion between SGF property value types and Go types

This package was designed with care about NOT to interfere with the user's code design as much as possible. Even a simple tree structure is provided as a base to parse SGF, though the user would find it is easy to alter it with their own design. Therefore NO concrete Go types defined for corresponding SGF property value types. Instead, convertion helpers are provided to make it easy to deal with SGF values. With that said, the tree structure provided in the package is simple and to the point, a good one to start with.
