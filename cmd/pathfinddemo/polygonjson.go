// Copyright 2023 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"image"
)

// polygonsFromJSON loads a polygon set from JSON data as exported by the
// [Polygon Constructor].
//
// [Polygon Constructor]: https://alaricus.github.io/PolygonConstructor/
func polygonsFromJSON(jsonData []byte) (ps [][]image.Point, size image.Point, err error) {
	var jsonStruct polygonToolJSON
	err = json.Unmarshal(jsonData, &jsonStruct)
	if err != nil {
		return nil, image.Point{}, fmt.Errorf("could not unmarshal polygon JSON: %w", err)
	}
	return jsonStruct.polygons(), jsonStruct.size(), nil
}

type polygonToolJSON struct {
	Canvas struct {
		W int `json:"w"`
		H int `json:"h"`
	} `json:"canvas"`
	Polygons [][]struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"polygons"`
}

func (pj polygonToolJSON) size() image.Point {
	return image.Pt(pj.Canvas.W, pj.Canvas.H)
}

func (pj polygonToolJSON) polygons() [][]image.Point {
	ps := make([][]image.Point, len(pj.Polygons))
	for i, jsonPolygon := range pj.Polygons {
		p := make([]image.Point, len(jsonPolygon))
		for j, v := range jsonPolygon {
			p[j] = image.Pt(int(v.X), int(v.Y))
		}
		ps[i] = p
	}
	return ps
}
