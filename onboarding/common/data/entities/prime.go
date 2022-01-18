package entities

import "time"

type PrimeDetails struct {
	GuesserID int64     `bson:"guesser_id"`
	Time      time.Time `bson:"time"`
	OriginNum int64     `bson:"origin_num"`
}

type Prime struct {
	Prime     int64     `bson:"prime"`
	PrimeDets []PrimeDetails `bson:"prime_dets"`
}
