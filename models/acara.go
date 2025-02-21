package models

import "time"

// Struct untuk data Acara
type Acara struct {
	ID                 string    `json:"id,omitempty" bson:"id"`
	NamaAcara          string    `json:"nama_acara" bson:"nama_acara"`
	TempatAcara        string    `json:"tempat_acara" bson:"tempat_acara"`
	LokasiPinGoogleMap string    `json:"lokasi_pin_google_map" bson:"lokasi_pin_google_map"`
	TanggalAcara       time.Time `json:"tanggal_acara" bson:"tanggal_acara"`
	Foto               string    `json:"foto" bson:"foto"`
	GoogleMapsEmbed    string    `json:"google_maps_embed" bson:"google_maps_embed"`
	Template           string    `json:"template" bson:"template"`
}
