package models

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	type fields struct {
		Model    Model
		Name     string
		Email    string
		Password string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Valid User",
			fields: fields{Email: "gJohnDoe@gmail.com", Password: "Uu1?1234", Name: "stephane"},
			want:   true,
		},
		{
			name:   "No Capital letter",
			fields: fields{Email: "gJohnDoe@gmail.com", Password: "uu1?1234", Name: "stephane"},
			want:   false,
		},
		{
			name:   "No  lowercase letter",
			fields: fields{Email: "gJohnDoe@gmail.com", Password: "UU1?1234", Name: "stephane"},
			want:   false,
		},
		{
			name:   "Empty Mail",
			fields: fields{Email: "", Password: "Uu1?1234", Name: "stephane"},
			want:   false,
		},
		{
			name:   "Empty UserName",
			fields: fields{Email: "gJohnDoe@gmail.com", Password: "Uu1?1234", Name: " "},
			want:   false,
		},
		{
			name:   "short UserName",
			fields: fields{Email: "gJohnDoe@gmail.com", Password: "Uu1?1234", Name: "a"},
			want:   false,
		},
		{
			name:   "Empty password",
			fields: fields{Email: "", Password: "", Name: "stephane"},
			want:   false,
		},
		{
			name:   "Short password",
			fields: fields{Email: "gJohnDoe@gmail.com", Password: "Uu1?1", Name: ""},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				Model:    tt.fields.Model,
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			got, _ := u.Validate()
			if got != tt.want {
				t.Errorf("User.Validate() [%s] got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_isEmailValid(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Valid gmail address", args: args{e: "gJohnDoe@gmail.com"}, want: true},
		{name: "dot address", args: args{e: "gmt.stephane@gmail.com"}, want: true},
		{name: "no @", args: args{e: "gmt.stephanegmail.com"}, want: false},
		{name: "EmptyString", args: args{e: ""}, want: false},
		{name: "EmptyString", args: args{e: ""}, want: false},
		{name: "No Domain", args: args{e: "gmt.stephanegmail.com"}, want: false},
		{name: "With space", args: args{e: "gmtste phanegmail.com"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmailValid(tt.args.e); got != tt.want {
				t.Errorf("isEmailValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPasswordValid(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"NoCharacterAtAll",
			args{""},
			false,
		},
		{
			"JustEmptyStringAndWhitespace",
			args{" \n\t\r\v\f "},
			false,
		},
		{
			"MixtureOfEmptyStringAndWhitespace",
			args{"U u\n1\t?\r1\v2\f34"},
			false,
		},
		{
			"MissingUpperCaseString",
			args{"uu1?1234"},
			false,
		},
		{
			"MissingLowerCaseString",
			args{"UU1?1234"},
			false,
		},
		{
			"MissingNumber",
			args{"Uua?aaaa"},
			false,
		},
		{
			"MissingSymbol",
			args{"Uu101234"},
			false,
		},
		{
			"LessThanRequiredMinimumLength",
			args{"Uu1?123"},
			false,
		},
		{
			"ValidPassword",
			args{"Uu1?1234"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPasswordValid(tt.args.pass); got != tt.want {
				t.Errorf("isPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUserNameValid(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty userName",
			args: args{e: ""},
			want: false,
		},
		{
			name: "UserName special char",
			args: args{e: "?/&"},
			want: false,
		},
		{
			name: "UserName to long",
			args: args{e: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			want: false,
		},
		{
			name: "UserName with space",
			args: args{e: "John doe"},
			want: false,
		},
		{
			name: "Valid  username",
			args: args{e: "John"},
			want: true,
		},
		{
			name: "Valid username with accent",
			args: args{e: "testÀ"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUserNameValid(tt.args.e); got != tt.want {
				t.Errorf("isUserNameValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
