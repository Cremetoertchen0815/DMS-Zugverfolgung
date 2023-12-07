package bl

type TrainId int

type Vector2 struct {
	X float64
	Y float64
}

type Train struct {
	Id             TrainId
	LastCheckpoint *Checkpoint
	Velocity       Vector2
	Length         float64
}

type Checkpoint struct {
	DistanceFromStart float64
	Parent            *TrackSection
	Scanner           int
}

type TrackSection struct {
	Id          int
	Name        string
	Length      float64
	Checkpoints []Checkpoint
}

type TrackLink struct {
	LinkType          string
	ConnectedSections []int
}

type Track struct {
	Sections []*TrackSection
	Links    []*TrackLink
}
