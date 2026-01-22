package check

import (
	"testing"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"test@example.com", false},
		{"user.name@domain.co.uk", false},
		{"user+tag@example.com", false},
		{"invalid", true},
		{"@example.com", true},
		{"test@", true},
	}
	for _, tt := range tests {
		v := Email(tt.input, "email")
		if v.Failed() != tt.wantErr {
			t.Errorf("Email(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://example.com/path?q=1", false},
		{"ftp://files.example.com", false},
		{"not-a-url", true},
		{"://missing-scheme.com", true},
	}
	for _, tt := range tests {
		v := URL(tt.input, "url")
		if v.Failed() != tt.wantErr {
			t.Errorf("URL(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestHTTPOrHTTPS(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://example.com", false},
		{"ftp://example.com", true},
		{"not-a-url", true},
	}
	for _, tt := range tests {
		v := HTTPOrHTTPS(tt.input, "url")
		if v.Failed() != tt.wantErr {
			t.Errorf("HTTPOrHTTPS(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", false},
		{"550e8400-e29b-11d4-a716-446655440000", false}, // v1
		{"550E8400-E29B-41D4-A716-446655440000", false}, // uppercase
		{"not-a-uuid", true},
		{"550e8400e29b41d4a716446655440000", true}, // no hyphens
	}
	for _, tt := range tests {
		v := UUID(tt.input, "uuid")
		if v.Failed() != tt.wantErr {
			t.Errorf("UUID(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestUUID4(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", false},
		{"550e8400-e29b-11d4-a716-446655440000", true}, // v1 not allowed
	}
	for _, tt := range tests {
		v := UUID4(tt.input, "uuid")
		if v.Failed() != tt.wantErr {
			t.Errorf("UUID4(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestIP(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"192.168.1.1", false},
		{"::1", false},
		{"2001:db8::1", false},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := IP(tt.input, "ip")
		if v.Failed() != tt.wantErr {
			t.Errorf("IP(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestIPv4(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"192.168.1.1", false},
		{"::1", true},
	}
	for _, tt := range tests {
		v := IPv4(tt.input, "ip")
		if v.Failed() != tt.wantErr {
			t.Errorf("IPv4(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestIPv6(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"::1", false},
		{"2001:db8::1", false},
		{"192.168.1.1", true},
	}
	for _, tt := range tests {
		v := IPv6(tt.input, "ip")
		if v.Failed() != tt.wantErr {
			t.Errorf("IPv6(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestCIDR(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"192.168.1.0/24", false},
		{"10.0.0.0/8", false},
		{"192.168.1.1", true},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := CIDR(tt.input, "cidr")
		if v.Failed() != tt.wantErr {
			t.Errorf("CIDR(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestMAC(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"01:23:45:67:89:AB", false},
		{"01-23-45-67-89-AB", false},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := MAC(tt.input, "mac")
		if v.Failed() != tt.wantErr {
			t.Errorf("MAC(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestHostname(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"example.com", false},
		{"sub.example.com", false},
		{"localhost", false},
		{"-invalid.com", true},
	}
	for _, tt := range tests {
		v := Hostname(tt.input, "hostname")
		if v.Failed() != tt.wantErr {
			t.Errorf("Hostname(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestPort(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"80", false},
		{"443", false},
		{"65535", false},
		{"0", true},
		{"65536", true},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := Port(tt.input, "port")
		if v.Failed() != tt.wantErr {
			t.Errorf("Port(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestHostPort(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"example.com:80", false},
		{"localhost:443", false},
		{"example.com", true},
		{":80", true},
	}
	for _, tt := range tests {
		v := HostPort(tt.input, "hostport")
		if v.Failed() != tt.wantErr {
			t.Errorf("HostPort(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestHexColor(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"#fff", false},
		{"#ffffff", false},
		{"#ffffffff", false},
		{"fff", true},
		{"#gg0000", true},
	}
	for _, tt := range tests {
		v := HexColor(tt.input, "color")
		if v.Failed() != tt.wantErr {
			t.Errorf("HexColor(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestBase64(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"SGVsbG8gV29ybGQ=", false},
		{"!!!invalid", true},
	}
	for _, tt := range tests {
		v := Base64(tt.input, "data")
		if v.Failed() != tt.wantErr {
			t.Errorf("Base64(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{`{"key": "value"}`, false},
		{`[1, 2, 3]`, false},
		{`"string"`, false},
		{`invalid`, true},
		{`{key: "value"}`, true},
	}
	for _, tt := range tests {
		v := JSON(tt.input, "json")
		if v.Failed() != tt.wantErr {
			t.Errorf("JSON(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestSemver(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"1.0.0", false},
		{"v1.0.0", false},
		{"1.0.0-alpha", false},
		{"1.0.0-alpha.1", false},
		{"1.0.0+build", false},
		{"1.0", true},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := Semver(tt.input, "version")
		if v.Failed() != tt.wantErr {
			t.Errorf("Semver(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestE164(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"+14155552671", false},
		{"+442071234567", false},
		{"14155552671", true},
		{"+0123456789", true},
	}
	for _, tt := range tests {
		v := E164(tt.input, "phone")
		if v.Failed() != tt.wantErr {
			t.Errorf("E164(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestCreditCard(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"4111111111111111", false},    // valid Visa
		{"4111 1111 1111 1111", false}, // with spaces
		{"4111-1111-1111-1111", false}, // with dashes
		{"4111111111111112", true},     // invalid checksum
		{"123", true},                  // too short
	}
	for _, tt := range tests {
		v := CreditCard(tt.input, "card")
		if v.Failed() != tt.wantErr {
			t.Errorf("CreditCard(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestLatitude(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"0", false},
		{"45.5", false},
		{"-90", false},
		{"90", false},
		{"91", true},
		{"-91", true},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := Latitude(tt.input, "lat")
		if v.Failed() != tt.wantErr {
			t.Errorf("Latitude(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestLongitude(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"0", false},
		{"-122.4", false},
		{"-180", false},
		{"180", false},
		{"181", true},
		{"-181", true},
		{"invalid", true},
	}
	for _, tt := range tests {
		v := Longitude(tt.input, "lng")
		if v.Failed() != tt.wantErr {
			t.Errorf("Longitude(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestCountryCode2(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"US", false},
		{"GB", false},
		{"us", true},  // lowercase
		{"USA", true}, // too long
		{"U", true},   // too short
	}
	for _, tt := range tests {
		v := CountryCode2(tt.input, "country")
		if v.Failed() != tt.wantErr {
			t.Errorf("CountryCode2(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestCountryCode3(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"USA", false},
		{"GBR", false},
		{"US", true}, // too short
	}
	for _, tt := range tests {
		v := CountryCode3(tt.input, "country")
		if v.Failed() != tt.wantErr {
			t.Errorf("CountryCode3(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestLanguageCode(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"en", false},
		{"fr", false},
		{"EN", true},  // uppercase
		{"eng", true}, // too long
	}
	for _, tt := range tests {
		v := LanguageCode(tt.input, "lang")
		if v.Failed() != tt.wantErr {
			t.Errorf("LanguageCode(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestCurrencyCode(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"USD", false},
		{"EUR", false},
		{"usd", true}, // lowercase
		{"US", true},  // too short
	}
	for _, tt := range tests {
		v := CurrencyCode(tt.input, "currency")
		if v.Failed() != tt.wantErr {
			t.Errorf("CurrencyCode(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestHex(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"deadbeef", false},
		{"DEADBEEF", false},
		{"0123456789abcdef", false},
		{"ghij", true},
	}
	for _, tt := range tests {
		v := Hex(tt.input, "hex")
		if v.Failed() != tt.wantErr {
			t.Errorf("Hex(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestDataURI(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"data:text/plain,hello", false},
		{"data:text/plain;base64,SGVsbG8=", false},
		{"not-a-data-uri", true},
		{"data:missing-comma", true},
	}
	for _, tt := range tests {
		v := DataURI(tt.input, "uri")
		if v.Failed() != tt.wantErr {
			t.Errorf("DataURI(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}

func TestFilePath(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"/path/to/file", false},
		{"file.txt", false},
		{"", true},
	}
	for _, tt := range tests {
		v := FilePath(tt.input, "path")
		if v.Failed() != tt.wantErr {
			t.Errorf("FilePath(%q) failed = %v, wantErr %v", tt.input, v.Failed(), tt.wantErr)
		}
	}
}
