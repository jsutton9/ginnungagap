package environment

type Space struct {
	Size int
	BreakdownRate float64
	TerrainField TerrainField
	FluidField FluidField
	Blocks [][]Block
	Season float64
}

type Block struct {
	Liquid bool
	Chemicals [8]float64
	Pressure float64
	NextPressure float64
	FlowRate [2]float64
	MineralCapacity float64
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

	MineralCapacityCurrent [][]float64
	MineralCapacityConstant [][]float64
	MineralCapacityAmplitude [][]float64
	MineralCapacityPhase [][]float64
}

type FluidField struct {
	FlowBiasCurrent [][]float64
	FlowBiasConstant [][]float64
	FlowBiasAmplitude [][]float64
	FlowBiasPhase [][]float64
}

type EnergyField struct {
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
}
