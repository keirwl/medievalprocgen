package tilemap

type TileType int

//go:generate stringer -type=TileType

const (
	None TileType = iota
	Grass
	GrassTrees
	GrassForest
	GrassHills
	GrassTreeHills
	Mountain
	WaterShallow
	WaterDeep
	GrassTown
	GrassCity
	GrassCastle
	Farm
	MarshTrees
	Marsh
	MarshLand
	MarshWater
	Snow
	SnowTrees
	SnowForest
	SnowHills
	SnowTreeHills
	SnowWater
	SnowTown
	SnowCastle
	Desert
	DesertHills
	DesertDunes
	DesertMountain
	DesertOasis
	DesertTown
	DesertCity
	DesertCastle
	EnchantedForest
	GrassCave
	SnowCave
	DesertCave
	PortLeft
	PortRight
	Lighthouse
	GrassRuin
	MarshGraveyard
	Max
)
