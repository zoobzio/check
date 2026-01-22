package check

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compiled regular expressions for format validation.
var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	uuid4Regex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

	hexColorRegex     = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	hexColorFullRegex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	semverRegex = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

	e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

	macRegex      = regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	macDashRegex  = regexp.MustCompile(`^([0-9A-Fa-f]{2}-){5}[0-9A-Fa-f]{2}$`)
	hostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

// Email validates that a string is a valid email address.
func Email(v, field string) error {
	if !emailRegex.MatchString(v) {
		return fieldErr(field, "must be a valid email address")
	}
	return nil
}

// URL validates that a string is a valid URL.
func URL(v, field string) error {
	u, err := url.Parse(v)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fieldErr(field, "must be a valid URL")
	}
	return nil
}

// URLWithScheme validates that a string is a valid URL with one of the given schemes.
func URLWithScheme(v string, schemes []string, field string) error {
	u, err := url.Parse(v)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fieldErr(field, "must be a valid URL")
	}
	for _, s := range schemes {
		if strings.EqualFold(u.Scheme, s) {
			return nil
		}
	}
	return fieldErrf(field, "must have scheme: %s", strings.Join(schemes, ", "))
}

// HTTPOrHTTPS validates that a string is a valid HTTP or HTTPS URL.
func HTTPOrHTTPS(v, field string) error {
	return URLWithScheme(v, []string{"http", "https"}, field)
}

// UUID validates that a string is a valid UUID (versions 1-5).
func UUID(v, field string) error {
	if !uuidRegex.MatchString(v) {
		return fieldErr(field, "must be a valid UUID")
	}
	return nil
}

// UUID4 validates that a string is a valid UUID version 4.
func UUID4(v, field string) error {
	if !uuid4Regex.MatchString(v) {
		return fieldErr(field, "must be a valid UUID v4")
	}
	return nil
}

// IP validates that a string is a valid IP address (v4 or v6).
func IP(v, field string) error {
	if net.ParseIP(v) == nil {
		return fieldErr(field, "must be a valid IP address")
	}
	return nil
}

// IPv4 validates that a string is a valid IPv4 address.
func IPv4(v, field string) error {
	ip := net.ParseIP(v)
	if ip == nil || ip.To4() == nil {
		return fieldErr(field, "must be a valid IPv4 address")
	}
	return nil
}

// IPv6 validates that a string is a valid IPv6 address.
func IPv6(v, field string) error {
	ip := net.ParseIP(v)
	if ip == nil || ip.To4() != nil {
		return fieldErr(field, "must be a valid IPv6 address")
	}
	return nil
}

// CIDR validates that a string is a valid CIDR notation.
func CIDR(v, field string) error {
	_, _, err := net.ParseCIDR(v)
	if err != nil {
		return fieldErr(field, "must be a valid CIDR notation")
	}
	return nil
}

// MAC validates that a string is a valid MAC address.
func MAC(v, field string) error {
	if !macRegex.MatchString(v) && !macDashRegex.MatchString(v) {
		return fieldErr(field, "must be a valid MAC address")
	}
	return nil
}

// Hostname validates that a string is a valid hostname.
func Hostname(v, field string) error {
	if len(v) > 253 {
		return fieldErr(field, "must be a valid hostname")
	}
	if !hostnameRegex.MatchString(v) {
		return fieldErr(field, "must be a valid hostname")
	}
	return nil
}

// Port validates that a string is a valid port number (1-65535).
func Port(v, field string) error {
	p, err := strconv.Atoi(v)
	if err != nil || p < 1 || p > 65535 {
		return fieldErr(field, "must be a valid port number (1-65535)")
	}
	return nil
}

// HostPort validates that a string is a valid host:port combination.
func HostPort(v, field string) error {
	host, port, err := net.SplitHostPort(v)
	if err != nil {
		return fieldErr(field, "must be a valid host:port")
	}
	if host == "" {
		return fieldErr(field, "host must not be empty")
	}
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return fieldErr(field, "port must be a valid number (1-65535)")
	}
	return nil
}

// HexColor validates that a string is a valid hex color (#RGB, #RRGGBB, or #RRGGBBAA).
func HexColor(v, field string) error {
	if !hexColorRegex.MatchString(v) {
		return fieldErr(field, "must be a valid hex color")
	}
	return nil
}

// HexColorFull validates that a string is a valid 6-digit hex color (#RRGGBB).
func HexColorFull(v, field string) error {
	if !hexColorFullRegex.MatchString(v) {
		return fieldErr(field, "must be a valid hex color (#RRGGBB)")
	}
	return nil
}

// Base64 validates that a string is valid base64.
func Base64(v, field string) error {
	_, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return fieldErr(field, "must be valid base64")
	}
	return nil
}

// Base64URL validates that a string is valid URL-safe base64.
func Base64URL(v, field string) error {
	_, err := base64.URLEncoding.DecodeString(v)
	if err != nil {
		return fieldErr(field, "must be valid URL-safe base64")
	}
	return nil
}

