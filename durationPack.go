package hserial

type ReadDurationPack struct {
	HeadDuration int
	ContentDuration int
}

func InitReadDurationPack(d1, d2 int)ReadDurationPack{
	return ReadDurationPack{
		HeadDuration: d1,
		ContentDuration: d2,
	}
}