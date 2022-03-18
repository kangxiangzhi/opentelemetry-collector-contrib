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
	"context"
	"crypto/md5"
	"reflect"
	"strings"
	"time"

	traceTranslator "github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/tracetranslator"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.6.1"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/collector/model/pdata"
	commonV3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

var (
	blankSkyWalkingProtoSpan = new(v3.SpanObject)
	ComponentMap             map[int32]string
)

//SegmentToInternalTraces converts skywalking proto SegmentObject to internal traces
func SegmentToInternalTraces(ctx context.Context, segmentObject *v3.SegmentObject) pdata.Traces {
	traceData := pdata.NewTraces()

	if strings.Compare(segmentObject.GetService(), "") == 0 &&
		len(segmentObject.GetSpans()) == 0 {
		return traceData
	}

	rss := traceData.ResourceSpans()
	rss.EnsureCapacity(1)
	segmentToResourceSpans(ctx, segmentObject, rss.AppendEmpty())

	return traceData
}

//SegmentCollectionToInternalTraces converts multiple skywalking proto segmentCollection to internal traces
func SegmentCollectionToInternalTraces(ctx context.Context, segmentCollection *v3.SegmentCollection) pdata.Traces {
	traceData := pdata.NewTraces()
	if len(segmentCollection.Segments) == 0 {
		return traceData
	}

	segmentsLength := len(segmentCollection.Segments)

	rss := traceData.ResourceSpans()
	rss.EnsureCapacity(segmentsLength)

	index := 0
	for _, segmentObject := range segmentCollection.Segments {
		if strings.Compare(segmentObject.GetService(), "") == 0 &&
			len(segmentObject.GetSpans()) == 0 {
			continue
		}
		segmentToResourceSpans(ctx, segmentObject, rss.AppendEmpty())
		index++
	}
	return traceData
}

func segmentToResourceSpans(ctx context.Context, segmentObject *v3.SegmentObject, dest pdata.ResourceSpans) {
	segmentSpans := segmentObject.Spans
	segmentToInternalResource(ctx, segmentObject, dest.Resource())

	if len(segmentSpans) == 0 {
		return
	}

	groupByLibrary := spansToInternal(segmentObject)
	ilss := dest.InstrumentationLibrarySpans()

	for _, v := range groupByLibrary {
		ils := ilss.AppendEmpty()
		v.MoveTo(ils)
	}
}

func segmentToInternalResource(ctx context.Context, segmentObject *v3.SegmentObject, dest pdata.Resource) {
	if segmentObject.Service == "" {
		return
	}
	headers, ok := metadata.FromIncomingContext(ctx)
	if ok {
		attrs := dest.Attributes()
		attrs.EnsureCapacity(ResourceAttributesLength + len(headers))
		for key, header := range headers {
			if len(header) > 0 {
				attrs.InsertString(key, header[0])
			}
		}
		attrs.InsertString(AttributeBeginTime, time.Now().String())
		attrs.InsertString(conventions.AttributeServiceName, segmentObject.Service)
		portIpToken := strings.Split(segmentObject.ServiceInstance, "@")
		if len(portIpToken) == 2 {
			attrs.InsertString(AttributeLocalIp, portIpToken[1])
		} else {
			// 兜底策略
			attrs.InsertString(AttributeLocalIp, segmentObject.ServiceInstance)
		}
		// 添加语言tag
		addLanguageTag(attrs)
	}
}

func addLanguageTag(attrs pdata.AttributeMap) {
	userAgentValue, ok := attrs.Get(AttributeUserAgent)
	if ok {
		userAgentStringValue := userAgentValue.StringVal()
		split := strings.Split(userAgentStringValue, "/")
		languageValue := "unknown"
		if split[0] == ResourceUserAgentJava {
			languageValue = "java"
		} else if split[0] == ResourceUserAgentGo {
			languageValue = "go"
		}
		attrs.InsertString(AttributeLanguage, languageValue)
	}

}

//SpansToInternal
func spansToInternal(
	segmentObject *v3.SegmentObject) map[instrumentationLibrary]pdata.InstrumentationLibrarySpans {

	spansByLibrary := make(map[instrumentationLibrary]pdata.InstrumentationLibrarySpans)

	for _, span := range segmentObject.Spans {
		if span == nil || reflect.DeepEqual(span, blankSkyWalkingProtoSpan) {
			continue
		}
		pSpan, library := spanToInternal(segmentObject, span)
		ils, found := spansByLibrary[library]
		if !found {
			ils = pdata.NewInstrumentationLibrarySpans()
			spansByLibrary[library] = ils

			if library.name != "" {
				ils.InstrumentationLibrary().SetName(library.name)
				ils.InstrumentationLibrary().SetVersion(library.version)
			}
		}
		is := ils.Spans().AppendEmpty()
		pSpan.MoveTo(is)
		//ils.Spans().Append(pSpan)
	}
	return spansByLibrary
}

type instrumentationLibrary struct {
	name, version string
}

