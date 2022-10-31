// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/v2"
	"github.com/pb33f/libopenapi/datamodel/low/v3"
	"reflect"
)

type ResponsesChanges struct {
	PropertyChanges
	ResponseChanges  map[string]*ResponseChanges
	DefaultChanges   *ResponseChanges
	ExtensionChanges *ExtensionChanges
}

func (r *ResponsesChanges) TotalChanges() int {
	c := r.PropertyChanges.TotalChanges()
	for k := range r.ResponseChanges {
		c += r.ResponseChanges[k].TotalChanges()
	}
	if r.DefaultChanges != nil {
		c += r.DefaultChanges.TotalChanges()
	}
	if r.ExtensionChanges != nil {
		c += r.ExtensionChanges.TotalChanges()
	}
	return c
}

func (r *ResponsesChanges) TotalBreakingChanges() int {
	c := r.PropertyChanges.TotalBreakingChanges()
	for k := range r.ResponseChanges {
		c += r.ResponseChanges[k].TotalBreakingChanges()
	}
	if r.DefaultChanges != nil {
		c += r.DefaultChanges.TotalBreakingChanges()
	}
	return c
}

func CompareResponses(l, r any) *ResponsesChanges {

	var changes []*Change

	rc := new(ResponsesChanges)

	// swagger
	if reflect.TypeOf(&v2.Responses{}) == reflect.TypeOf(l) &&
		reflect.TypeOf(&v2.Responses{}) == reflect.TypeOf(r) {

		lResponses := l.(*v2.Responses)
		rResponses := r.(*v2.Responses)

		// perform hash check to avoid further processing
		if low.AreEqual(lResponses, rResponses) {
			return nil
		}

		if !lResponses.Default.IsEmpty() && !rResponses.Default.IsEmpty() {
			rc.DefaultChanges = CompareResponse(lResponses.Default.Value, rResponses.Default.Value)
		}
		if !lResponses.Default.IsEmpty() && rResponses.Default.IsEmpty() {
			CreateChange(&changes, ObjectRemoved, v3.DefaultLabel,
				lResponses.Default.ValueNode, nil, true,
				lResponses.Default.Value, nil)
		}
		if lResponses.Default.IsEmpty() && !rResponses.Default.IsEmpty() {
			CreateChange(&changes, ObjectAdded, v3.DefaultLabel,
				nil, rResponses.Default.ValueNode, false,
				nil, lResponses.Default.Value)
		}

		rc.ResponseChanges = CheckMapForChanges(lResponses.Codes, rResponses.Codes,
			&changes, v3.CodesLabel, CompareResponseV2)

		rc.ExtensionChanges = CompareExtensions(lResponses.Extensions, rResponses.Extensions)
	}

	// openapi
	if reflect.TypeOf(&v3.Responses{}) == reflect.TypeOf(l) &&
		reflect.TypeOf(&v3.Responses{}) == reflect.TypeOf(r) {

		lResponses := l.(*v3.Responses)
		rResponses := r.(*v3.Responses)

		// perform hash check to avoid further processing
		if low.AreEqual(lResponses, rResponses) {
			return nil
		}

		if !lResponses.Default.IsEmpty() && !rResponses.Default.IsEmpty() {
			rc.DefaultChanges = CompareResponse(lResponses.Default.Value, rResponses.Default.Value)
		}
		if !lResponses.Default.IsEmpty() && rResponses.Default.IsEmpty() {
			CreateChange(&changes, ObjectRemoved, v3.DefaultLabel,
				lResponses.Default.ValueNode, nil, true,
				lResponses.Default.Value, nil)
		}
		if lResponses.Default.IsEmpty() && !rResponses.Default.IsEmpty() {
			CreateChange(&changes, ObjectAdded, v3.DefaultLabel,
				nil, rResponses.Default.ValueNode, false,
				nil, lResponses.Default.Value)
		}

		rc.ResponseChanges = CheckMapForChanges(lResponses.Codes, rResponses.Codes,
			&changes, v3.CodesLabel, CompareResponseV3)

		rc.ExtensionChanges = CompareExtensions(lResponses.Extensions, rResponses.Extensions)

	}

	rc.Changes = changes
	return rc
}