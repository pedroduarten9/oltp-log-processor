package logs

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/logs/v1"
)

type Attributable interface {
	GetAttributes() []*commonpb.KeyValue
}

type LogProcessor struct {
	key  string
	repo Repo
}

func NewLogProcessor(key string, repo Repo) *LogProcessor {
	return &LogProcessor{
		key:  key,
		repo: repo,
	}
}

func (lp *LogProcessor) ProcessLog(log *v1.ResourceLogs) {
	attrValue := lp.extractAttributeValue(log)
	if attrValue == "" {
		slog.Debug("unable to extract attribute", "log", log)
		return
	}

	lp.repo.IncrementAttribute(attrValue)
}

func (lp *LogProcessor) extractAttributeValue(log *v1.ResourceLogs) string {
	if log == nil {
		return ""
	}

	attr := lp.extractFromAttributable(log.Resource)
	if attr != "" {
		return attr
	}

	for _, scopeLogs := range log.ScopeLogs {
		attr = lp.extractFromAttributable(scopeLogs.Scope)
		if attr != "" {
			return attr
		}
		for _, record := range scopeLogs.LogRecords {
			attr = lp.extractFromAttributable(record)
			if attr != "" {
				return attr
			}
		}
	}

	return "unknown"
}

func (lp *LogProcessor) extractFromAttributable(attributable Attributable) string {
	if attributable == nil {
		return ""
	}

	for _, attr := range attributable.GetAttributes() {
		if attr.Key == lp.key {
			return stringifyAnyValue(attr.Value)
		}
	}

	return ""
}

func stringifyAnyValue(value *commonpb.AnyValue) string {
	switch v := value.Value.(type) {
	case *commonpb.AnyValue_StringValue:
		return v.StringValue
	case *commonpb.AnyValue_IntValue:
		return strconv.FormatInt(v.IntValue, 10)
	case *commonpb.AnyValue_DoubleValue:
		return strconv.FormatFloat(v.DoubleValue, 'f', -1, 64)
	case *commonpb.AnyValue_BoolValue:
		return strconv.FormatBool(v.BoolValue)
	case *commonpb.AnyValue_ArrayValue:
		return fmt.Sprintf("%v", v.ArrayValue.Values)
	case *commonpb.AnyValue_KvlistValue:
		return fmt.Sprintf("%v", v.KvlistValue.Values)
	default:
		slog.Debug("unsupported value type found", "type", v)
		return ""
	}
}

func (lp *LogProcessor) ReportAndReset() string {
	oldMap := lp.repo.Reset()
	if len(oldMap) == 0 {
		slog.Debug("No logs found")
		return ""
	}

	var buffer bytes.Buffer
	for value, count := range oldMap {
		_, err := buffer.WriteString(fmt.Sprintf("%s - %d\n", value, count))
		if err != nil {
			slog.Debug("Error printing log", "value", value)
			return ""
		}
	}

	return buffer.String()
}