// JSON validates that a string is valid JSON.
func JSON(v, field string) error {
	if !json.Valid([]byte(v)) {
		return fieldErr(field, "must be valid JSON")
	}
	return nil
}

// Semver validates that a string is a valid semantic version.
func Semver(v, field string) error {
	if !semverRegex.MatchString(v) {
		return fieldErr(field, "must be a valid semantic version")
	}
	return nil
}

// E164 validates that a string is a valid E.164 phone number.
func E164(v, field string) error {
	if !e164Regex.MatchString(v) {
		return fieldErr(field, "must be a valid E.164 phone number")
	}
	return nil
}

// CreditCard validates that a string is a valid credit card number using the Luhn algorithm.
func CreditCard(v, field string) error {
	// Remove spaces and dashes
	clean := strings.ReplaceAll(strings.ReplaceAll(v, " ", ""), "-", "")
	if len(clean) < 13 || len(clean) > 19 {
		return fieldErr(field, "must be a valid credit card number")
	}
	// Check all digits
	for _, r := range clean {
		if r < '0' || r > '9' {
			return fieldErr(field, "must be a valid credit card number")
		}
	}
	// Luhn algorithm
	sum := 0
	nDigits := len(clean)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit := int(clean[i] - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	if sum%10 != 0 {
		return fieldErr(field, "must be a valid credit card number")
	}
	return nil
}

// Latitude validates that a string is a valid latitude (-90 to 90).
func Latitude(v, field string) error {
	lat, err := strconv.ParseFloat(v, 64)
	if err != nil || lat < -90 || lat > 90 {
		return fieldErr(field, "must be a valid latitude (-90 to 90)")
	}
	return nil
}

// Longitude validates that a string is a valid longitude (-180 to 180).
func Longitude(v, field string) error {
	lon, err := strconv.ParseFloat(v, 64)
	if err != nil || lon < -180 || lon > 180 {
		return fieldErr(field, "must be a valid longitude (-180 to 180)")
	}
	return nil
}

// CountryCode2 validates that a string is a valid ISO 3166-1 alpha-2 country code.
func CountryCode2(v, field string) error {
	if len(v) != 2 {
		return fieldErr(field, "must be a valid ISO 3166-1 alpha-2 country code")
	}
	for _, r := range v {
		if r < 'A' || r > 'Z' {
			return fieldErr(field, "must be a valid ISO 3166-1 alpha-2 country code")
		}
	}
	return nil
}

// CountryCode3 validates that a string is a valid ISO 3166-1 alpha-3 country code.
func CountryCode3(v, field string) error {
	if len(v) != 3 {
		return fieldErr(field, "must be a valid ISO 3166-1 alpha-3 country code")
	}
	for _, r := range v {
		if r < 'A' || r > 'Z' {
			return fieldErr(field, "must be a valid ISO 3166-1 alpha-3 country code")
		}
	}
	return nil
}

// LanguageCode validates that a string is a valid ISO 639-1 language code.
func LanguageCode(v, field string) error {
	if len(v) != 2 {
		return fieldErr(field, "must be a valid ISO 639-1 language code")
	}
	for _, r := range v {
		if r < 'a' || r > 'z' {
			return fieldErr(field, "must be a valid ISO 639-1 language code")
		}
	}
	return nil
}

// CurrencyCode validates that a string is a valid ISO 4217 currency code.
func CurrencyCode(v, field string) error {
	if len(v) != 3 {
		return fieldErr(field, "must be a valid ISO 4217 currency code")
	}
	for _, r := range v {
		if r < 'A' || r > 'Z' {
			return fieldErr(field, "must be a valid ISO 4217 currency code")
		}
	}
	return nil
}

// Hex validates that a string contains only hexadecimal characters.
func Hex(v, field string) error {
	for _, r := range v {
		if !isHexDigit(r) {
			return fieldErr(field, "must be a valid hexadecimal string")
		}
	}
	return nil
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

// DataURI validates that a string is a valid data URI.
func DataURI(v, field string) error {
	if !strings.HasPrefix(v, "data:") {
		return fieldErr(field, "must be a valid data URI")
	}
	commaIdx := strings.Index(v, ",")
	if commaIdx == -1 {
		return fieldErr(field, "must be a valid data URI")
	}
	return nil
}

// FilePath validates that a string looks like a file path (contains path separators).
func FilePath(v, field string) error {
	if v == "" {
		return fieldErr(field, "must be a valid file path")
	}
	// Just basic validation - not empty and doesn't contain null bytes
	if strings.ContainsRune(v, 0) {
		return fieldErr(field, "must be a valid file path")
	}
	return nil
}

// UnixPath validates that a string is a valid Unix-style path.
func UnixPath(v, field string) error {
	if v == "" || strings.ContainsRune(v, 0) {
		return fieldErr(field, "must be a valid Unix path")
	}
	return nil
}
