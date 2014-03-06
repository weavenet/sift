package engine

import (
  "reflect"
  "testing"
)

var lineageTCs = []struct {
  name      string
  hasParent bool
  lineage   []string
  parent    string
}{
  {"instance", false, []string{"instance"}, ""},
  {"bucket", false, []string{"bucket"}, ""},
  {"bucket-object", true, []string{"bucket", "bucket-object"}, "bucket"},
  {"bucket-object-policy", true, []string{"bucket", "bucket-object", "bucket-object-policy"}, "bucket-object"},
}

func TestLineage(t *testing.T) {
  for _, tc := range lineageTCs {
    t.Logf("Test case '%s'.", tc.name)
    c := newCollection(tc.name)
    if c.hasParent() != tc.hasParent {
      t.Fatalf("Error checking for parent.")
    }
    if !reflect.DeepEqual(c.lineage(), tc.lineage) {
      t.Fatalf("Error testing lineage. Expected '%+v', got '%+v'", c.lineage(), tc.lineage)
    }
    if c.hasParent() {
      tcParent := newCollection(tc.parent)
      if !reflect.DeepEqual(c.parent(), tcParent) {
        t.Fatalf("Error retrieving parent. Expected '%+v', got '%+v'", tcParent, c.parent())
      }
    }
  }
}
