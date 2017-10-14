package environment

type World struct {
	Size int
	BreakdownRate float64
	TerrainField TerrainField
	FluidField FluidField
	EnergyField EnergyField
	Blocks [][]Block
	Season float64
}

type Block struct {
	Solid bool
	Chemicals [8]float64
	ChemicalCapacity float64
}

type TerrainField struct {
	SolidityCurrent [][]float64
	SolidityConstant [][]float64
	SolidityAmplitude [][]float64
	SolidityPhase [][]float64

	HardnessCurrent [][]float64
	HardnessConstant [][]float64
	HardnessAmplitude [][]float64
	HardnessPhase [][]float64

	ChemicalCapacityCurrent [][]float64
	ChemicalCapacityConstant [][]float64
	ChemicalCapacityAmplitude [][]float64
	ChemicalCapacityPhase [][]float64
}

type FluidField struct {
	FlowBiasCurrent [][][2]float64
	FlowBiasConstant [][][2]float64
	FlowBiasAmplitude [][][2]float64
	FlowBiasPhase [][][2]float64
}

type EnergyField struct {
	PowerCurrent [][]float64
	PowerConstant [][]float64
	PowerAmplitude [][]float64
	PowerPhase [][]float64
}

type Rock struct {
	Block
	Health float64
}

type Water struct {
	Block
	Pressure float64
	NextPressure float64
	FlowRate [2]float64
}
