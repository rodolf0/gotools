# bisect: a simple array bisection algorithm

Source *http://docs.python.org/library/bisect.html*

## Example

    package main

    import (
      "bisect"
      "fmt"
    )

    type Score int

    func (i Score) Less(other bisect.Elem) bool {
      return i < other.(Score)
    }

    func main() {
      var breakpoints = []bisect.Elem{Score(60), Score(70), Score(80), Score(90)}
      var grades = "FDCBA"

      for _, test_score := range []Score{93, 3, 83, 61, 72, 55, 99, 110, 100} {
        var grade_idx = bisect.Bisect(breakpoints, test_score)
        fmt.Printf("Score %v is graded with %c\n", test_score, grades[grade_idx])
      }
    }
