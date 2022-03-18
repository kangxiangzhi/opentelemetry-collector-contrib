// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package skywalkingv2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/model/pdata"
)

func TestStringToTraceID(t *testing.T) {
	skywalkingTraceId := "21de58a4f4b7426b886e7bf2f285897e.41.16197491006830000"
	hexLength := len(skywalkingTraceId)
	assert.Greater(t, hexLength, 32)
	traceID := stringToTraceID(skywalkingTraceId)
	assert.NotEqual(t, "21de58a4f4b7426b886e7bf2f285897e", traceID.HexString())
	assert.NotEqual(t, pdata.InvalidTraceID(), traceID)
}

func TestStringToSpanID(t *testing.T) {
	skywalkingSegmentId := "21de58a4f4b7426b886e7bf2f285897e.41.16197491006830000"
	var spanId int32
	spanId = 3
	spanID3 := stringToSpanID(skywalkingSegmentId, spanId)
	assert.NotEqual(t, pdata.InvalidSpanID(), spanID3)

	spanId = 6
	spanID6 := stringToSpanID(skywalkingSegmentId, spanId)
	assert.NotEqual(t, spanID3, spanID6)

	spanId = 9
	spanID9 := stringToSpanID(skywalkingSegmentId, spanId)
	assert.NotEqual(t, spanID9, spanID6)
}
