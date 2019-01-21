## 0.0.4 (unreleased)

- Will now clean the expressions before evaluating them

      There can still be legacy formulas that use the old AND, OR and '='
      syntax, the clean method will replace this with:
             AND => &&
             OR => ||
             = => == (only for equality comparison)

- Improved test runner and test coverage (Run `script/test` or `script/cibuild` on CI)
- CLI is a bit more friendly
