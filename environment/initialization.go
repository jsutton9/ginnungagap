package environment

import (
)

const ROCK_PORTION = 0.3
const BASE_HARDNESS = 100.0

func makeTerrainField(size int) TerrainField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := TerrainField{
		SolidityCurrent: zeroSquare(size),
		SolidityConstant: applyGridDistribution(randomSquare(size), logNormal),
		SolidityAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		SolidityPhase: applyGridDistribution(randomSquare(size), uniformConverter),

		HardnessCurrent: zeroSquare(size),
		HardnessConstant: applyGridDistribution(randomSquare(size), logNormal),
		HardnessAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		HardnessPhase: applyGridDistribution(randomSquare(size), uniformConverter),

		ChemicalCapacityCurrent: zeroSquare(size),
		ChemicalCapacityConstant: applyGridDistribution(randomSquare(size), logNormal),
		ChemicalCapacityAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		ChemicalCapacityPhase: applyGridDistribution(randomSquare(size), uniformConverter),
	}

	return field
}

func makeFluidField(size int) FluidField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := FluidField{
		FlowBiasCurrent: zeroVectorSquare(size),
		FlowBiasConstant: applyVectorMagnitudeDistribution(
			randomVectorSquare(size), logNormal),
		FlowBiasAmplitude: applyVectorMagnitudeDistribution(
			randomVectorSquare(size), logNormal),
		FlowBiasPhase: applyVectorMagnitudeDistribution(
			randomVectorSquare(size), uniformConverter),
	}

	return field
}

func makeEnergyField(size int) EnergyField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := EnergyField{
		PowerCurrent: zeroSquare(size),
		PowerConstant: applyGridDistribution(randomSquare(size), logNormal),
		PowerAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		PowerPhase: applyGridDistribution(randomSquare(size), uniformConverter),
	}

	return field
}

/*func makeBlocks(terrain TerrainField) [][]Block {
	blocks := make([][]Block, len(terrain))
	for i:=0; i<len(terrain); i++ {
		blocks[i] = make([]Block, len(terrain[i]))
	}
}*/