//SpanToInternal
func spanToInternal(segmentObject *v3.SegmentObject,
	span *v3.SpanObject) (pdata.Span, instrumentationLibrary) {
	dest := pdata.NewSpan()
	traceID := stringToTraceID(segmentObject.TraceId)
	dest.SetTraceID(traceID)
	dest.SetSpanID(stringToSpanID(segmentObject.TraceSegmentId, span.SpanId))
	dest.SetName(span.OperationName)
	dest.SetStartTimestamp(millisecondsToUnixNano(span.StartTime))
	if span.StartTime == span.EndTime {
		dest.SetEndTimestamp(millisecondsToUnixNano(span.EndTime + defaultOneMillisecond))
	} else {
		dest.SetEndTimestamp(millisecondsToUnixNano(span.EndTime))
	}
	//log.Println("service:", segmentObject.Service, "name:", span.OperationName,
	//	"startTime:", span.StartTime, "duration:", span.EndTime-span.StartTime)

	//spanType -> kind
	dest.SetKind(spanTypeToSpanKind(span))

	if span.ParentSpanId > -1 {
		parentSpanID := stringToSpanID(segmentObject.TraceSegmentId, span.ParentSpanId)
		dest.SetParentSpanID(parentSpanID)
		if len(span.Refs) == 0 {
			spanParentIdToSpanLinks(traceID, parentSpanID, dest.Links())
		}
	}

	ipPort := strings.Split(span.Peer, ":")
	attrs := dest.Attributes()
	peerLength := len(ipPort)
	if span.Peer != "" && peerLength == 2 {
		peerLength++
		attrs.EnsureCapacity(len(span.Tags) + peerLength)
		attrs.UpsertString(conventions.AttributeNetPeerIP, ipPort[0])
		attrs.UpsertString(conventions.AttributeNetPeerPort, ipPort[1])
	} else {
		attrs.EnsureCapacity(len(span.Tags) + 1)
	}
	//componentId to attributes
	component := ComponentMap[span.ComponentId]
	attrs.UpsertString(AttributeComponent, component)
	attrs.UpsertString(TagTraceID, segmentObject.TraceId)

	segmentReferencesToSpanLinks(span.Refs, span.ParentSpanId, dest.Links(), attrs)

	//tags to attributes
	tagsToInternalAttributes(span.Tags, attrs)
	componentTagsChange(span, attrs)

	il := instrumentationLibrary{}
	il.name, il.version = getLibraryNameAndVersion(attrs)
	//isError -> Tag
	spanErrorToSpanStatus(span.IsError, attrs, dest.Status())

	if attrs.Len() == 0 {
		attrs = pdata.NewAttributeMapFromMap(nil)
	}

	logsToSpanEvents(span.Logs, dest.Events())

	return dest, il
}

//tagsToInternalAttributes sets internal span Tags based on skyWalking span tags
func tagsToInternalAttributes(tags []*commonV3.KeyStringValuePair, dest pdata.AttributeMap) {
	for _, tag := range tags {
		dest.UpsertString(tag.Key, tag.Value)
	}
}

func componentTagsChange(span *v3.SpanObject, attrs pdata.AttributeMap) {
	if span.SpanLayer == v3.SpanLayer_MQ {
		attributesBroker, bFind := attrs.Get(AttributeMqBroker)
		if bFind {
			peerService := ComponentMap[span.ComponentId] + "-" + attributesBroker.StringVal()
			attrs.UpsertString(conventions.AttributePeerService, peerService)
		}
	}

	_, bFind := attrs.Get(AttributeDbType)
	if bFind {
		peerService := ComponentMap[span.ComponentId] + "-" + span.Peer
		attrs.UpsertString(conventions.AttributePeerService, peerService)
		peerIp, bOk := attrs.Get(conventions.AttributeNetPeerIP)
		if bOk {
			attrs.UpsertString(AttributeDbIp, peerIp.StringVal())
		}
	}
}

//logsToSpanEvents
func logsToSpanEvents(logs []*v3.Log, dest pdata.SpanEventSlice) {
	if len(logs) == 0 {
		return
	}
	length := len(logs)
	index := 0
	dest.EnsureCapacity(length)
	for reverseIndex := length - 1; reverseIndex >= 0; reverseIndex-- {
		log := logs[reverseIndex]
		event := dest.AppendEmpty()
		event.SetTimestamp(millisecondsToUnixNano(log.Time))
		if len(log.Data) == 0 {
			continue
		}

		attrs := event.Attributes()
		attrs.EnsureCapacity(len(log.Data))
		tagsToInternalAttributes(log.Data, attrs)
		if errorValue, okError := attrs.Get(TagEvent); okError {
			if errorValue.StringVal() == traceTranslator.TagError {
				if name, ok := attrs.Get(traceTranslator.TagMessage); ok {
					//event.SetName(name.StringVal())
					attrs.Upsert(AttributeErrorMessage, name)
					attrs.Delete(traceTranslator.TagMessage)
				}
			}
		}
		index++
	}
}

