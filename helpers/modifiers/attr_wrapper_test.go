package modifiers

import (
	"context"
	"errors"
	"testing"

	tt "github.com/vingarcia/ksql/helpers/testtools"
	"github.com/vingarcia/ksql/ksqlmodifiers"
)

func TestAttrScanWrapper(t *testing.T) {
	ctx := context.Background()

	var scanArgs map[string]interface{}
	wrapper := AttrScanWrapper{
		Ctx:     ctx,
		AttrPtr: "fakeAttrPtr",
		ScanFn: func(ctx context.Context, opInfo ksqlmodifiers.OpInfo, attrPtr interface{}, dbValue interface{}) error {
			scanArgs = map[string]interface{}{
				"opInfo":  opInfo,
				"attrPtr": attrPtr,
				"dbValue": dbValue,
			}
			return errors.New("fakeScanErrMsg")
		},
		OpInfo: ksqlmodifiers.OpInfo{
			Method:     "fakeMethod",
			DriverName: "fakeDriverName",
		},
	}

	err := wrapper.Scan("fakeDbValue")
	tt.AssertErrContains(t, err, "fakeScanErrMsg")
	tt.AssertEqual(t, scanArgs, map[string]interface{}{
		"opInfo": ksqlmodifiers.OpInfo{
			Method:     "fakeMethod",
			DriverName: "fakeDriverName",
		},
		"attrPtr": "fakeAttrPtr",
		"dbValue": "fakeDbValue",
	})
}

func TestAttrWrapper(t *testing.T) {
	ctx := context.Background()

	var valueArgs map[string]interface{}
	wrapper := AttrValueWrapper{
		Ctx:  ctx,
		Attr: "fakeAttr",
		ValueFn: func(ctx context.Context, opInfo ksqlmodifiers.OpInfo, inputValue interface{}) (outputValue interface{}, _ error) {
			valueArgs = map[string]interface{}{
				"opInfo":     opInfo,
				"inputValue": inputValue,
			}
			return "fakeOutputValue", errors.New("fakeValueErrMsg")
		},
		OpInfo: ksqlmodifiers.OpInfo{
			Method:     "fakeMethod",
			DriverName: "fakeDriverName",
		},
	}

	value, err := wrapper.Value()
	tt.AssertErrContains(t, err, "fakeValueErrMsg")
	tt.AssertEqual(t, valueArgs, map[string]interface{}{
		"opInfo": ksqlmodifiers.OpInfo{
			Method:     "fakeMethod",
			DriverName: "fakeDriverName",
		},
		"inputValue": "fakeAttr",
	})
	tt.AssertEqual(t, value, "fakeOutputValue")
}
