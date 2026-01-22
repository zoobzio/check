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
		err := Email(tt.input, "email")
		if (err != nil) != tt.wantErr {
			t.Errorf("Email(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := URL(tt.input, "url")
		if (err != nil) != tt.wantErr {
			t.Errorf("URL(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := HTTPOrHTTPS(tt.input, "url")
		if (err != nil) != tt.wantErr {
			t.Errorf("HTTPOrHTTPS(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := UUID(tt.input, "uuid")
		if (err != nil) != tt.wantErr {
			t.Errorf("UUID(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := UUID4(tt.input, "uuid")
		if (err != nil) != tt.wantErr {
			t.Errorf("UUID4(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := IP(tt.input, "ip")
		if (err != nil) != tt.wantErr {
			t.Errorf("IP(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := IPv4(tt.input, "ip")
		if (err != nil) != tt.wantErr {
			t.Errorf("IPv4(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := IPv6(tt.input, "ip")
		if (err != nil) != tt.wantErr {
			t.Errorf("IPv6(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := CIDR(tt.input, "cidr")
		if (err != nil) != tt.wantErr {
			t.Errorf("CIDR(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := MAC(tt.input, "mac")
		if (err != nil) != tt.wantErr {
			t.Errorf("MAC(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Hostname(tt.input, "hostname")
		if (err != nil) != tt.wantErr {
			t.Errorf("Hostname(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Port(tt.input, "port")
		if (err != nil) != tt.wantErr {
			t.Errorf("Port(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := HostPort(tt.input, "hostport")
		if (err != nil) != tt.wantErr {
			t.Errorf("HostPort(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := HexColor(tt.input, "color")
		if (err != nil) != tt.wantErr {
			t.Errorf("HexColor(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Base64(tt.input, "data")
		if (err != nil) != tt.wantErr {
			t.Errorf("Base64(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := JSON(tt.input, "json")
		if (err != nil) != tt.wantErr {
			t.Errorf("JSON(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Semver(tt.input, "version")
		if (err != nil) != tt.wantErr {
			t.Errorf("Semver(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := E164(tt.input, "phone")
		if (err != nil) != tt.wantErr {
			t.Errorf("E164(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestCreditCard(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"4111111111111111", false},        // valid Visa
		{"4111 1111 1111 1111", false},     // with spaces
		{"4111-1111-1111-1111", false},     // with dashes
		{"4111111111111112", true},         // invalid checksum
		{"123", true},                       // too short
	}
	for _, tt := range tests {
		err := CreditCard(tt.input, "card")
		if (err != nil) != tt.wantErr {
			t.Errorf("CreditCard(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Latitude(tt.input, "lat")
		if (err != nil) != tt.wantErr {
			t.Errorf("Latitude(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Longitude(tt.input, "lng")
		if (err != nil) != tt.wantErr {
			t.Errorf("Longitude(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		{"us", true},   // lowercase
		{"USA", true},  // too long
		{"U", true},    // too short
	}
	for _, tt := range tests {
		err := CountryCode2(tt.input, "country")
		if (err != nil) != tt.wantErr {
			t.Errorf("CountryCode2(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		{"US", true},   // too short
	}
	for _, tt := range tests {
		err := CountryCode3(tt.input, "country")
		if (err != nil) != tt.wantErr {
			t.Errorf("CountryCode3(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		{"EN", true},   // uppercase
		{"eng", true},  // too long
	}
	for _, tt := range tests {
		err := LanguageCode(tt.input, "lang")
		if (err != nil) != tt.wantErr {
			t.Errorf("LanguageCode(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		{"usd", true},  // lowercase
		{"US", true},   // too short
	}
	for _, tt := range tests {
		err := CurrencyCode(tt.input, "currency")
		if (err != nil) != tt.wantErr {
			t.Errorf("CurrencyCode(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := Hex(tt.input, "hex")
		if (err != nil) != tt.wantErr {
			t.Errorf("Hex(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := DataURI(tt.input, "uri")
		if (err != nil) != tt.wantErr {
			t.Errorf("DataURI(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
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
		err := FilePath(tt.input, "path")
		if (err != nil) != tt.wantErr {
			t.Errorf("FilePath(%q) = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}
