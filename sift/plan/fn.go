package plan

import (
  "fmt"
  log "github.com/cihub/seelog"
)

func executeFn(data map[string]interface{}, p plan) ([]string, error) {
  var r []string

  if len(data) > 1 {
    return r, fmt.Errorf("Fns must be single key value pair.")
  }

  var name string
  var args []string

  for n, a := range data {
    if w, ok := a.([]interface{}); ok {
      log.Debugf("Fn '%s' with arguments '%a'.", n, a)
      name = n
      for _, i := range w {
        args = append(args, i.(string))
      }
    } else {
      return r, fmt.Errorf("Fn expects an array of strings. Received '%s'.", a)
    }
  }

  switch name {
  case "Fn::List":
    {
      return listFn(args, p)
    }
  case "Fn::ListSub":
    {
      return listSubFn(args, p)
    }
  case "Fn::ListOnly":
    {
      return listOnlyFn(args, p)
    }
  default:
    {
      return r, fmt.Errorf("Fn '%s' is not a valid fn.", name)
    }
  }
}

func listFn(args []string, p plan) ([]string, error) {
  if len(args) != 1 {
    return []string{}, fmt.Errorf("Fn::List expects single list name as argument. Received '%s'.", args)
  }
  listName := args[0]

  if list, ok := p.Lists[listName]; ok {
    log.Debugf("Returning entries from list '%s'.", listName)
    return list.all(), nil
  }
  return []string{}, fmt.Errorf("List '%s' not found by Fn::List.", listName)
}

func listSubFn(args []string, p plan) ([]string, error) {
  if len(args) != 2 {
    return []string{}, fmt.Errorf("Fn::ListSub expects list name and tag to override default (when present) as arguments. Received '%s'.", args)
  }
  listName := args[0]
  tag := args[1]

  if list, ok := p.Lists[listName]; ok {
    log.Debugf("Returning entries from list '%s' substituting tag '%s' when present.", listName, tag)
    return list.sub(tag), nil
  }
  return []string{}, fmt.Errorf("List '%s' not found by Fn::List.", listName)
}

func listOnlyFn(args []string, p plan) ([]string, error) {
  if len(args) != 2 {
    return []string{}, fmt.Errorf("Fn::ListOnly expects list name and tag who's value will be returned (if presenet). Received '%s'.", args)
  }
  listName := args[0]
  tag := args[1]

  if list, ok := p.Lists[listName]; ok {
    log.Debugf("Returning only entries from list '%s' with tag '%s'.", listName, tag)
    return list.only(tag), nil
  }
  return []string{}, fmt.Errorf("List '%s' not found by Fn::List.", listName)
}
