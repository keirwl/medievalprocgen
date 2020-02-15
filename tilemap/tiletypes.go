package tilemap

type TileType int

const (
	Grass TileType = iota
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
