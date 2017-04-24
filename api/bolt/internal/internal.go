package internal

import (
	"github.com/sparcs-kaist/whale"

	"encoding/binary"
	"encoding/json"
)

// MarshalUser encodes a user to binary format.
func MarshalUser(user *whale.User) ([]byte, error) {
	return json.Marshal(user)
}

// UnmarshalUser decodes a user from a binary data.
func UnmarshalUser(data []byte, user *whale.User) error {
	return json.Unmarshal(data, user)
}

// MarshalEndpoint encodes an endpoint to binary format.
func MarshalEndpoint(endpoint *whale.Endpoint) ([]byte, error) {
	return json.Marshal(endpoint)
}

// UnmarshalEndpoint decodes an endpoint from a binary data.
func UnmarshalEndpoint(data []byte, endpoint *whale.Endpoint) error {
	return json.Unmarshal(data, endpoint)
}

// MarshalResourceControl encodes a resource control object to binary format.
func MarshalResourceControl(rc *whale.ResourceControl) ([]byte, error) {
	return json.Marshal(rc)
}

// UnmarshalResourceControl decodes a resource control object from a binary data.
func UnmarshalResourceControl(data []byte, rc *whale.ResourceControl) error {
	return json.Unmarshal(data, rc)
}

// Itob returns an 8-byte big endian representation of v.
// This function is typically used for encoding integer IDs to byte slices
// so that they can be used as BoltDB keys.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
