package environment

type World struct {
	Size int
	TerrainField TerrainField
	PressureField PressureField
	EnergyField EnergyField
	Blocks [][]Block
	Season float64
}

type water struct {
	Pressure float64
	NextPressure float64
	Flow [][][2]float64
}

type rock struct {
	Health float64
}

type Block struct {
	Solid bool
	Chemicals [8]float64
	ChemicalCapacity float64
	Water water
	Rock rock
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

type PressureField struct {
	PressureBiasCurrent [][]float64
	PressureBiasConstant [][]float64
	PressureBiasAmplitude [][]float64
	PressureBiasPhase [][]float64
}

type EnergyField struct {
	PowerCurrent [][]float64
	PowerConstant [][]float64
	PowerAmplitude [][]float64
	PowerPhase [][]float64
}

type PhysicalConstants struct {
	RockPortion float64
	BaseHealth float64
	BaseCapacity float64
	BreakdownRate float64
}
