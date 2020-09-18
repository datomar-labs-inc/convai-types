package ctypes

import (
	"net/url"
	"reflect"
	"testing"
)

func TestResourceQueryFromURL(t *testing.T) {
	url1, _ := url.Parse("https://test.url?limit=6&mode=nor&offset=86665&eq.test=floobs&eq.test=floobs2&!ex.testfield&!sw.bubs=floobitybits&CONT.stuff=yeah&ew.YeAh=55&rx.ffld=%5Eflub%5C((.*)%20YEE%5C)%24&sort=fields&sortdesc=fields2")

	type args struct {
		url         *url.URL
		fieldPrefix string
	}
	tests := []struct {
		name string
		args args
		want *ResourceQuery
	}{
		{
			name: "t1",
			args: args{
				url:         url1,
				fieldPrefix: "fp.",
			},
			want: (&ResourceQuery{
				Limit:  6,
				Offset: 86665,
				Mode:   RQNor,
				Sort: []RQSort{
					{
						Field:     "fp.fields",
						Ascending: true,
					},
					{
						Field:     "fp.fields2",
						Ascending: false,
					},
				},
				Queries: []RQQ{
					{
						Field:    "fp.test",
						Operator: RQEquals,
						Value:    StrPtr("floobs"),
					},
					{
						Field:    "fp.test",
						Operator: RQEquals,
						Value:    StrPtr("floobs2"),
					},
					{
						Field:    "fp.testfield",
						Operator: RQExists,
						Negate:   true,
					},
					{
						Field:    "fp.bubs",
						Operator: RQStartsWith,
						Value:    StrPtr("floobitybits"),
						Negate:   true,
					},
					{
						Field:    "fp.stuff",
						Operator: RQContains,
						Value:    StrPtr("yeah"),
					},
					{
						Field:    "fp.YeAh",
						Operator: RQEndsWith,
						Value:    StrPtr("55"),
					},
					{
						Field:    "fp.ffld",
						Operator: RQRegex,
						Value:    StrPtr("^flub\\((.*) YEE\\)$"),
					},
				},
			}).sortQueriesByField(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ResourceQueryFromURL(tt.args.url, tt.args.fieldPrefix); err != nil || !reflect.DeepEqual(got, tt.want) {
				if err != nil {
					t.Error(err)
				} else if len(got.Queries) != len(tt.want.Queries) {
					t.Errorf("queries arrays differ in length: got %d, want %d", len(got.Queries), len(tt.want.Queries))
				} else {
					for i, s := range got.Queries {
						if !reflect.DeepEqual(s, tt.want.Queries[i]) {
							t.Errorf("Unmatching queries %d: %v, want %v", i, s, tt.want.Queries[i])
						}
					}
				}

				t.Errorf("ResourceQueryFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRQQ_ValueAsTyped(t *testing.T) {
	type fields struct {
		Field    string
		Operator int
		Value    *string
		Negate   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "Normal Integer",
			fields: fields{
				Value: StrPtr("5"),
			},
			want: int64(5),
		},
		{
			name: "Negative Integer",
			fields: fields{
				Value: StrPtr("-5068"),
			},
			want: int64(-5068),
		},
		{
			name: "Normal Float",
			fields: fields{
				Value: StrPtr("5.897"),
			},
			want: float64(5.897),
		},
		{
			name: "Negative Float",
			fields: fields{
				Value: StrPtr("-5068.5566"),
			},
			want: float64(-5068.5566),
		},
		{
			name: "Overflow Integer",
			fields: fields{
				Value: StrPtr("10223372036854775807"),
			},
			want: "10223372036854775807",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RQQ{
				Field:    tt.fields.Field,
				Operator: tt.fields.Operator,
				Value:    tt.fields.Value,
				Negate:   tt.fields.Negate,
			}
			if got := r.ValueAsTyped(); got != tt.want {
				t.Errorf("ValueAsTyped() = %v (%s), want %v (%s)", got, reflect.ValueOf(got).Kind().String(), tt.want, reflect.ValueOf(tt.want).Kind().String())
			}
		})
	}
}
