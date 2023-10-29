package pkg

type Message struct {
	Type    string  `json:"type"`
	ID      string  `json:"id"`
	Uri     string  `json:"uri"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	ForcePairing bool     `json:"forcePairing"`
	Volume       int      `json:"volume"`
	PairingType  string   `json:"pairingType"`
	Manifest     Manifest `json:"manifest"`
	ClientKey    string   `json:"client-key"`
	Message      string   `json:"message"`
}

type Manifest struct {
	ManifestVersion int         `json:"manifestVersion"`
	AppVersion      string      `json:"appVersion"`
	Signed          Signed      `json:"signed"`
	Permissions     []string    `json:"permissions"`
	Signatures      []Signature `json:"signatures"`
}

type Signed struct {
	Created              string            `json:"created"`
	AppId                string            `json:"appId"`
	VendorId             string            `json:"vendorId"`
	LocalizedAppNames    map[string]string `json:"localizedAppNames"`
	LocalizedVendorNames map[string]string `json:"localizedVendorNames"`
	Permissions          []string          `json:"permissions"`
	Serial               string            `json:"serial"`
}

type Signature struct {
	SignatureVersion int    `json:"signatureVersion"`
	Signature        string `json:"signature"`
}
