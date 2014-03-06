package engine

type reportResult struct {
  Passed bool `json:"pass"`
}

type verificationResult struct {
  Id       string `json:"id"`
  ParentId string `json:"parent_id"`
  Passed   bool   `json:"pass"`
}

func newReportResult(pass bool) reportResult {
  return reportResult{Passed: pass}
}

func newVerificationResult(id string, parentId string, pass bool) verificationResult {
  return verificationResult{Id: id, ParentId: parentId, Passed: pass}
}

func (r reportResult) Pass() bool {
  return r.Passed
}

func (r reportResult) Fail() bool {
  return !r.Pass()
}

func (v verificationResult) Pass() bool {
  return v.Passed
}

func (v verificationResult) Fail() bool {
  return !v.Pass()
}
