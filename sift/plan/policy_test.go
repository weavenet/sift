package plan

import (
  "fmt"
  log "github.com/cihub/seelog"
  "reflect"
  "testing"
)

var matrixTCs = []struct {
  args   map[string][]string
  result []map[string]string
}{
  {
    map[string][]string{
      "arg1": []string{"a"},
      "arg2": []string{"alpha", "beta"},
    },
    []map[string]string{
      map[string]string{"arg1": "a", "arg2": "alpha"},
      map[string]string{"arg1": "a", "arg2": "beta"},
    },
  },
  {
    map[string][]string{
      "arg1": []string{"a", "b"},
      "arg2": []string{"alpha", "beta"},
    },
    []map[string]string{
      map[string]string{"arg1": "a", "arg2": "alpha"},
      map[string]string{"arg1": "a", "arg2": "beta"},
      map[string]string{"arg1": "b", "arg2": "alpha"},
      map[string]string{"arg1": "b", "arg2": "beta"},
    },
  },
  {
    map[string][]string{
      "arg1": []string{"a", "b"},
      "arg2": []string{"1"},
      "arg3": []string{"alpha", "beta"},
    },
    []map[string]string{
      map[string]string{"arg1": "a", "arg2": "1", "arg3": "alpha"},
      map[string]string{"arg1": "a", "arg2": "1", "arg3": "beta"},
      map[string]string{"arg1": "b", "arg2": "1", "arg3": "alpha"},
      map[string]string{"arg1": "b", "arg2": "1", "arg3": "beta"},
    },
  },
  {
    map[string][]string{
      "arg1": []string{"a", "b"},
      "arg2": []string{"1", "2", "3"},
      "arg3": []string{"alpha", "beta"},
    },
    []map[string]string{
      map[string]string{"arg1": "a", "arg2": "1", "arg3": "alpha"},
      map[string]string{"arg1": "a", "arg2": "1", "arg3": "beta"},
      map[string]string{"arg1": "a", "arg2": "2", "arg3": "alpha"},
      map[string]string{"arg1": "a", "arg2": "2", "arg3": "beta"},
      map[string]string{"arg1": "a", "arg2": "3", "arg3": "alpha"},
      map[string]string{"arg1": "a", "arg2": "3", "arg3": "beta"},
      map[string]string{"arg1": "b", "arg2": "1", "arg3": "alpha"},
      map[string]string{"arg1": "b", "arg2": "1", "arg3": "beta"},
      map[string]string{"arg1": "b", "arg2": "2", "arg3": "alpha"},
      map[string]string{"arg1": "b", "arg2": "2", "arg3": "beta"},
      map[string]string{"arg1": "b", "arg2": "3", "arg3": "alpha"},
      map[string]string{"arg1": "b", "arg2": "3", "arg3": "beta"},
    },
  },
}

func TestBuildArgumentMatrix(t *testing.T) {
  for count, tc := range matrixTCs {
    t.Logf("Testing matric #%d.", count)
    matrix := buildArgumentMatrix(tc.args)
    logConfig := fmt.Sprintf(`<seelog type="sync">`)
    logger, _ := log.LoggerFromConfigAsBytes([]byte(logConfig))
    log.ReplaceLogger(logger)
    if !reflect.DeepEqual(matrix, tc.result) {
      t.Errorf("Error buidling matrix from %+v", tc.args)
      t.Fatalf("Error buidling arument matrix. Expected '%+v', got '%+v'.", tc.result, matrix)
    }
  }
}
