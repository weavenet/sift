package plan

import (
  "encoding/json"
  log "github.com/cihub/seelog"
  "io/ioutil"
  "os"
  "reflect"
  "testing"
)

func TestLoadRepo(t *testing.T) {
  defer log.Flush()
  wd, _ := os.Getwd()

  var plans = []string{"minimal", "full", "dynamic_filter"}

  for _, n := range plans {
    t.Logf("Testing loading repo from plan '%s'", n)
    repoPlan := NewPlan()
    jsonPlan := NewPlan()

    if err := repoPlan.LoadRepo(wd + "/test/repo_to_plan/valid/" + n + "/repo"); err != nil {
      t.Fatalf("Received error loading repo '%s'.", err)
    }
    jsonData, _ := ioutil.ReadFile(wd + "/test/repo_to_plan/valid/" + n + "/plan.json")
    if err := jsonPlan.LoadJSON(jsonData); err != nil {
      t.Fatalf("Received error loading plan from file '%s'.", err)
    }
    if !reflect.DeepEqual(repoPlan.Accounts, jsonPlan.Accounts) {
      t.Fatalf("Accounts '%s' does not match '%s'", repoPlan.Accounts, jsonPlan.Accounts)
    }
    if !reflect.DeepEqual(repoPlan.Lists, jsonPlan.Lists) {
      t.Fatalf("Lists '%s' does not match '%s'.", repoPlan.Lists, jsonPlan.Lists)
    }
    if !reflect.DeepEqual(repoPlan.Policies, jsonPlan.Policies) {
      t.Fatalf("Policies '%s' does not match '%s'.", repoPlan.Policies, jsonPlan.Policies)
    }
    if !reflect.DeepEqual(repoPlan.Sources, jsonPlan.Sources) {
      t.Fatalf("Sources '%s' does not match '%s'.", repoPlan.Sources, jsonPlan.Sources)
    }
    if !reflect.DeepEqual(repoPlan.Filters, jsonPlan.Filters) {
      t.Fatalf("Filters '%s' does not match '%s'.", repoPlan.Sources, jsonPlan.Sources)
    }

  }
}

func TestPlanErrorsLoadingRepo(t *testing.T) {
  defer log.Flush()
  wd, _ := os.Getwd()
  plans := map[string]string{
    "no_accounts":    "No accounts found in repo.",
    "no_account_dir": "Accounts directory does not exist in repo.",
    "foo":            "Repo directory does not exist.",
  }

  for plan, e := range plans {
    t.Logf("Testing in-valid plan %s", plan)
    repoPlan := NewPlan()

    if err := repoPlan.LoadRepo(wd + "/test/repo_to_plan/invalid/" + plan + "/repo"); err.Error() != e {
      t.Fatalf("Failed to receive error '%s' got '%s'", e, err)
    }
  }
}

func TestCreateEvaluations(t *testing.T) {
  defer log.Flush()
  wd, _ := os.Getwd()

  var plans = []string{"full", "minimal", "dynamic_filter1", "dynamic_filter2", "dynamic_filter3", "multi_source"}

  for _, n := range plans {
    t.Logf("Testing creating evaluations from plan '%s'", n)

    p := NewPlan()

    jsonData, _ := ioutil.ReadFile(wd + "/test/plan_to_evaluations/valid/" + n + "/plan.json")
    if err := p.LoadJSON(jsonData); err != nil {
      t.Fatalf("Received error loading plan from file '%s'.", err)
    }

    e, _ := p.Evaluations()
    planEvaluations, err := json.Marshal(e)
    if err != nil {
      t.Fatalf("Error loading evaluations.")
    }
    jsonEvaluations, _ := ioutil.ReadFile(wd + "/test/plan_to_evaluations/valid/" + n + "/evaluations.json")

    var pe interface{}
    if err := json.Unmarshal(planEvaluations, &pe); err != nil {
      t.Fatalf("Error converting evaluations from repo to JSON '%s'.", err)
    }

    var je interface{}
    if err := json.Unmarshal(jsonEvaluations, &je); err != nil {
      t.Fatalf("Error converting evaluation.json file to JSON '%s'.", err)
    }

    if !reflect.DeepEqual(pe, je) {
      t.Fatalf("Repo converted to evaluations \n '%s' \n does not match \n '%s'.", planEvaluations, jsonEvaluations)
    }
  }
}
