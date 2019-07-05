package resource

import "testing"

func TestPropertyRendering(t *testing.T) {
	tests := []struct {
		prop Property
		want string
	}{
		{NewStringProperty(""), "\"\""},
		{NewStringProperty("foo"), "\"foo\""},
		{NewIntegerProperty(0), "0"},
		{NewIntegerProperty(999), "999"},
		{NewIntegerProperty(-1000), "-1000"},
		{NewRawProperty(""), ""},
		{NewRawProperty("foo"), "foo"},
		{NewRawProperty("resource.name.id"), "resource.name.id"},
		{NewArrayProperty(), "[  ]"},
		{NewArrayProperty(
			NewIntegerProperty(1),
			NewIntegerProperty(2),
			NewIntegerProperty(3),
		), "[ 1, 2, 3 ]"},
		{NewArrayProperty(
			NewStringProperty("1"),
			NewStringProperty("2"),
			NewStringProperty("3"),
		), "[ \"1\", \"2\", \"3\" ]"},
		{NewArrayProperty(
			NewRawProperty("res.name-1.id"),
			NewRawProperty("res.name-2.id"),
			NewRawProperty("res.name-3.id"),
		), "[ res.name-1.id, res.name-2.id, res.name-3.id ]"},
		{NewArrayProperty(
			NewIntegerProperty(1),
			NewStringProperty("2"),
			NewRawProperty("res.name.id"),
			NewArrayProperty(
				NewIntegerProperty(0),
				NewIntegerProperty(1),
			),
		), "[ 1, \"2\", res.name.id, [ 0, 1 ] ]"},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			got := test.prop.Render()
			if got != test.want {
				t.Errorf("render prop: got '%s', want '%s'", got, test.want)
			}
		})
	}
}
