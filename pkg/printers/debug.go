// Copyright 2023 Authzed, Inc.
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//        http://www.apache.org/licenses/LICENSE-2.0
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package printers

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/spicedb/pkg/tuple"
	"github.com/gookit/color"
)

// DisplayCheckTrace prints out the check trace found in the given debug message.
func DisplayCheckTrace(checkTrace *v1.CheckDebugTrace, tp *TreePrinter, hasError bool) {
	displayCheckTrace(checkTrace, tp, hasError, map[string]struct{}{})
}

func displayCheckTrace(checkTrace *v1.CheckDebugTrace, tp *TreePrinter, hasError bool, encountered map[string]struct{}) {
	red := color.FgRed.Render
	green := color.FgGreen.Render
	cyan := color.FgCyan.Render
	white := color.FgWhite.Render
	faint := color.FgGray.Render
	magenta := color.FgMagenta.Render

	orange := color.C256(166).Sprint
	purple := color.C256(99).Sprint
	lightgreen := color.C256(35).Sprint
	caveatColor := color.C256(198).Sprint

	hasPermission := green("✓")
	resourceColor := white
	permissionColor := color.FgWhite.Render

	if checkTrace.PermissionType == v1.CheckDebugTrace_PERMISSION_TYPE_PERMISSION {
		permissionColor = lightgreen
	} else if checkTrace.PermissionType == v1.CheckDebugTrace_PERMISSION_TYPE_RELATION {
		permissionColor = orange
	}

	if checkTrace.Result == v1.CheckDebugTrace_PERMISSIONSHIP_CONDITIONAL_PERMISSION {
		switch checkTrace.CaveatEvaluationInfo.Result {
		case v1.CaveatEvalInfo_RESULT_FALSE:
			hasPermission = red("⨉")
			resourceColor = faint
			permissionColor = faint

		case v1.CaveatEvalInfo_RESULT_MISSING_SOME_CONTEXT:
			hasPermission = magenta("?")
			resourceColor = faint
			permissionColor = faint
		}
	} else if checkTrace.Result != v1.CheckDebugTrace_PERMISSIONSHIP_HAS_PERMISSION {
		hasPermission = red("⨉")
		resourceColor = faint
		permissionColor = faint
	}

	additional := ""
	if checkTrace.GetWasCachedResult() {
		additional = cyan(" (cached)")
	} else if hasError && isPartOfCycle(checkTrace, map[string]struct{}{}) {
		hasPermission = orange("!")
		resourceColor = white
	}

	isEndOfCycle := false
	if hasError {
		key := cycleKey(checkTrace)
		_, isEndOfCycle = encountered[key]
		if isEndOfCycle {
			additional = color.C256(166).Sprint(" (cycle)")
		}
		encountered[key] = struct{}{}
	}

	timing := ""
	if checkTrace.Duration != nil {
		timing = fmt.Sprintf(" (%s)", checkTrace.Duration.AsDuration().String())
	}

	tp = tp.Child(
		fmt.Sprintf(
			"%s %s:%s %s%s%s",
			hasPermission,
			resourceColor(checkTrace.Resource.ObjectType),
			resourceColor(checkTrace.Resource.ObjectId),
			permissionColor(checkTrace.Permission),
			additional,
			timing,
		),
	)

	if isEndOfCycle {
		return
	}

	if checkTrace.GetCaveatEvaluationInfo() != nil {
		indicator := ""
		exprColor := color.FgWhite.Render
		switch checkTrace.CaveatEvaluationInfo.Result {
		case v1.CaveatEvalInfo_RESULT_FALSE:
			indicator = red("⨉")
			exprColor = faint

		case v1.CaveatEvalInfo_RESULT_TRUE:
			indicator = green("✓")

		case v1.CaveatEvalInfo_RESULT_MISSING_SOME_CONTEXT:
			indicator = magenta("?")
		}

		white := color.HEXStyle("fff")
		white.SetOpts(color.Opts{color.OpItalic})

		contextMap := checkTrace.CaveatEvaluationInfo.Context.AsMap()
		caveatName := checkTrace.CaveatEvaluationInfo.CaveatName

		c := tp.Child(fmt.Sprintf("%s %s %s", indicator, exprColor(checkTrace.CaveatEvaluationInfo.Expression), caveatColor(caveatName)))
		if len(contextMap) > 0 {
			contextJSON, _ := json.MarshalIndent(contextMap, "", "  ")
			c.Child(string(contextJSON))
		} else {
			if checkTrace.CaveatEvaluationInfo.Result != v1.CaveatEvalInfo_RESULT_MISSING_SOME_CONTEXT {
				c.Child(faint("(no matching context found)"))
			}
		}

		if checkTrace.CaveatEvaluationInfo.Result == v1.CaveatEvalInfo_RESULT_MISSING_SOME_CONTEXT {
			c.Child(fmt.Sprintf("missing context: %s", strings.Join(checkTrace.CaveatEvaluationInfo.PartialCaveatInfo.MissingRequiredContext, ", ")))
		}
	}

	if checkTrace.GetSubProblems() != nil {
		for _, subProblem := range checkTrace.GetSubProblems().Traces {
			displayCheckTrace(subProblem, tp, hasError, encountered)
		}
	} else if checkTrace.Result == v1.CheckDebugTrace_PERMISSIONSHIP_HAS_PERMISSION {
		tp.Child(purple(fmt.Sprintf("%s:%s %s", checkTrace.Subject.Object.ObjectType, checkTrace.Subject.Object.ObjectId, checkTrace.Subject.OptionalRelation)))
	}
}

func cycleKey(checkTrace *v1.CheckDebugTrace) string {
	return fmt.Sprintf("%s#%s", tuple.StringObjectRef(checkTrace.Resource), checkTrace.Permission)
}

func isPartOfCycle(checkTrace *v1.CheckDebugTrace, encountered map[string]struct{}) bool {
	if checkTrace.GetSubProblems() == nil {
		return false
	}

	encounteredCopy := make(map[string]struct{}, len(encountered))
	for k, v := range encountered {
		encounteredCopy[k] = v
	}

	key := cycleKey(checkTrace)
	if _, ok := encounteredCopy[key]; ok {
		return true
	}

	encounteredCopy[key] = struct{}{}

	for _, subProblem := range checkTrace.GetSubProblems().Traces {
		if isPartOfCycle(subProblem, encounteredCopy) {
			return true
		}
	}

	return false
}