//segmentReferencesToSpanLinks
func segmentReferencesToSpanLinks(refs []*v3.SegmentReference, excludeParentID int32,
	dest pdata.SpanLinkSlice, destAttribute pdata.AttributeMap) {
	if len(refs) == 0 || len(refs) == 1 && refs[0].ParentSpanId == excludeParentID {
		return
	}
	dest.EnsureCapacity(len(refs))
	index := 0
	for _, ref := range refs {
		link := dest.AppendEmpty()
		if ref.ParentSpanId == excludeParentID {
			continue
		}
		link.SetTraceID(stringToTraceID(ref.TraceId))
		link.SetSpanID(stringToSpanID(ref.ParentTraceSegmentId, ref.ParentSpanId))
		if ref.ParentService != "" {
			destAttribute.InsertString(conventions.AttributePeerService, ref.ParentService)
		}
		if ref.ParentServiceInstance != "" {
			instanceIp := strings.Split(ref.ParentServiceInstance, "@")
			if len(instanceIp) == 2 {
				destAttribute.InsertString(conventions.AttributeNetPeerIP, instanceIp[1])
			}
		}
		index++
	}

	/*if index < len(refs) {
		dest.Resize(index)
	}*/
}

func spanParentIdToSpanLinks(traceID pdata.TraceID, parentSpanID pdata.SpanID, dest pdata.SpanLinkSlice) {
	dest.EnsureCapacity(1)
	link := dest.AppendEmpty()
	link.SetTraceID(traceID)
	link.SetSpanID(parentSpanID)
}

//spanErrorToSpanStatus
func spanErrorToSpanStatus(isError bool, attrs pdata.AttributeMap, dest pdata.SpanStatus) {
	statusCode := pdata.StatusCodeOk
	statusMessage := ""
	statusExists := false
	if isError {
		statusExists = true
		statusCode = pdata.StatusCodeError
	}

	if _, ok := attrs.Get(TagStatusCode); ok {
		statusExists = true
		if msgAttr, ok := attrs.Get(TagStatusMsg); ok {
			statusMessage = msgAttr.StringVal()
			attrs.Delete(TagStatusMsg)
		}
	} else if _, ok := attrs.Get(conventions.AttributeHTTPStatusCode); ok {
		statusExists = true
		if msgAttr, ok := attrs.Get(traceTranslator.TagHTTPStatusMsg); ok && isError {
			statusMessage = msgAttr.StringVal()
		}
	}
	dest.SetCode(statusCode)
	if statusExists {
		dest.SetMessage(statusMessage)
	}
}

//spanTypeToSpanKind
func spanTypeToSpanKind(span *v3.SpanObject) pdata.SpanKind {
	if span.SpanType == v3.SpanType_Exit {
		if span.SpanLayer == v3.SpanLayer_MQ {
			return pdata.SpanKindProducer
		} else {
			return pdata.SpanKindClient
		}
	}

	if span.SpanType == v3.SpanType_Entry {
		if span.SpanLayer == v3.SpanLayer_MQ {
			return pdata.SpanKindConsumer
		} else {
			return pdata.SpanKindServer
		}
	}

	if span.SpanType == v3.SpanType_Local {
		return pdata.SpanKindInternal
	}

	return pdata.SpanKindUnspecified
}

//stringToTraceID
func stringToTraceID(traceId string) pdata.TraceID {
	md5Hash := md5.New()
	md5Hash.Write([]byte(traceId))
	hashTraceId := md5Hash.Sum([]byte(""))
	var byteTraceId [16]byte
	copy(byteTraceId[:], hashTraceId)
	return pdata.NewTraceID(byteTraceId)
}

//stringToSpanID
func stringToSpanID(segmentId string, spanId int32) pdata.SpanID {
	md5Hash := md5.New()
	md5Hash.Write([]byte(segmentId + "." + string(spanId)))
	hashSpanId := md5Hash.Sum([]byte(""))
	var byteSpanId [8]byte
	copy(byteSpanId[:], hashSpanId)
	return pdata.NewSpanID(byteSpanId)
	//ids := strings.Split(segmentId, ".")
	//spanID := pdata.InvalidSpanID()
	//if len(ids) == 3 {
	//	float64Id, err1 := strconv.ParseFloat(ids[2], 64)
	//	threadId, err2 := strconv.ParseUint(ids[1], 16, 64)
	//	if err1 == nil && err2 == nil {
	//		//pdata.SpanID只有八个字节大小，容量限制了解决方案
	//		unit64Id := uint64(float64Id) + threadId + uint64(spanId)
	//		spanID = traceTranslator.UInt64ToSpanID(unit64Id)
	//	}
	//}
	//return spanID
}

func getLibraryNameAndVersion(attrs pdata.AttributeMap) (string, string) {
	var name, version string
	if libraryName, ok := attrs.Get(conventions.OtelLibraryName); ok {
		name = libraryName.StringVal()
		attrs.Delete(conventions.OtelLibraryName)
		if libraryVersion, ok := attrs.Get(conventions.OtelLibraryVersion); ok {
			version = libraryVersion.StringVal()
			attrs.Delete(conventions.OtelLibraryVersion)
		}
	}
	return name, version
}

// milliSecondsToUnixNano converts epoch microseconds to pdata.Timestamp
func millisecondsToUnixNano(ms int64) pdata.Timestamp {
	return pdata.Timestamp(uint64(ms) * 1e6)
}
