## 0.0.4 (unreleased)

- Add support for two new functions: `array` and `eval_array`

      eval_array('a', (1,2,3), 'b', (2,3,4), 'b - a')

      'a'				=> key 1
      '(1,2,3)'	=> array 1
      'b'				=> key 2
      '(2,3,4)		=> array 2 (needs to be same length as array 1)
      'b-a'			=> metaformula

      Will loop over the arrays and apply the formula to each index, so
      in this example would result in:

            (2-1, 3-2,4-3)
            (1,1,1)

      array(1,2,3,4,5)

      A noop function in this context, mainly added for api parity with
      dentaku, so arrays can be explicilty marked as arrays.

      ARRAY(1,2,3) => (1,2,3,4,5)

- Add support for profiling (see README for instructions)
- Will now clean the expressions before evaluating them

      There can still be legacy formulas that use the old AND, OR and '='
      syntax, the clean method will replace this with:
             AND => &&
             OR => ||
             = => == (only for equality comparison)

- Improved test runner and test coverage (Run `script/test` or `script/cibuild` on CI)
- CLI is a bit more friendly
