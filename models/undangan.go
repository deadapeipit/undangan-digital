package models

// Struct untuk data Undangan
type Undangan struct {
	ID          string `json:"id,omitempty" bson:"id"`
	Nama        string `json:"nama" bson:"nama"`
	JumlahOrang int    `json:"jumlah_orang" bson:"jumlah_orang"`
	AcaraID     string `json:"acara_id" bson:"acara_id"`
	Hadir       bool   `json:"hadir" bson:"hadir"`
	Wishes      string `json:"wishes" bson:"wishes"`
}
