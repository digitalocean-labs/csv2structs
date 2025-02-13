package csv2structs

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSetup(t *testing.T, header []string, rows ...[]string) *os.File {
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	assert.NoError(t, err)

	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	assert.NoError(t, writer.Write(header))
	assert.NoError(t, writer.WriteAll(rows))

	pos, err := tmpFile.Seek(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), pos)

	t.Cleanup(func() {
		// we might already be closed
		_ = tmpFile.Close()
		assert.NoError(t, os.Remove(tmpFile.Name()))
	})

	return tmpFile
}

func TestParse_Simple(t *testing.T) {
	type simpleObj struct {
		Ident      int
		Name       string
		OtherThing string

		hiddenFieldsAreIgnored uint64
	}

	tests := []struct {
		name        string
		header      []string
		rows        [][]string
		options     []Option
		expected    []*simpleObj
		expectedErr error
	}{
		{
			"happy path",
			[]string{"ident", "name", "other_thing"},
			[][]string{
				{"1", "foo", "bar"},
				{"2", "bar", "baz"},
			},
			nil,
			[]*simpleObj{
				{
					1,
					"foo",
					"bar",
					0,
				},
				{
					2,
					"bar",
					"baz",
					0,
				},
			},
			nil,
		},
		{
			"snake case header transformation is default",
			[]string{"Ident", "Name", "OtherThing"},
			[][]string{
				{"1", "foo", "bar"},
				{"2", "bar", "baz"},
			},
			nil,
			nil,
			fmt.Errorf("missing headers: OtherThing, Otherthing"),
		},
		{
			"happy path without transforming headers",
			[]string{"Ident", "Name", "OtherThing"},
			[][]string{
				{"1", "foo", "bar"},
				{"2", "bar", "baz"},
			},
			[]Option{WithHeaderType(HeaderTypeNone)},
			[]*simpleObj{
				{
					1,
					"foo",
					"bar",
					0,
				},
				{
					2,
					"bar",
					"baz",
					0,
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := testSetup(t, tt.header, tt.rows...)

			actual, err := Parse[simpleObj](tmpFile, tt.options...)

			assert.EqualValues(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)

			for _, obj := range actual {
				assert.Equal(t, uint64(0), obj.hiddenFieldsAreIgnored)
			}
		})
	}
}

func TestParse_Complete(t *testing.T) {
	type completeObj struct {
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		String  string
		Bool    bool
	}

	tests := []struct {
		name        string
		header      []string
		rows        [][]string
		expected    []*completeObj
		expectedErr error
	}{
		{
			"happy path",
			[]string{"Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Float32", "Float64", "String", "Bool"},
			[][]string{
				{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11.1", "12.2", "foo", "true"},
			},
			[]*completeObj{
				{
					1,
					2,
					3,
					4,
					5,
					6,
					7,
					8,
					9,
					10,
					11.1,
					12.2,
					"foo",
					true,
				},
			},
			nil,
		},
		{
			"missing header is an error",
			[]string{"Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Float32", "Float64", "String"},
			[][]string{
				{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11.1", "12.2", "foo", "false"},
			},
			nil,
			fmt.Errorf("missing header: Bool"),
		},
		{
			"lower case headers are title cased, underscores are removed",
			[]string{"int", "int_8", "int_16", "int_32", "int_64", "uint", "uint_8", "uint_16", "uint_32", "uint_64", "float_32", "float_64", "string", "bool"},
			[][]string{
				{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11.1", "12.2", "bar", "false"},
			},
			[]*completeObj{
				{
					1,
					2,
					3,
					4,
					5,
					6,
					7,
					8,
					9,
					10,
					11.1,
					12.2,
					"bar",
					false,
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := testSetup(t, tt.header, tt.rows...)

			actual, err := Parse[completeObj](tmpFile)

			assert.EqualValues(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParse_ErrorInterface(t *testing.T) {
	type Foo interface{}
	resp, err := Parse[Foo](nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errInvalidType)
}

func TestParse_ErrorNoFields(t *testing.T) {
	type NoFields struct{}
	resp, err := Parse[NoFields](nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errNoVisibleFields)
}

func TestParse_ErrorNoExportedFields(t *testing.T) {
	type NoExportedFields struct {
		_ string
	}
	resp, err := Parse[NoExportedFields](nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, errNoExportedFields)
}

func TestParse_ErrorInvalidType(t *testing.T) {
	type InvalidType struct {
		// slices as a field type don't really make sense for CSV...
		// you probably want a string then json load that later or something
		FieldA []string
	}

	csvFile := testSetup(t, []string{"field_a"}, []string{`["foo","bar","baz"]`})
	resp, err := Parse[InvalidType](csvFile)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "unsupported type: slice", err.Error())
}

func TestParse_ErrorRead(t *testing.T) {
	type Foo struct {
		Bar string
	}

	tmpFile := testSetup(t, []string{"bar"}, []string{"baz"})

	// close the file so it can't be read from
	assert.NoError(t, tmpFile.Close())

	resp, err := Parse[Foo](tmpFile)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file already closed")
}

func TestParse_ErrorIntConversion(t *testing.T) {
	type Foo struct {
		Bar int
	}

	tmpFile := testSetup(t, []string{"bar"}, []string{"baz"})

	resp, err := Parse[Foo](tmpFile)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), `failed to convert "baz" to int`)
}

func TestParse_ErrorBoolConversion(t *testing.T) {
	type Foo struct {
		Bar bool
	}

	tmpFile := testSetup(t, []string{"bar"}, []string{"baz"})

	resp, err := Parse[Foo](tmpFile)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), `failed to convert "baz" to bool`)
}

func TestParse_ErrorFloatConversion(t *testing.T) {
	type Foo struct {
		Bar float64
	}

	tmpFile := testSetup(t, []string{"bar"}, []string{"baz"})

	resp, err := Parse[Foo](tmpFile)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), `failed to convert "baz" to float`)
}

func TestParse_ErrorUintConversion(t *testing.T) {
	type Foo struct {
		Bar uint
	}

	tmpFile := testSetup(t, []string{"bar"}, []string{"baz"})

	resp, err := Parse[Foo](tmpFile)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), `failed to convert "baz" to uint`)
}

func TestParse_CustomTransform(t *testing.T) {
	type Foo struct {
		BAR string
	}

	tmpFile := testSetup(t, []string{"bar"}, []string{"baz"})

	screamer := func(s string) string {
		return strings.ToUpper(s)
	}

	resp, err := Parse[Foo](tmpFile, WithHeaderTransform(screamer))

	assert.NoError(t, err)
	assert.Equal(t, []*Foo{{"baz"}}, resp)
}

func TestParse_LeadingByteOrderMark(t *testing.T) {
	type Foo struct {
		Foo string
		Bar string
	}

	reader := strings.NewReader(
		"\xEF\xBB\xBF" + // BOM
			"foo,bar\n" + // header
			"baz,foo\n", // row
	)

	resp, err := Parse[Foo](reader)

	assert.NoError(t, err)
	assert.Equal(t, []*Foo{{"baz", "foo"}}, resp)
}
